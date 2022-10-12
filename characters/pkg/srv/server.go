package srv

import (
	"context"
	"fmt"
	accountspb "github.com/ShatteredRealms/Accounts/pkg/pb"
	accountssrv "github.com/ShatteredRealms/Accounts/pkg/srv"
	"github.com/ShatteredRealms/Characters/internal/log"
	"github.com/ShatteredRealms/Characters/internal/option"
	"github.com/ShatteredRealms/Characters/pkg/pb"
	"github.com/ShatteredRealms/Characters/pkg/service"
	"github.com/ShatteredRealms/GoUtils/pkg/interceptor"
	utilService "github.com/ShatteredRealms/GoUtils/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"time"
)

const (
	retryAfter = time.Second * 10
)

func NewServer(
	characterService service.CharacterService,
	jwt utilService.JWTService,
	logger log.LoggerService,
	config option.Config,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	conn, err := dialAuthorizationServer(config)
	if err != nil {
		return nil, nil, err
	}

	publicRPCs := make(map[string]struct{})

	authInterceptor := interceptor.NewAuthInterceptor(
		jwt,
		publicRPCs,
		accountssrv.GetPermissions(accountspb.NewAuthorizationServiceClient(conn), jwt, "shatteredrealmsonline.com/characters/v1"),
	)

	go processRoleUpdates(accountspb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, logger)
	go processUserUpdates(accountspb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), logger.UnaryLogRequest(log.Info)),
		grpc.ChainStreamInterceptor(authInterceptor.Stream(), logger.StreamLogRequest(log.Info)),
	)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	characterServiceServer := NewCharacterServiceServer(characterService, jwt)
	pb.RegisterCharactersServiceServer(grpcServer, characterServiceServer)
	err = pb.RegisterCharactersServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	return grpcServer, gwmux, nil
}

func dialAuthorizationServer(
	config option.Config,
) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return grpc.Dial(config.AccountsAddress(), opts...)
}

func processUserUpdates(
	authorizationClient accountspb.AuthorizationServiceClient,
	interceptor *interceptor.AuthInterceptor,
	jwtService utilService.JWTService,
	logger log.LoggerService,
) {
	userUpdatesClient, err := authorizationClient.SubscribeUserUpdates(serverAuthContext(jwtService), &emptypb.Empty{})
	if err != nil {
		logger.Logf(log.Error, "Unable to subscribe to user updates from authorization client. Retrying in %d seconds", retryAfter/time.Second)
		time.Sleep(retryAfter)
		processUserUpdates(authorizationClient, interceptor, jwtService, logger)
		return
	}
	logger.Info("Successfully subscribed to user updates from authorization server.")
	for {
		msg, err := userUpdatesClient.Recv()

		if err == nil {
			logger.Logf(log.Debug, "Update to user %d permissions. Clearing permissions cache for that user.", msg.Id)
			err = interceptor.ClearUserCache(uint(msg.Id))
			if err != nil {
				logger.Logf(log.Warning, "Clearing cache: %v", err)
			}
		} else if err == io.EOF {
			logger.Infof("User updates stream ended. Retrying in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			processUserUpdates(authorizationClient, interceptor, jwtService, logger)
			return
		} else {
			logger.Logf(log.Error, "User updates: %v.", err)
			logger.Infof("Retrying connection in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			processUserUpdates(authorizationClient, interceptor, jwtService, logger)
			return
		}
	}
}

func processRoleUpdates(
	authorizationClient accountspb.AuthorizationServiceClient,
	interceptor *interceptor.AuthInterceptor,
	jwtService utilService.JWTService,
	logger log.LoggerService,
) {
	roleUpdatesClient, err := authorizationClient.SubscribeRoleUpdates(serverAuthContext(jwtService), &emptypb.Empty{})
	if err != nil {
		logger.Logf(log.Error, "Unable to subscribe to role updates from authorization client. Retrying in %d seconds", retryAfter/time.Second)
		time.Sleep(retryAfter)
		processRoleUpdates(authorizationClient, interceptor, jwtService, logger)
		return
	}
	logger.Info("Successfully subscribed to role updates from authorization server.")
	for {
		msg, err := roleUpdatesClient.Recv()
		if err == nil {
			logger.Logf(log.Debug, "Update to role %d permissions. Clearing permissions cache for all users.", msg.Id)
			err = interceptor.ClearCache()
			if err != nil {
				logger.Logf(log.Warning, "Clearing cache: %v", err)
			}
		} else if err == io.EOF {
			logger.Infof("Role updates stream ended. Retrying in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			processRoleUpdates(authorizationClient, interceptor, jwtService, logger)
			return
		} else {
			logger.Logf(log.Error, "Role Updates: %v.", err)
			logger.Infof("Retrying connection in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			processRoleUpdates(authorizationClient, interceptor, jwtService, logger)
			return
		}
	}
}

func serverAuthContext(jwtService utilService.JWTService) context.Context {
	md := metadata.New(
		map[string]string{
			"authorization": fmt.Sprintf(
				"Bearer %s", generateTemporaryServerToken(jwtService, "shatteredrealmsonline.com/characters/v1"),
			),
		},
	)
	return metadata.NewOutgoingContext(context.Background(), md)
}

func generateTemporaryServerToken(jwtService utilService.JWTService, requestingHost string) string {
	out, _ := jwtService.Create(time.Second*10, requestingHost, jwt.MapClaims{"sub": 0})
	return out
}

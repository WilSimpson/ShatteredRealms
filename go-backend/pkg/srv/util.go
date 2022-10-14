package srv

import (
	"context"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"time"
)

const (
	retryAfter = time.Second * 10
)

func ConvertRolePermission(inPermission *model.RolePermission) *pb.UserPermission {
	return &pb.UserPermission{
		Permission: &wrapperspb.StringValue{Value: inPermission.Permission},
		Other:      inPermission.Other,
	}
}

func ConvertRolePermissions(inPermissions []*model.RolePermission) []*pb.UserPermission {
	permissions := make([]*pb.UserPermission, len(inPermissions))
	for i, permission := range inPermissions {
		permissions[i] = ConvertRolePermission(permission)
	}

	return permissions
}

func ConvertUserPermission(inPermission *model.UserPermission) *pb.UserPermission {
	return &pb.UserPermission{
		Permission: &wrapperspb.StringValue{Value: inPermission.Permission},
		Other:      inPermission.Other,
	}
}

func ConvertUserPermissions(inPermissions []*model.UserPermission) []*pb.UserPermission {
	permissions := make([]*pb.UserPermission, len(inPermissions))
	for i, permission := range inPermissions {
		permissions[i] = ConvertUserPermission(permission)
	}

	return permissions
}

func ConvertRoleWithoutPermissions(inRole *model.Role) *pb.UserRole {
	return &pb.UserRole{
		Id:   uint64(inRole.ID),
		Name: &wrapperspb.StringValue{Value: inRole.Name},
	}
}

func ConvertRolesWithoutPermissions(inRoles []*model.Role) []*pb.UserRole {
	roles := make([]*pb.UserRole, len(inRoles))
	for i, role := range inRoles {
		roles[i] = ConvertRoleWithoutPermissions(role)
	}

	return roles
}

func ConvertRoleNameOnly(inRole *model.Role) *pb.UserRole {
	return &pb.UserRole{
		Id:   uint64(inRole.ID),
		Name: &wrapperspb.StringValue{Value: inRole.Name},
	}
}

func ConvertRolesNamesOnly(inRoles []*model.Role) []*pb.UserRole {
	roles := make([]*pb.UserRole, len(inRoles))
	for i, role := range inRoles {
		roles[i] = ConvertRoleNameOnly(role)
	}

	return roles
}

func GetPermissions(
	authorizationService pb.AuthorizationServiceClient,
	jwtService service.JWTService,
	requestingHost string,
) func(userID uint) map[string]bool {
	return func(userID uint) map[string]bool {
		md := metadata.New(
			map[string]string{
				"authorization": fmt.Sprintf(
					"Bearer %s", generateTemporaryServerToken(jwtService, requestingHost),
				),
			},
		)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		authorizations, err := authorizationService.GetAuthorization(ctx, &pb.IDMessage{Id: uint64(userID)})

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return map[string]bool{}
		}

		resp := make(map[string]bool)

		for _, role := range authorizations.Roles {
			for _, rolePermission := range role.Permissions {
				resp[rolePermission.Permission.Value] = resp[rolePermission.Permission.Value] || rolePermission.Other
			}
		}

		for _, userPermission := range authorizations.Permissions {
			resp[userPermission.Permission.Value] = resp[userPermission.Permission.Value] || userPermission.Other
		}

		return resp
	}
}

func generateTemporaryServerToken(jwtService service.JWTService, requestingHost string) string {
	out, _ := jwtService.Create(time.Second, requestingHost, jwt.MapClaims{"sub": 0})
	return out
}
func ProcessUserUpdates(
	authorizationClient pb.AuthorizationServiceClient,
	interceptor *interceptor.AuthInterceptor,
	jwtService service.JWTService,
	serviceAuthName string,
) {
	userUpdatesClient, err := authorizationClient.SubscribeUserUpdates(serverAuthContext(jwtService, serviceAuthName), &emptypb.Empty{})
	if err != nil {
		log.Errorf("Unable to subscribe to user updates from authorization client. Retrying in %d seconds", retryAfter/time.Second)
		time.Sleep(retryAfter)
		ProcessUserUpdates(authorizationClient, interceptor, jwtService, serviceAuthName)
		return
	}
	log.Info("Successfully subscribed to user updates from authorization server.")
	for {
		msg, err := userUpdatesClient.Recv()

		if err == nil {
			log.Debugf("Update to user %d permissions. Clearing permissions cache for that user.", msg.Id)
			err = interceptor.ClearUserCache(uint(msg.Id))
			if err != nil {
				log.Warning("Clearing cache: %v", err)
			}
		} else if err == io.EOF {
			log.Infof("User updates stream ended. Retrying in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			ProcessUserUpdates(authorizationClient, interceptor, jwtService, serviceAuthName)
			return
		} else {
			log.Errorf("User updates: %v.", err)
			log.Infof("Retrying connection in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			ProcessUserUpdates(authorizationClient, interceptor, jwtService, serviceAuthName)
			return
		}
	}
}

func ProcessRoleUpdates(
	authorizationClient pb.AuthorizationServiceClient,
	interceptor *interceptor.AuthInterceptor,
	jwtService service.JWTService,
	serviceAuthName string,
) {
	roleUpdatesClient, err := authorizationClient.SubscribeRoleUpdates(serverAuthContext(jwtService, serviceAuthName), &emptypb.Empty{})
	if err != nil {
		log.Errorf("Unable to subscribe to role updates from authorization client. Retrying in %d seconds", retryAfter/time.Second)
		time.Sleep(retryAfter)
		ProcessRoleUpdates(authorizationClient, interceptor, jwtService, serviceAuthName)
		return
	}
	log.Info("Successfully subscribed to role updates from authorization server.")
	for {
		msg, err := roleUpdatesClient.Recv()
		if err == nil {
			log.Debug("Update to role %d permissions. Clearing permissions cache for all users.", msg.Id)
			err = interceptor.ClearCache()
			if err != nil {
				log.Warning("Clearing cache: %v", err)
			}
		} else if err == io.EOF {
			log.Infof("Role updates stream ended. Retrying in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			ProcessRoleUpdates(authorizationClient, interceptor, jwtService, serviceAuthName)
			return
		} else {
			log.Errorf("Role Updates: %v.", err)
			log.Infof("Retrying connection in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			ProcessRoleUpdates(authorizationClient, interceptor, jwtService, serviceAuthName)
			return
		}
	}
}

func serverAuthContext(jwtService service.JWTService, authorizer string) context.Context {
	md := metadata.New(
		map[string]string{
			"authorization": fmt.Sprintf(
				"Bearer %s", generateTemporaryServerToken(jwtService, authorizer),
			),
		},
	)
	return metadata.NewOutgoingContext(context.Background(), md)
}

func DialOtelGrpc(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, insecureOtelGrpcDialOpts()...)
}

func otelGrpcDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}
}

func insecureOtelGrpcDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}
}

func CreateGrpcServerWithAuth(jwt service.JWTService, accountsAddress string, publicRPCs map[string]struct{}) (*grpc.Server, *runtime.ServeMux, []grpc.DialOption, error) {
	if publicRPCs == nil {
		publicRPCs = make(map[string]struct{})
	}

	conn, err := grpc.Dial(
		accountsAddress,
		insecureOtelGrpcDialOpts()...,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	authInterceptor := interceptor.NewAuthInterceptor(
		jwt,
		publicRPCs,
		GetPermissions(pb.NewAuthorizationServiceClient(conn), jwt, "sro.com/characters/v1"),
	)

	go ProcessRoleUpdates(pb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, "sro.com/characters/v1")
	go ProcessUserUpdates(pb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, "sro.com/characters/v1")

	return grpc.NewServer(
			grpc.ChainUnaryInterceptor(authInterceptor.Unary(), helpers.UnaryLogRequest()),
			grpc.ChainStreamInterceptor(authInterceptor.Stream(), helpers.StreamLogRequest()),
		),
		runtime.NewServeMux(),
		insecureOtelGrpcDialOpts(),
		nil
}

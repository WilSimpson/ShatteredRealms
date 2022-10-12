package srv

import (
	"context"
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	"github.com/ShatteredRealms/Accounts/pkg/pb"
	accountService "github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/ShatteredRealms/GoUtils/pkg/interceptor"
	"github.com/ShatteredRealms/GoUtils/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"time"
)

func NewServer(
	u accountService.UserService,
	p accountService.PermissionService,
	r accountService.RoleService,
	jwt service.JWTService,
	logger log.LoggerService,
	config option.Config,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	publicRPCs := make(map[string]struct{})
	publicRPCs["/sro.accounts.HealthService/Health"] = struct{}{}
	publicRPCs["/sro.accounts.AuthenticationService/Login"] = struct{}{}
	publicRPCs["/sro.accounts.AuthenticationService/Register"] = struct{}{}

	authorizationServiceServer := NewAuthorizationServiceServer(u, p, r, logger)
	authInterceptor := interceptor.NewAuthInterceptor(jwt, publicRPCs, getPermissions(authorizationServiceServer))
	authorizationServiceServer.AddAuthInterceptor(authInterceptor)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), logger.UnaryLogRequest(log.Info)),
		grpc.ChainStreamInterceptor(authInterceptor.Stream(), logger.StreamLogRequest(log.Info)),
	)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	authenticationServiceServer := NewAuthenticationServiceServer(u, jwt, p, logger)
	pb.RegisterAuthenticationServiceServer(grpcServer, authenticationServiceServer)
	err := pb.RegisterAuthenticationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	userServiceServer := NewUserServiceServer(u, p, jwt, logger)
	pb.RegisterUserServiceServer(grpcServer, userServiceServer)
	err = pb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	healthServiceServer := NewHealthServiceServer()
	pb.RegisterHealthServiceServer(grpcServer, healthServiceServer)
	err = pb.RegisterHealthServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	pb.RegisterAuthorizationServiceServer(grpcServer, authorizationServiceServer)
	err = pb.RegisterAuthorizationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	// Compute the AllPermissions method once and save in memory
	authorizationServiceServer.SetupAllPermissions(grpcServer.GetServiceInfo())

	return grpcServer, gwmux, nil
}

func getPermissions(
	server *authorizationServiceServer,
) func(userID uint) map[string]bool {
	return func(userID uint) map[string]bool {
		// UserID 0 is for server communication
		if userID == 0 {
			allPerms := make(map[string]bool, len(server.allPermissions.Permissions))
			for _, perm := range server.allPermissions.Permissions {
				allPerms[perm.Permission.Value] = true
			}

			return allPerms
		}

		user := server.userService.FindById(userID)
		if user == nil || !user.Exists() {
			return map[string]bool{}
		}

		resp := make(map[string]bool)

		for _, role := range user.Roles {
			for _, rolePermission := range server.permissionService.FindPermissionsForRoleID(role.ID) {
				resp[rolePermission.Permission] = resp[rolePermission.Permission] || rolePermission.Other
			}
		}

		for _, userPermission := range server.permissionService.FindPermissionsForUserID(userID) {
			resp[userPermission.Permission] = resp[userPermission.Permission] || userPermission.Other
		}

		return resp
	}
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

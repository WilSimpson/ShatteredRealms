package main

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	accountService "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/srv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(
	u accountService.UserService,
	p accountService.PermissionService,
	r accountService.RoleService,
	jwt service.JWTService,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	publicRPCs := make(map[string]struct{})
	publicRPCs["/sro.accounts.HealthService/Health"] = struct{}{}
	publicRPCs["/sro.accounts.AuthenticationService/Login"] = struct{}{}
	publicRPCs["/sro.accounts.AuthenticationService/Register"] = struct{}{}

	authorizationServiceServer := srv.NewAuthorizationServiceServer(u, p, r)
	authInterceptor := interceptor.NewAuthInterceptor(jwt, publicRPCs, getPermissions(authorizationServiceServer))
	authorizationServiceServer.AddAuthInterceptor(authInterceptor)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), helpers.UnaryLogRequest()),
		grpc.ChainStreamInterceptor(authInterceptor.Stream(), helpers.StreamLogRequest()),
	)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	authenticationServiceServer := srv.NewAuthenticationServiceServer(u, jwt, p)
	pb.RegisterAuthenticationServiceServer(grpcServer, authenticationServiceServer)
	err := pb.RegisterAuthenticationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	userServiceServer := srv.NewUserServiceServer(u, p, jwt)
	pb.RegisterUserServiceServer(grpcServer, userServiceServer)
	err = pb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	healthServiceServer := srv.NewHealthServiceServer()
	pb.RegisterHealthServiceServer(grpcServer, healthServiceServer)
	err = pb.RegisterHealthServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	pb.RegisterAuthorizationServiceServer(grpcServer, authorizationServiceServer)
	err = pb.RegisterAuthorizationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Address(),
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
	server *srv.AuthorizationServiceServer,
) func(userID uint) map[string]bool {
	return func(userID uint) map[string]bool {
		// UserID 0 is for server communication
		if userID == 0 {
			allPerms := make(map[string]bool, len(server.AllPermissions.Permissions))
			for _, perm := range server.AllPermissions.Permissions {
				allPerms[perm.Permission.Value] = true
			}

			return allPerms
		}

		user := server.UserService.FindById(userID)
		if user == nil || !user.Exists() {
			return map[string]bool{}
		}

		resp := make(map[string]bool)

		for _, role := range user.Roles {
			for _, rolePermission := range server.PermissionService.FindPermissionsForRoleID(role.ID) {
				resp[rolePermission.Permission] = resp[rolePermission.Permission] || rolePermission.Other
			}
		}

		for _, userPermission := range server.PermissionService.FindPermissionsForUserID(userID) {
			resp[userPermission.Permission] = resp[userPermission.Permission] || userPermission.Other
		}

		return resp
	}
}

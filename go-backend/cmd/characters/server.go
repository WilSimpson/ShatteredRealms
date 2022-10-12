package main

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/srv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(
	characterService service.CharacterService,
	jwt service.JWTService,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	conn, err := grpc.Dial(conf.AccountsAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	publicRPCs := make(map[string]struct{})

	authInterceptor := interceptor.NewAuthInterceptor(
		jwt,
		publicRPCs,
		srv.GetPermissions(pb.NewAuthorizationServiceClient(conn), jwt, "sro.com/characters/v1"),
	)

	go srv.ProcessRoleUpdates(pb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, "sro.com/characters/v1")
	go srv.ProcessUserUpdates(pb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, "sro.com/characters/v1")

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), helpers.UnaryLogRequest()),
		grpc.ChainStreamInterceptor(authInterceptor.Stream(), helpers.StreamLogRequest()),
	)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	characterServiceServer := srv.NewCharacterServiceServer(characterService, jwt)
	pb.RegisterCharactersServiceServer(grpcServer, characterServiceServer)
	err = pb.RegisterCharactersServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	return grpcServer, gwmux, nil
}

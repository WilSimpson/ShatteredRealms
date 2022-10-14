package main

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/srv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func NewServer(
	characterService service.CharacterService,
	jwt service.JWTService,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	grpcServer, gwmux, opts, err := srv.CreateGrpcServerWithAuth(jwt, conf.Accounts.Address(), nil)

	characterServiceServer := srv.NewCharacterServiceServer(characterService, jwt)
	pb.RegisterCharactersServiceServer(grpcServer, characterServiceServer)
	err = pb.RegisterCharactersServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Accounts.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	return grpcServer, gwmux, nil
}

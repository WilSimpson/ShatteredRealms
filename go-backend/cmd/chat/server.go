package main

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func NewServer(
	jwt service.JWTService,
) (*grpc.Server, *runtime.ServeMux, error) {
	//ctx := context.Background()
	//
	//conn, err := srv.DialClientWithTelemetry(conf.Accounts.Address())
	//if err != nil {
	//    return nil, nil, err
	//}
	//
	//publicRPCs := make(map[string]struct{})
	//
	//authInterceptor := interceptor.NewAuthInterceptor(
	//    jwt,
	//    publicRPCs,
	//    srv.GetPermissions(pb.NewAuthorizationServiceClient(conn), jwt, "sro.com/chat/v1"),
	//)
	//
	//go srv.ProcessRoleUpdates(pb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, "sro.com/chat/v1")
	//go srv.ProcessUserUpdates(pb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, "sro.com/chat/v1")
	//
	//grpcServer := grpc.NewServer(
	//    grpc.ChainUnaryInterceptor(authInterceptor.Unary(), helpers.UnaryLogRequest()),
	//    grpc.ChainStreamInterceptor(authInterceptor.Stream(), helpers.StreamLogRequest()),
	//)
	//
	//gwmux := runtime.NewServeMux()
	//opts := []grpc.DialOption{
	//    grpc.WithTransportCredentials(insecure.NewCredentials()),
	//}

	return nil, nil, nil
}

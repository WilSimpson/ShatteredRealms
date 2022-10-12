package main

import (
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
)

func main() {
	jwtService, err := service.NewJWTService(conf.KeyDir)
	if err != nil {
		log.Error(fmt.Sprintf("jwt service: %v", err))
		os.Exit(1)
	}

	grpcServer, gwmux, err := NewServer(jwtService)
	if err != nil {
		log.Error(fmt.Sprintf("grpc server: %v", err))
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", conf.Address())
	if err != nil {
		log.Error(fmt.Sprintf("listen: %v", err))
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    conf.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")

	err = server.Serve(lis)

	if err != nil {
		log.Error(fmt.Sprintf("listen: %v", err))
		os.Exit(1)
	}
}

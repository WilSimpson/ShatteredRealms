package main

import (
    "fmt"
    "github.com/ShatteredRealms/GoUtils/pkg/helpers"
    utilService "github.com/ShatteredRealms/GoUtils/pkg/service"
    "github.com/ShatteredRealms/gamebackend/internal/log"
    "github.com/ShatteredRealms/gamebackend/internal/option"
    "github.com/ShatteredRealms/gamebackend/pkg/srv"
    "net"
    "net/http"
    "os"
)

func main() {
    config := option.NewConfig()

    var logger log.LoggerService
    if config.IsRelease() {
        logger = log.NewLogger(log.Info, "")
    } else {
        logger = log.NewLogger(log.Debug, "")
    }

    jwtService, err := utilService.NewJWTService(config.KeyDir.Value)
    if err != nil {
        logger.Log(log.Error, fmt.Sprintf("jwt service: %v", err))
        os.Exit(1)
    }

    grpcServer, gwmux, err := srv.NewServer(jwtService, config, logger)
    if err != nil {
        logger.Log(log.Error, fmt.Sprintf("grpc server: %v", err))
        os.Exit(1)
    }

    lis, err := net.Listen("tcp", config.Address())
    if err != nil {
        logger.Log(log.Error, fmt.Sprintf("listen: %v", err))
        os.Exit(1)
    }

    server := &http.Server{
        Addr:    config.Address(),
        Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
    }

    logger.Info("Server starting")

    err = server.Serve(lis)

    if err != nil {
        logger.Log(log.Error, fmt.Sprintf("listen: %v", err))
        os.Exit(1)
    }
}

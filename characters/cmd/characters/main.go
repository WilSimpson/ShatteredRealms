package main

import (
    "fmt"
    "github.com/ShatteredRealms/Characters/internal/log"
    "github.com/ShatteredRealms/Characters/internal/option"
    repository2 "github.com/ShatteredRealms/Characters/pkg/repository"
    "github.com/ShatteredRealms/Characters/pkg/service"
    "github.com/ShatteredRealms/Characters/pkg/srv"
    "github.com/ShatteredRealms/GoUtils/pkg/helpers"
    "github.com/ShatteredRealms/GoUtils/pkg/repository"
    utilService "github.com/ShatteredRealms/GoUtils/pkg/service"
    "gopkg.in/yaml.v3"
    "io/ioutil"
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

    file, err := ioutil.ReadFile(*config.DBFile.Value)
    if err != nil {
        logger.Log(log.Error, fmt.Sprintf("reading db file: %v", err))
        os.Exit(1)
    }

    c := &repository.DBConnections{}
    err = yaml.Unmarshal(file, c)
    if err != nil {
        logger.Log(log.Error, fmt.Sprintf("parsing db file: %v", err))
        os.Exit(1)
    }

    db, err := repository.DBConnect(*c)
    if err != nil {
        logger.Log(log.Error, fmt.Sprintf("db: %v", err))
        os.Exit(1)
    }

    jwtService, err := utilService.NewJWTService(*config.KeyDir.Value)
    if err != nil {
        logger.Log(log.Error, fmt.Sprintf("jwt service: %v", err))
        os.Exit(1)
    }

    characterRepo := repository2.NewCharacterRepository(db)
    if err := characterRepo.Migrate(); err != nil {
        logger.Log(log.Error, fmt.Sprintf("character repo: %v", err))
        os.Exit(1)
    }
    characterService := service.NewCharacterService(characterRepo)

    grpcServer, gwmux, err := srv.NewServer(characterService, jwtService, logger, config)
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

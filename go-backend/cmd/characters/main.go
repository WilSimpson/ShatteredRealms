package main

import (
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

func main() {
	file, err := ioutil.ReadFile(conf.DBFile)
	if err != nil {
		log.Error(fmt.Sprintf("reading db file: %v", err))
		os.Exit(1)
	}

	c := &repository.DBConnections{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Error(fmt.Sprintf("parsing db file: %v", err))
		os.Exit(1)
	}

	db, err := repository.DBConnect(*c)
	if err != nil {
		log.Error(fmt.Sprintf("db: %v", err))
		os.Exit(1)
	}

	jwtService, err := service.NewJWTService(conf.KeyDir)
	if err != nil {
		log.Error(fmt.Sprintf("jwt service: %v", err))
		os.Exit(1)
	}

	characterRepo := repository.NewCharacterRepository(db)
	if err := characterRepo.Migrate(); err != nil {
		log.Error(fmt.Sprintf("character repo: %v", err))
		os.Exit(1)
	}
	characterService := service.NewCharacterService(characterRepo)

	grpcServer, gwmux, err := NewServer(characterService, jwtService)
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

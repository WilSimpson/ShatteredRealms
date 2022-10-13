package main

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/config"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type appConfig struct {
	Characters config.Server `yaml:"characters"`
	Accounts   config.Server `yaml:"accounts"`
	KeyDir     string        `yaml:"keyDir"`
	DBFile     string        `yaml:"dbFile"`
}

var (
	conf = &appConfig{
		Characters: config.Server{
			Port:     8081,
			Host:     "",
			Mode:     "development",
			LogLevel: log.InfoLevel,
		},
		Accounts: config.Server{
			Port: 8080,
			Host: "",
		},
		KeyDir: "/etc/sro/auth",
		DBFile: "/etc/sro/db.yaml",
	}
)

func main() {
	db, err := repository.ConnectFromFile(conf.DBFile)
	helpers.Check(err, "db connect from file")

	characterRepo := repository.NewCharacterRepository(db)
	helpers.Check(characterRepo.Migrate(), "character repo")

	characterService := service.NewCharacterService(characterRepo)
	jwtService, err := service.NewJWTService(conf.KeyDir)
	helpers.Check(err, "jwt service")

	grpcServer, gwmux, err := NewServer(characterService, jwtService)
	helpers.Check(err, "create grpc server")

	lis, err := net.Listen("tcp", conf.Characters.Address())
	helpers.Check(err, "listen")

	server := &http.Server{
		Addr:    conf.Characters.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")

	err = server.Serve(lis)
	helpers.Check(err, "serve")
}

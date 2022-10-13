package main

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/config"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type appConfig struct {
	GameBackend config.Server `yaml:"gameBackend"`
	Accounts    config.Server `yaml:"accounts"`
	Characters  config.Server `yaml:"characters"`
	KeyDir      string        `yaml:"keyDir"`
	DBFile      string        `yaml:"dbFile"`
	Agones      agonesConfig  `yaml:"agones"`
}

type agonesConfig struct {
	KeyFile    string        `yaml:"keyFile"`
	CertFile   string        `yaml:"certFile"`
	CaCertFile string        `yaml:"caCertFile"`
	Namespace  string        `yaml:"namespace"`
	Allocator  config.Server `yaml:"allocator"`
}

var (
	conf = &appConfig{
		GameBackend: config.Server{
			Port:     8082,
			Host:     "",
			Mode:     "development",
			LogLevel: log.InfoLevel,
		},
		Characters: config.Server{
			Port: 8081,
			Host: "",
		},
		Accounts: config.Server{
			Port: 8080,
			Host: "",
		},
		KeyDir: "/etc/sro/auth",
		DBFile: "/etc/sro/db.yaml",
		Agones: agonesConfig{
			KeyFile:    "/etc/sro/auth/agones/client/key",
			CertFile:   "/etc/sro/auth/agones/client/key",
			CaCertFile: "/etc/sro/auth/agones/ca/ca",
			Namespace:  "default",
			Allocator: config.Server{
				Port: 443,
				Host: "",
			},
		},
	}
)

func main() {
	jwtService, err := service.NewJWTService(conf.KeyDir)
	helpers.Check(err, "jwt service")

	grpcServer, gwmux, err := NewServer(jwtService)
	helpers.Check(err, "create grpc server")

	lis, err := net.Listen("tcp", conf.GameBackend.Address())
	helpers.Check(err, "listen")

	server := &http.Server{
		Addr:    conf.GameBackend.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")

	err = server.Serve(lis)
	helpers.Check(err, "serve")
}

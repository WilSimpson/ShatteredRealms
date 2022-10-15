package main

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/config"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
	"net"
	"net/http"
)

type appConfig struct {
	GameBackend config.Server        `yaml:"gameBackend"`
	Accounts    config.Server        `yaml:"accounts"`
	Characters  config.Server        `yaml:"characters"`
	KeyDir      string               `yaml:"keyDir"`
	DBFile      string               `yaml:"dbFile"`
	Agones      agonesConfig         `yaml:"agones"`
	Uptrace     config.UptraceConfig `yaml:"uptrace"`
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
			Local: config.ServerAddress{
				Port: 8082,
				Host: "",
			},
			Remote: config.ServerAddress{
				Port: 8082,
				Host: "",
			},
			Mode:     "development",
			LogLevel: log.InfoLevel,
			DB: repository.DBPoolConfig{
				Master: repository.DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "gamebackend",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []repository.DBConfig{},
			},
		},
		Characters: config.Server{
			Remote: config.ServerAddress{
				Port: 8081,
				Host: "",
			},
		},
		Accounts: config.Server{
			Remote: config.ServerAddress{
				Port: 8080,
				Host: "",
			},
		},
		KeyDir: "/etc/sro/auth",
		DBFile: "/etc/sro/db.yaml",
		Agones: agonesConfig{
			KeyFile:    "/etc/sro/auth/agones/client/key",
			CertFile:   "/etc/sro/auth/agones/client/cert",
			CaCertFile: "/etc/sro/auth/agones/ca/ca",
			Namespace:  "default",
			Allocator: config.Server{
				Remote: config.ServerAddress{
					Port: 443,
					Host: "",
				},
			},
		},
	}
)

func init() {
	helpers.SetupLogs()
	config.SetupConfig(conf)
	log.Infof("config: %+v", *conf)
}

func main() {
	ctx := context.Background()
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(conf.Uptrace.DSN()),
		uptrace.WithServiceName("gamebackend_service"),
		uptrace.WithServiceVersion("v1.0.0"),
	)
	defer uptrace.Shutdown(ctx)

	jwtService, err := service.NewJWTService(conf.KeyDir)
	helpers.Check(err, "jwt service")

	grpcServer, gwmux, err := NewServer(jwtService)
	helpers.Check(err, "create grpc server")

	lis, err := net.Listen("tcp", conf.GameBackend.Local.Address())
	helpers.Check(err, "listen")

	server := &http.Server{
		Addr:    conf.GameBackend.Local.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")

	err = server.Serve(lis)
	helpers.Check(err, "serve")
}

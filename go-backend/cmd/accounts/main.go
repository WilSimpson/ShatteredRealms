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
	Accounts config.Server `yaml:"accounts"`
	KeyDir   string        `yaml:"keyDir"`
	DBFile   string        `yaml:"dbFile"`
}

var (
	conf = &appConfig{
		Accounts: config.Server{
			Port:     8080,
			Host:     "",
			Mode:     "development",
			LogLevel: log.InfoLevel,
		},
		KeyDir: "/etc/sro/auth",
		DBFile: "/etc/sro/db.yaml",
	}
)

func init() {
	helpers.SetupLogs()
	config.SetupConfig(conf)
}

func main() {
	db, err := repository.ConnectFromFile(conf.DBFile)
	helpers.Check(err, "db connect from file")

	permissionRepository := repository.NewPermissionRepository(db)
	helpers.Check(permissionRepository.Migrate(), "permissions repo")
	roleRepository := repository.NewRoleRepository(db)
	helpers.Check(roleRepository.Migrate(), "role repo")
	userRepository := repository.NewUserRepository(db)
	helpers.Check(userRepository.Migrate(), "user repo")

	permissionService := service.NewPermissionService(permissionRepository)
	roleService := service.NewRoleService(roleRepository)
	userService := service.NewUserService(userRepository)
	jwtService, err := service.NewJWTService(conf.KeyDir)
	helpers.Check(err, "jwt service")

	grpcServer, gwmux, err := NewServer(userService, permissionService, roleService, jwtService)
	helpers.Check(err, "create grpc server")

	seedDatabaseIfNeeded(userService, permissionService, roleService, grpcServer.GetServiceInfo())

	lis, err := net.Listen("tcp", conf.Accounts.Address())
	helpers.Check(err, "listen")

	server := &http.Server{
		Addr:    conf.Accounts.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")

	err = server.Serve(lis)
	helpers.Check(err, "serve")
}

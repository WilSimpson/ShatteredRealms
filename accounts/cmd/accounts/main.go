package main

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/ShatteredRealms/Accounts/pkg/srv"
	"github.com/ShatteredRealms/GoUtils/pkg/helpers"
	utilRepository "github.com/ShatteredRealms/GoUtils/pkg/repository"
	utilService "github.com/ShatteredRealms/GoUtils/pkg/service"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

var (
	starterPermissions = []string{
		"/sro.accounts.AuthenticationService/Login",
		"/sro.accounts.AuthenticationService/Register",
	}
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

	c := &utilRepository.DBConnections{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("parsing db file: %v", err))
		os.Exit(1)
	}

	db, err := utilRepository.DBConnect(*c)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("db: %v", err))
		os.Exit(1)
	}

	permissionRepository := repository.NewPermissionRepository(db)
	if err := permissionRepository.Migrate(); err != nil {
		logger.Log(log.Error, fmt.Sprintf("permission repo: %v", err))
		os.Exit(1)
	}

	roleRepository := repository.NewRoleRepository(db)
	if err := roleRepository.Migrate(); err != nil {
		logger.Log(log.Error, fmt.Sprintf("role repo: %v", err))
		os.Exit(1)
	}

	permissionService := service.NewPermissionService(permissionRepository)
	roleService := service.NewRoleService(roleRepository)

	userRepository := repository.NewUserRepository(db)
	if err := userRepository.Migrate(); err != nil {
		logger.Log(log.Error, fmt.Sprintf("user repo: %v", err))
		os.Exit(1)
	}
	userService := service.NewUserService(userRepository)

	jwtService, err := utilService.NewJWTService(*config.KeyDir.Value)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("jwt service: %v", err))
		os.Exit(1)
	}
	grpcServer, gwmux, err := srv.NewServer(userService, permissionService, roleService, jwtService, logger, config)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("server creation: %v", err))
		os.Exit(1)
	}

	seedDatabaseIfNeeded(userService, permissionService, roleService, grpcServer.GetServiceInfo(), logger)

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

func seedDatabaseIfNeeded(
	userService service.UserService,
	permissionService service.PermissionService,
	roleService service.RoleService,
	servicesInfo map[string]grpc.ServiceInfo,
	logger log.LoggerService,
) {
	var err error
	superAdminRole := roleService.FindByName("Super Admin")
	if superAdminRole.Model.ID == 0 {
		// Create Super Admin role
		superAdminRole, err = roleService.Create(&model.Role{
			Name: "Super Admin",
		})
		if err != nil {
			logger.Logf(log.Error, "creating super admin: %v", err)
			return
		}
	}

	currentPermissions := createSetOfPermissions(
		permissionService.FindPermissionsForRoleID(superAdminRole.ID),
	)

	// Assign all permissions with Other set to true
	for packageService, serviceInfo := range servicesInfo {
		for _, methodInfo := range serviceInfo.Methods {
			permission := fmt.Sprintf("/%s/%s", packageService, methodInfo.Name)
			// Only add the permission if it doesn't exist already
			if _, ok := currentPermissions[permission]; !ok {
				err = permissionService.AddPermissionForRole(&model.RolePermission{
					RoleID:     superAdminRole.ID,
					Permission: permission,
					Other:      true,
				})

				if err != nil {
					logger.Logf(log.Error, "creating permission %s for super admin: %v", permission, err)
				}
			}
		}
	}

	if len(userService.FindAll()) > 0 {
		return
	}

	_, err = userService.Create(&model.User{
		FirstName: "Wil",
		LastName:  "Simpson",
		Username:  "unreal",
		Email:     "wil@forever.dev",
		Password:  "password",
		Roles:     []*model.Role{superAdminRole},
	})

	if err != nil {
		logger.Logf(log.Error, "creating super user: %v", err)
	}
}

func createSetOfPermissions(permissions []*model.RolePermission) map[string]struct{} {
	out := make(map[string]struct{}, len(permissions))

	for _, v := range permissions {
		out[v.Permission] = struct{}{}
	}

	return out
}

package main

import (
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
	utilRepository "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	utilService "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
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

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

func main() {

	file, err := ioutil.ReadFile(conf.DBFile)
	if err != nil {
		log.Error(fmt.Sprintf("reading db file: %v", err))
		os.Exit(1)
	}

	c := &utilRepository.DBConnections{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Error(fmt.Sprintf("parsing db file: %v", err))
		os.Exit(1)
	}

	db, err := utilRepository.DBConnect(*c)
	if err != nil {
		log.Error(fmt.Sprintf("db: %v", err))
		os.Exit(1)
	}

	permissionRepository := repository.NewPermissionRepository(db)
	if err := permissionRepository.Migrate(); err != nil {
		log.Error(fmt.Sprintf("permission repo: %v", err))
		os.Exit(1)
	}

	roleRepository := repository.NewRoleRepository(db)
	if err := roleRepository.Migrate(); err != nil {
		log.Error(fmt.Sprintf("role repo: %v", err))
		os.Exit(1)
	}

	permissionService := service.NewPermissionService(permissionRepository)
	roleService := service.NewRoleService(roleRepository)

	userRepository := repository.NewUserRepository(db)
	if err := userRepository.Migrate(); err != nil {
		log.Error(fmt.Sprintf("user repo: %v", err))
		os.Exit(1)
	}
	userService := service.NewUserService(userRepository)

	jwtService, err := utilService.NewJWTService(conf.KeyDir)
	if err != nil {
		log.Error(fmt.Sprintf("jwt service: %v", err))
		os.Exit(1)
	}
	grpcServer, gwmux, err := NewServer(userService, permissionService, roleService, jwtService)
	if err != nil {
		log.Error(fmt.Sprintf("server creation: %v", err))
		os.Exit(1)
	}

	seedDatabaseIfNeeded(userService, permissionService, roleService, grpcServer.GetServiceInfo())

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

func seedDatabaseIfNeeded(
	userService service.UserService,
	permissionService service.PermissionService,
	roleService service.RoleService,
	servicesInfo map[string]grpc.ServiceInfo,
) {
	var err error
	superAdminRole := roleService.FindByName("Super Admin")
	if superAdminRole.Model.ID == 0 {
		// Create Super Admin role
		superAdminRole, err = roleService.Create(&model.Role{
			Name: "Super Admin",
		})
		if err != nil {
			log.Error(fmt.Sprintf("creating super admin: %v", err))
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
					log.Error(fmt.Sprintf("creating permission %s for super admin: %v", permission, err))
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
		log.Error(fmt.Sprintf("creating super user: %v", err))
	}
}

func createSetOfPermissions(permissions []*model.RolePermission) map[string]struct{} {
	out := make(map[string]struct{}, len(permissions))

	for _, v := range permissions {
		out[v.Permission] = struct{}{}
	}

	return out
}

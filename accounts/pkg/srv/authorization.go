package srv

import (
	"context"
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/pb"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	charactersPb "github.com/ShatteredRealms/Characters/pkg/pb"
	"github.com/ShatteredRealms/GoUtils/pkg/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type authorizationServiceServer struct {
	pb.UnimplementedAuthorizationServiceServer
	userService       service.UserService
	permissionService service.PermissionService
	roleService       service.RoleService
	logger            log.LoggerService
	allPermissions    *pb.UserPermissions
	authInterceptor   *interceptor.AuthInterceptor
	userUpdates       chan uint64
	roleUpdates       chan uint64
}

func NewAuthorizationServiceServer(
	u service.UserService,
	permissionService service.PermissionService,
	roleService service.RoleService,
	logger log.LoggerService,
) *authorizationServiceServer {
	return &authorizationServiceServer{
		userService:       u,
		permissionService: permissionService,
		roleService:       roleService,
		logger:            logger,
		userUpdates:       make(chan uint64),
		roleUpdates:       make(chan uint64),
	}
}

func (s *authorizationServiceServer) AddAuthInterceptor(interceptor *interceptor.AuthInterceptor) {
	s.authInterceptor = interceptor
}

func (s *authorizationServiceServer) GetAuthorization(
	ctx context.Context,
	message *pb.IDMessage,
) (*pb.AuthorizationMessage, error) {
	user := s.userService.FindById(uint(message.Id))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	permissions := ConvertUserPermissions(s.permissionService.FindPermissionsForUserID(user.ID))
	roles := ConvertRolesWithoutPermissions(user.Roles)
	for i, role := range roles {
		roles[i].Permissions = ConvertRolePermissions(s.permissionService.FindPermissionsForRoleID(uint(role.Id)))
	}

	resp := &pb.AuthorizationMessage{
		UserId:      message.Id,
		Roles:       roles,
		Permissions: permissions,
	}

	return resp, nil
}

func (s *authorizationServiceServer) AddAuthorization(ctx context.Context, message *pb.AuthorizationMessage) (*emptypb.Empty, error) {
	user := s.userService.FindById(uint(message.UserId))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, v := range message.Permissions {
		err := s.permissionService.AddPermissionForUser(&model.UserPermission{
			UserID:     user.ID,
			Permission: v.Permission.Value,
			Other:      v.Other,
		})

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	for _, v := range message.Roles {
		err := s.userService.AddToRole(
			user,
			&model.Role{
				Model: gorm.Model{
					ID: uint(v.Id),
				},
			})

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	s.userUpdates <- message.UserId
	err := s.authInterceptor.ClearUserCache(uint(message.UserId))
	return &emptypb.Empty{}, err
}

func (s *authorizationServiceServer) RemoveAuthorization(ctx context.Context, message *pb.AuthorizationMessage) (*emptypb.Empty, error) {
	user := s.userService.FindById(uint(message.UserId))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, v := range message.Permissions {
		err := s.permissionService.RemPermissionForUser(&model.UserPermission{
			UserID:     user.ID,
			Permission: v.Permission.Value,
			Other:      v.Other,
		})

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	for _, v := range message.Roles {
		err := s.userService.RemFromRole(
			user,
			&model.Role{
				Model: gorm.Model{
					ID: uint(v.Id),
				},
			})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	s.userUpdates <- message.UserId
	err := s.authInterceptor.ClearUserCache(uint(message.UserId))
	return &emptypb.Empty{}, err
}

func (s *authorizationServiceServer) GetRoles(ctx context.Context, message *emptypb.Empty) (*pb.UserRoles, error) {
	resp := &pb.UserRoles{
		Roles: ConvertRolesNamesOnly(s.roleService.FindAll()),
	}
	for i, v := range resp.Roles {
		resp.Roles[i].Permissions = ConvertRolePermissions(s.permissionService.FindPermissionsForRoleID(uint(v.Id)))
	}

	return resp, nil
}

func (s *authorizationServiceServer) GetRole(ctx context.Context, message *pb.IDMessage) (*pb.UserRole, error) {
	resp := ConvertRoleNameOnly(s.roleService.FindById(uint(message.Id)))
	resp.Permissions = ConvertRolePermissions(s.permissionService.FindPermissionsForRoleID(uint(message.Id)))

	return resp, nil
}

func (s *authorizationServiceServer) CreateRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	_, err := s.roleService.Create(&model.Role{
		Name: message.Name.Value,
	})

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *authorizationServiceServer) EditRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	if message.Name != nil {
		err := s.roleService.Update(&model.Role{
			Model: gorm.Model{
				ID: uint(message.Id),
			},
			Name: message.Name.Value,
		})
		if err != nil {
			return nil, err
		}
	}
	if message.Permissions != nil {
		newPermissions := make([]*model.RolePermission, len(message.Permissions))
		for i, permission := range message.Permissions {
			newPermissions[i] = &model.RolePermission{
				RoleID:     uint(message.Id),
				Permission: permission.Permission.Value,
				Other:      permission.Other,
			}
		}

		err := s.permissionService.ResetPermissionsForRole(uint(message.Id), newPermissions)
		if err != nil {
			return nil, err
		}
	}

	s.roleUpdates <- message.Id
	err := s.authInterceptor.ClearCache()
	return &emptypb.Empty{}, err
}

func (s *authorizationServiceServer) DeleteRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	err := s.roleService.Delete(&model.Role{
		Model: gorm.Model{
			ID: uint(message.Id),
		},
	})

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	s.roleUpdates <- message.Id
	err = s.authInterceptor.ClearCache()
	return &emptypb.Empty{}, err
}

func (s *authorizationServiceServer) GetAllPermissions(ctx context.Context, message *emptypb.Empty) (*pb.UserPermissions, error) {
	return s.allPermissions, nil
}

func (s *authorizationServiceServer) SetupAllPermissions(accountsServicesInfo map[string]grpc.ServiceInfo) {
	accountsPermissions := setupPermissions(accountsServicesInfo)
	charactersPermissions := setupPermissions(getCharactersServiceInfo())

	s.allPermissions = &pb.UserPermissions{Permissions: accountsPermissions}
	s.allPermissions.Permissions = append(s.allPermissions.Permissions, charactersPermissions...)
}

func (s *authorizationServiceServer) SubscribeUserUpdates(message *emptypb.Empty, stream pb.AuthorizationService_SubscribeUserUpdatesServer) error {
	for {
		select {
		case <-stream.Context().Done():
			s.logger.Log(log.Debug, "User subscribe context closed")
			return nil
		case userId := <-s.userUpdates:
			s.logger.Log(log.Debug, "Sending update")
			err := stream.Send(&pb.IDMessage{Id: userId})
			s.logger.Log(log.Debug, "Broadcast role update")
			if err != nil {
				return err
			}
		}
	}
}
func (s *authorizationServiceServer) SubscribeRoleUpdates(message *emptypb.Empty, stream pb.AuthorizationService_SubscribeRoleUpdatesServer) error {
	for {
		select {
		case <-stream.Context().Done():
			s.logger.Log(log.Debug, "Role subscribe context closed")
			return nil
		case userId := <-s.roleUpdates:
			s.logger.Log(log.Debug, "Sending update")
			err := stream.Send(&pb.IDMessage{Id: userId})
			s.logger.Log(log.Debug, "Broadcast role update")
			if err != nil {
				return err
			}
		}
	}
}

func getCharactersServiceInfo() map[string]grpc.ServiceInfo {
	grpcServer := grpc.NewServer()
	charactersPb.RegisterCharactersServiceServer(grpcServer, charactersPb.UnimplementedCharactersServiceServer{})
	return grpcServer.GetServiceInfo()
}

func setupPermissions(serviceInfos map[string]grpc.ServiceInfo) []*pb.UserPermission {

	count := 0
	for _, serviceInfo := range serviceInfos {
		count += len(serviceInfo.Methods)
	}

	methods := make([]*pb.UserPermission, count)
	index := 0
	for serviceName, serviceInfo := range serviceInfos {
		for _, method := range serviceInfo.Methods {
			methods[index] = &pb.UserPermission{
				Permission: &wrapperspb.StringValue{Value: fmt.Sprintf("/%s/%s", serviceName, method.Name)},
			}
			index++
		}
	}

	return methods
}

package srv

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	utilService "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// @TODO(wil): Change all errors to variables

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService       service.UserService
	permissionService service.PermissionService
	jwtService        utilService.JWTService
}

func NewUserServiceServer(
	u service.UserService,
	p service.PermissionService,
	j utilService.JWTService,
) *userServiceServer {
	return &userServiceServer{
		userService:       u,
		permissionService: p,
		jwtService:        j,
	}
}

func (s *userServiceServer) GetAll(
	ctx context.Context,
	message *emptypb.Empty,
) (*pb.GetAllUsersResponse, error) {
	users := s.userService.FindAll()
	resp := &pb.GetAllUsersResponse{
		Users: []*pb.UserMessage{},
	}
	for _, u := range users {
		resp.Users = append(resp.Users, &pb.UserMessage{
			Id:                 uint64(u.ID),
			Email:              u.Email,
			Username:           u.Username,
			Roles:              ConvertRolesNamesOnly(u.Roles),
			CreatedAt:          u.CreatedAt.String(),
			BannedAt:           u.BannedAtString(),
			CurrentCharacterId: uint64(u.CurrentCharacterId),
		})
	}

	return resp, nil
}

func (s *userServiceServer) Get(
	ctx context.Context,
	message *pb.GetUserMessage,
) (*pb.GetUserResponse, error) {
	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		Id:                 uint64(user.ID),
		Email:              user.Email,
		Username:           user.Username,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Roles:              ConvertRolesNamesOnly(user.Roles),
		Permissions:        ConvertUserPermissions(s.permissionService.FindPermissionsForUserID(user.ID)),
		CreatedAt:          user.CreatedAt.String(),
		BannedAt:           user.BannedAtString(),
		CurrentCharacterId: uint64(user.CurrentCharacterId),
	}, nil
}

func (s *userServiceServer) Edit(
	ctx context.Context,
	message *pb.UserDetails,
) (*emptypb.Empty, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, uint(message.UserId))
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	newUserData := model.User{
		FirstName: message.FirstName,
		LastName:  message.LastName,
		Username:  message.Username,
		Email:     message.Email,
		Password:  message.Password,
	}

	err = user.UpdateInfo(newUserData)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.userService.Save(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *userServiceServer) Ban(ctx context.Context, message *pb.GetUserMessage) (*emptypb.Empty, error) {
	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err := s.userService.Ban(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to ban user: %v", err.Error())
	}

	return &emptypb.Empty{}, nil
}
func (s *userServiceServer) UnBan(ctx context.Context, message *pb.GetUserMessage) (*emptypb.Empty, error) {
	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err := s.userService.UnBan(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to unban user: %v", err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *userServiceServer) GetStatus(ctx context.Context, message *pb.GetUserMessage) (*pb.StatusResponse, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, uint(message.UserId))
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.StatusResponse{CharacterId: uint64(user.CurrentCharacterId)}, nil
}

func (s *userServiceServer) SetStatus(ctx context.Context, message *pb.StatusRequest) (*emptypb.Empty, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, uint(message.UserId))
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	user.CurrentCharacterId = uint(message.CharacterId)
	_, err = s.userService.Save(user)
	if err != nil {
		if !user.Exists() {
			return nil, status.Error(codes.Internal, "unable to update user")
		}
	}

	return &emptypb.Empty{}, nil
}

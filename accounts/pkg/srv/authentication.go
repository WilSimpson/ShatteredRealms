package srv

import (
	"context"
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	accountModel "github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/pb"
	accountService "github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/ShatteredRealms/GoUtils/pkg/service"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type authenticationServiceServer struct {
	pb.UnimplementedAuthenticationServiceServer
	userService       accountService.UserService
	permissionService accountService.PermissionService
	jwtService        service.JWTService
	logger            log.LoggerService
}

func NewAuthenticationServiceServer(
	u accountService.UserService,
	jwt service.JWTService,
	permissionService accountService.PermissionService,
	logger log.LoggerService,
) *authenticationServiceServer {
	return &authenticationServiceServer{
		userService:       u,
		permissionService: permissionService,
		jwtService:        jwt,
		logger:            logger,
	}
}

func (s *authenticationServiceServer) Register(
	ctx context.Context,
	message *pb.RegisterAccountMessage,
) (*emptypb.Empty, error) {
	user := &accountModel.User{
		FirstName: message.FirstName,
		LastName:  message.LastName,
		Username:  message.Username,
		Email:     message.Email,
		Password:  message.Password,
	}

	user, err := s.userService.Create(user)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	s.logger.LogRegisterRequest()

	return &emptypb.Empty{}, nil
}

func (s *authenticationServiceServer) Login(
	ctx context.Context,
	message *pb.LoginMessage,
) (*pb.LoginResponse, error) {
	if message.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email cannot be empty")
	}

	if message.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password cannot be empty")
	}

	user := s.userService.FindByEmail(message.Email)
	if !user.Exists() || user.Login(message.Password) != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid username or password")
	}

	token, err := s.tokenForUser(user)
	if err != nil {
		s.logger.Log(log.Error, fmt.Sprintf("error signing jwt: %v", err))
		return nil, status.Error(codes.Internal, "Error signing validation token")
	}

	s.logger.LogLoginRequest()

	return &pb.LoginResponse{
		Token:     token,
		Id:        uint64(user.ID),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.String(),
		Roles:     ConvertRolesWithoutPermissions(user.Roles),
		BannedAt:  user.BannedAtString(),
	}, nil
}

func (s *authenticationServiceServer) tokenForUser(u *accountModel.User) (t string, err error) {
	claims := jwt.MapClaims{
		"sub":                u.ID,
		"preferred_username": u.Username,
		//"given_name":  u.FirstName,
		//"family_name": u.LastName,
		//"email":       u.Email,
	}

	t, err = s.jwtService.Create(time.Hour, "shatteredrealmsonline.com/accounts/v1", claims)
	return t, err
}

package srv

import (
    "context"
    accountModel "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
    "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
    "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
    accountService "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
    "github.com/golang-jwt/jwt"
    log "github.com/sirupsen/logrus"
    otelcodes "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/emptypb"
    "time"
)

var (
    ErrorEmailEmpty         = status.Error(codes.InvalidArgument, "Email cannot be empty")
    ErrorInvalidEmailOrPass = status.Error(codes.Unauthenticated, "Invalid username or password")
    ErrorCreatingToken      = status.Error(codes.Internal, "Error creating token")
)

type authenticationServiceServer struct {
    pb.UnimplementedAuthenticationServiceServer
    userService       accountService.UserService
    permissionService accountService.PermissionService
    jwtService        service.JWTService
}

func NewAuthenticationServiceServer(
    u accountService.UserService,
    jwt service.JWTService,
    permissionService accountService.PermissionService,
) *authenticationServiceServer {
    return &authenticationServiceServer{
        userService:       u,
        permissionService: permissionService,
        jwtService:        jwt,
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

    user, err := s.userService.Create(ctx, user)
    if err != nil {
        span := trace.SpanFromContext(ctx)
        span.RecordError(err)
        span.SetStatus(otelcodes.Error, "creating user")
        return nil, status.Error(codes.FailedPrecondition, err.Error())
    }

    return &emptypb.Empty{}, nil
}

func (s *authenticationServiceServer) Login(
    ctx context.Context,
    message *pb.LoginMessage,
) (*pb.LoginResponse, error) {
    span := trace.SpanFromContext(ctx)

    if message.Email == "" {
        span.RecordError(ErrorEmailEmpty)
        span.SetStatus(otelcodes.Error, "no email")
        return nil, ErrorEmailEmpty
    }

    if message.Password == "" {
        return nil, status.Error(codes.InvalidArgument, "Password cannot be empty")
    }

    user := s.userService.FindByEmail(ctx, message.Email)
    if !user.Exists() {
        span.RecordError(ErrorInvalidEmailOrPass)
        span.SetStatus(otelcodes.Error, "email not used")
    }

    if user.Login(message.Password) != nil {
        span.RecordError(ErrorInvalidEmailOrPass)
        span.SetStatus(otelcodes.Error, "invalid password")
        return nil, ErrorInvalidEmailOrPass
    }

    token, err := s.tokenForUser(ctx, user)
    if err != nil {
        log.WithContext(ctx).Errorf("error signing jwt: %v", err)
        span.RecordError(ErrorCreatingToken)
        span.SetStatus(otelcodes.Error, "creating token")
        return nil, ErrorCreatingToken
    }

    return &pb.LoginResponse{
        Token:     token,
        Id:        uint64(user.ID),
        Email:     user.Email,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Username:  user.Username,
        CreatedAt: user.CreatedAt.String(),
        Roles:     user.Roles.ToPB().Roles,
        BannedAt:  user.BannedAtWrapper(),
    }, nil
}

func (s *authenticationServiceServer) tokenForUser(ctx context.Context, u *accountModel.User) (t string, err error) {
    claims := jwt.MapClaims{
        "sub":                u.ID,
        "preferred_username": u.Username,
        //"given_name":  u.FirstName,
        //"family_name": u.LastName,
        //"email":       u.Email,
    }

    t, err = s.jwtService.Create(ctx, time.Hour, "sro.com/accounts/v1", claims)
    return t, err
}

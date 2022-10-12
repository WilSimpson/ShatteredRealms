// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: authorization.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthorizationServiceClient is the client API for AuthorizationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthorizationServiceClient interface {
	GetAuthorization(ctx context.Context, in *IDMessage, opts ...grpc.CallOption) (*AuthorizationMessage, error)
	AddAuthorization(ctx context.Context, in *AuthorizationMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RemoveAuthorization(ctx context.Context, in *AuthorizationMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetRoles(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserRoles, error)
	GetRole(ctx context.Context, in *IDMessage, opts ...grpc.CallOption) (*UserRole, error)
	CreateRole(ctx context.Context, in *UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error)
	EditRole(ctx context.Context, in *UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteRole(ctx context.Context, in *UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetAllPermissions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserPermissions, error)
	SubscribeUserUpdates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (AuthorizationService_SubscribeUserUpdatesClient, error)
	SubscribeRoleUpdates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (AuthorizationService_SubscribeRoleUpdatesClient, error)
}

type authorizationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthorizationServiceClient(cc grpc.ClientConnInterface) AuthorizationServiceClient {
	return &authorizationServiceClient{cc}
}

func (c *authorizationServiceClient) GetAuthorization(ctx context.Context, in *IDMessage, opts ...grpc.CallOption) (*AuthorizationMessage, error) {
	out := new(AuthorizationMessage)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/GetAuthorization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) AddAuthorization(ctx context.Context, in *AuthorizationMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/AddAuthorization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) RemoveAuthorization(ctx context.Context, in *AuthorizationMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/RemoveAuthorization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) GetRoles(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserRoles, error) {
	out := new(UserRoles)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/GetRoles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) GetRole(ctx context.Context, in *IDMessage, opts ...grpc.CallOption) (*UserRole, error) {
	out := new(UserRole)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/GetRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) CreateRole(ctx context.Context, in *UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/CreateRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) EditRole(ctx context.Context, in *UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/EditRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) DeleteRole(ctx context.Context, in *UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/DeleteRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) GetAllPermissions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserPermissions, error) {
	out := new(UserPermissions)
	err := c.cc.Invoke(ctx, "/sro.accounts.AuthorizationService/GetAllPermissions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorizationServiceClient) SubscribeUserUpdates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (AuthorizationService_SubscribeUserUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &AuthorizationService_ServiceDesc.Streams[0], "/sro.accounts.AuthorizationService/SubscribeUserUpdates", opts...)
	if err != nil {
		return nil, err
	}
	x := &authorizationServiceSubscribeUserUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AuthorizationService_SubscribeUserUpdatesClient interface {
	Recv() (*IDMessage, error)
	grpc.ClientStream
}

type authorizationServiceSubscribeUserUpdatesClient struct {
	grpc.ClientStream
}

func (x *authorizationServiceSubscribeUserUpdatesClient) Recv() (*IDMessage, error) {
	m := new(IDMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *authorizationServiceClient) SubscribeRoleUpdates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (AuthorizationService_SubscribeRoleUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &AuthorizationService_ServiceDesc.Streams[1], "/sro.accounts.AuthorizationService/SubscribeRoleUpdates", opts...)
	if err != nil {
		return nil, err
	}
	x := &authorizationServiceSubscribeRoleUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AuthorizationService_SubscribeRoleUpdatesClient interface {
	Recv() (*IDMessage, error)
	grpc.ClientStream
}

type authorizationServiceSubscribeRoleUpdatesClient struct {
	grpc.ClientStream
}

func (x *authorizationServiceSubscribeRoleUpdatesClient) Recv() (*IDMessage, error) {
	m := new(IDMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AuthorizationServiceServer is the server API for AuthorizationService service.
// All implementations must embed UnimplementedAuthorizationServiceServer
// for forward compatibility
type AuthorizationServiceServer interface {
	GetAuthorization(context.Context, *IDMessage) (*AuthorizationMessage, error)
	AddAuthorization(context.Context, *AuthorizationMessage) (*emptypb.Empty, error)
	RemoveAuthorization(context.Context, *AuthorizationMessage) (*emptypb.Empty, error)
	GetRoles(context.Context, *emptypb.Empty) (*UserRoles, error)
	GetRole(context.Context, *IDMessage) (*UserRole, error)
	CreateRole(context.Context, *UserRole) (*emptypb.Empty, error)
	EditRole(context.Context, *UserRole) (*emptypb.Empty, error)
	DeleteRole(context.Context, *UserRole) (*emptypb.Empty, error)
	GetAllPermissions(context.Context, *emptypb.Empty) (*UserPermissions, error)
	SubscribeUserUpdates(*emptypb.Empty, AuthorizationService_SubscribeUserUpdatesServer) error
	SubscribeRoleUpdates(*emptypb.Empty, AuthorizationService_SubscribeRoleUpdatesServer) error
	mustEmbedUnimplementedAuthorizationServiceServer()
}

// UnimplementedAuthorizationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthorizationServiceServer struct {
}

func (UnimplementedAuthorizationServiceServer) GetAuthorization(context.Context, *IDMessage) (*AuthorizationMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthorization not implemented")
}
func (UnimplementedAuthorizationServiceServer) AddAuthorization(context.Context, *AuthorizationMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAuthorization not implemented")
}
func (UnimplementedAuthorizationServiceServer) RemoveAuthorization(context.Context, *AuthorizationMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveAuthorization not implemented")
}
func (UnimplementedAuthorizationServiceServer) GetRoles(context.Context, *emptypb.Empty) (*UserRoles, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoles not implemented")
}
func (UnimplementedAuthorizationServiceServer) GetRole(context.Context, *IDMessage) (*UserRole, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRole not implemented")
}
func (UnimplementedAuthorizationServiceServer) CreateRole(context.Context, *UserRole) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRole not implemented")
}
func (UnimplementedAuthorizationServiceServer) EditRole(context.Context, *UserRole) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditRole not implemented")
}
func (UnimplementedAuthorizationServiceServer) DeleteRole(context.Context, *UserRole) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRole not implemented")
}
func (UnimplementedAuthorizationServiceServer) GetAllPermissions(context.Context, *emptypb.Empty) (*UserPermissions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPermissions not implemented")
}
func (UnimplementedAuthorizationServiceServer) SubscribeUserUpdates(*emptypb.Empty, AuthorizationService_SubscribeUserUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeUserUpdates not implemented")
}
func (UnimplementedAuthorizationServiceServer) SubscribeRoleUpdates(*emptypb.Empty, AuthorizationService_SubscribeRoleUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeRoleUpdates not implemented")
}
func (UnimplementedAuthorizationServiceServer) mustEmbedUnimplementedAuthorizationServiceServer() {}

// UnsafeAuthorizationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthorizationServiceServer will
// result in compilation errors.
type UnsafeAuthorizationServiceServer interface {
	mustEmbedUnimplementedAuthorizationServiceServer()
}

func RegisterAuthorizationServiceServer(s grpc.ServiceRegistrar, srv AuthorizationServiceServer) {
	s.RegisterService(&AuthorizationService_ServiceDesc, srv)
}

func _AuthorizationService_GetAuthorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).GetAuthorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/GetAuthorization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).GetAuthorization(ctx, req.(*IDMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_AddAuthorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizationMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).AddAuthorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/AddAuthorization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).AddAuthorization(ctx, req.(*AuthorizationMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_RemoveAuthorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizationMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).RemoveAuthorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/RemoveAuthorization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).RemoveAuthorization(ctx, req.(*AuthorizationMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_GetRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).GetRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/GetRoles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).GetRoles(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_GetRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).GetRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/GetRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).GetRole(ctx, req.(*IDMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_CreateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRole)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).CreateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/CreateRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).CreateRole(ctx, req.(*UserRole))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_EditRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRole)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).EditRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/EditRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).EditRole(ctx, req.(*UserRole))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_DeleteRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRole)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).DeleteRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/DeleteRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).DeleteRole(ctx, req.(*UserRole))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_GetAllPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServiceServer).GetAllPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.accounts.AuthorizationService/GetAllPermissions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServiceServer).GetAllPermissions(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorizationService_SubscribeUserUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AuthorizationServiceServer).SubscribeUserUpdates(m, &authorizationServiceSubscribeUserUpdatesServer{stream})
}

type AuthorizationService_SubscribeUserUpdatesServer interface {
	Send(*IDMessage) error
	grpc.ServerStream
}

type authorizationServiceSubscribeUserUpdatesServer struct {
	grpc.ServerStream
}

func (x *authorizationServiceSubscribeUserUpdatesServer) Send(m *IDMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _AuthorizationService_SubscribeRoleUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AuthorizationServiceServer).SubscribeRoleUpdates(m, &authorizationServiceSubscribeRoleUpdatesServer{stream})
}

type AuthorizationService_SubscribeRoleUpdatesServer interface {
	Send(*IDMessage) error
	grpc.ServerStream
}

type authorizationServiceSubscribeRoleUpdatesServer struct {
	grpc.ServerStream
}

func (x *authorizationServiceSubscribeRoleUpdatesServer) Send(m *IDMessage) error {
	return x.ServerStream.SendMsg(m)
}

// AuthorizationService_ServiceDesc is the grpc.ServiceDesc for AuthorizationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthorizationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sro.accounts.AuthorizationService",
	HandlerType: (*AuthorizationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAuthorization",
			Handler:    _AuthorizationService_GetAuthorization_Handler,
		},
		{
			MethodName: "AddAuthorization",
			Handler:    _AuthorizationService_AddAuthorization_Handler,
		},
		{
			MethodName: "RemoveAuthorization",
			Handler:    _AuthorizationService_RemoveAuthorization_Handler,
		},
		{
			MethodName: "GetRoles",
			Handler:    _AuthorizationService_GetRoles_Handler,
		},
		{
			MethodName: "GetRole",
			Handler:    _AuthorizationService_GetRole_Handler,
		},
		{
			MethodName: "CreateRole",
			Handler:    _AuthorizationService_CreateRole_Handler,
		},
		{
			MethodName: "EditRole",
			Handler:    _AuthorizationService_EditRole_Handler,
		},
		{
			MethodName: "DeleteRole",
			Handler:    _AuthorizationService_DeleteRole_Handler,
		},
		{
			MethodName: "GetAllPermissions",
			Handler:    _AuthorizationService_GetAllPermissions_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeUserUpdates",
			Handler:       _AuthorizationService_SubscribeUserUpdates_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeRoleUpdates",
			Handler:       _AuthorizationService_SubscribeRoleUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "authorization.proto",
}

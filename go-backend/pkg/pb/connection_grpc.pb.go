// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: connection.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ConnectionServiceClient is the client API for ConnectionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectionServiceClient interface {
	ConnectGameServer(ctx context.Context, in *ConnectGameServerRequest, opts ...grpc.CallOption) (*ConnectGameServerResponse, error)
}

type connectionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectionServiceClient(cc grpc.ClientConnInterface) ConnectionServiceClient {
	return &connectionServiceClient{cc}
}

func (c *connectionServiceClient) ConnectGameServer(ctx context.Context, in *ConnectGameServerRequest, opts ...grpc.CallOption) (*ConnectGameServerResponse, error) {
	out := new(ConnectGameServerResponse)
	err := c.cc.Invoke(ctx, "/sro.gamebackend.ConnectionService/ConnectGameServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectionServiceServer is the server API for ConnectionService service.
// All implementations must embed UnimplementedConnectionServiceServer
// for forward compatibility
type ConnectionServiceServer interface {
	ConnectGameServer(context.Context, *ConnectGameServerRequest) (*ConnectGameServerResponse, error)
	mustEmbedUnimplementedConnectionServiceServer()
}

// UnimplementedConnectionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConnectionServiceServer struct {
}

func (UnimplementedConnectionServiceServer) ConnectGameServer(context.Context, *ConnectGameServerRequest) (*ConnectGameServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectGameServer not implemented")
}
func (UnimplementedConnectionServiceServer) mustEmbedUnimplementedConnectionServiceServer() {}

// UnsafeConnectionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectionServiceServer will
// result in compilation errors.
type UnsafeConnectionServiceServer interface {
	mustEmbedUnimplementedConnectionServiceServer()
}

func RegisterConnectionServiceServer(s grpc.ServiceRegistrar, srv ConnectionServiceServer) {
	s.RegisterService(&ConnectionService_ServiceDesc, srv)
}

func _ConnectionService_ConnectGameServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectGameServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServiceServer).ConnectGameServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.gamebackend.ConnectionService/ConnectGameServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServiceServer).ConnectGameServer(ctx, req.(*ConnectGameServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConnectionService_ServiceDesc is the grpc.ServiceDesc for ConnectionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConnectionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sro.gamebackend.ConnectionService",
	HandlerType: (*ConnectionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnectGameServer",
			Handler:    _ConnectionService_ConnectGameServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "connection.proto",
}

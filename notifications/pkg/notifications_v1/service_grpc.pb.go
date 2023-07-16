// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: service.proto

package notifications_v1

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

const (
	Notifications_GetHistoryWithPeriod_FullMethodName = "/notifications.Notifications/GetHistoryWithPeriod"
)

// NotificationsClient is the client API for Notifications service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationsClient interface {
	GetHistoryWithPeriod(ctx context.Context, in *GetHistoryWithPeriodRequest, opts ...grpc.CallOption) (*GetHistoryWithPeriodResponse, error)
}

type notificationsClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationsClient(cc grpc.ClientConnInterface) NotificationsClient {
	return &notificationsClient{cc}
}

func (c *notificationsClient) GetHistoryWithPeriod(ctx context.Context, in *GetHistoryWithPeriodRequest, opts ...grpc.CallOption) (*GetHistoryWithPeriodResponse, error) {
	out := new(GetHistoryWithPeriodResponse)
	err := c.cc.Invoke(ctx, Notifications_GetHistoryWithPeriod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationsServer is the server API for Notifications service.
// All implementations must embed UnimplementedNotificationsServer
// for forward compatibility
type NotificationsServer interface {
	GetHistoryWithPeriod(context.Context, *GetHistoryWithPeriodRequest) (*GetHistoryWithPeriodResponse, error)
	mustEmbedUnimplementedNotificationsServer()
}

// UnimplementedNotificationsServer must be embedded to have forward compatible implementations.
type UnimplementedNotificationsServer struct {
}

func (UnimplementedNotificationsServer) GetHistoryWithPeriod(context.Context, *GetHistoryWithPeriodRequest) (*GetHistoryWithPeriodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHistoryWithPeriod not implemented")
}
func (UnimplementedNotificationsServer) mustEmbedUnimplementedNotificationsServer() {}

// UnsafeNotificationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationsServer will
// result in compilation errors.
type UnsafeNotificationsServer interface {
	mustEmbedUnimplementedNotificationsServer()
}

func RegisterNotificationsServer(s grpc.ServiceRegistrar, srv NotificationsServer) {
	s.RegisterService(&Notifications_ServiceDesc, srv)
}

func _Notifications_GetHistoryWithPeriod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHistoryWithPeriodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationsServer).GetHistoryWithPeriod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notifications_GetHistoryWithPeriod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationsServer).GetHistoryWithPeriod(ctx, req.(*GetHistoryWithPeriodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Notifications_ServiceDesc is the grpc.ServiceDesc for Notifications service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Notifications_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notifications.Notifications",
	HandlerType: (*NotificationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHistoryWithPeriod",
			Handler:    _Notifications_GetHistoryWithPeriod_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

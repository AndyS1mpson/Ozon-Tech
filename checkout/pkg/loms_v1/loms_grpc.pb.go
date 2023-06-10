// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: loms.proto

package loms_v1

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
	Loms_Stocks_FullMethodName      = "/loms.Loms/Stocks"
	Loms_CreateOrder_FullMethodName = "/loms.Loms/CreateOrder"
)

// LomsClient is the client API for Loms service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LomsClient interface {
	Stocks(ctx context.Context, in *StocksRequest, opts ...grpc.CallOption) (*StocksResponse, error)
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
}

type lomsClient struct {
	cc grpc.ClientConnInterface
}

func NewLomsClient(cc grpc.ClientConnInterface) LomsClient {
	return &lomsClient{cc}
}

func (c *lomsClient) Stocks(ctx context.Context, in *StocksRequest, opts ...grpc.CallOption) (*StocksResponse, error) {
	out := new(StocksResponse)
	err := c.cc.Invoke(ctx, Loms_Stocks_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, Loms_CreateOrder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LomsServer is the server API for Loms service.
// All implementations must embed UnimplementedLomsServer
// for forward compatibility
type LomsServer interface {
	Stocks(context.Context, *StocksRequest) (*StocksResponse, error)
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	mustEmbedUnimplementedLomsServer()
}

// UnimplementedLomsServer must be embedded to have forward compatible implementations.
type UnimplementedLomsServer struct {
}

func (UnimplementedLomsServer) Stocks(context.Context, *StocksRequest) (*StocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stocks not implemented")
}
func (UnimplementedLomsServer) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedLomsServer) mustEmbedUnimplementedLomsServer() {}

// UnsafeLomsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LomsServer will
// result in compilation errors.
type UnsafeLomsServer interface {
	mustEmbedUnimplementedLomsServer()
}

func RegisterLomsServer(s grpc.ServiceRegistrar, srv LomsServer) {
	s.RegisterService(&Loms_ServiceDesc, srv)
}

func _Loms_Stocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).Stocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Loms_Stocks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).Stocks(ctx, req.(*StocksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Loms_CreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Loms_ServiceDesc is the grpc.ServiceDesc for Loms service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Loms_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "loms.Loms",
	HandlerType: (*LomsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Stocks",
			Handler:    _Loms_Stocks_Handler,
		},
		{
			MethodName: "CreateOrder",
			Handler:    _Loms_CreateOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "loms.proto",
}
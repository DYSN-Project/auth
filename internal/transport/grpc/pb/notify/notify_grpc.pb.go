// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: notify.proto

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

// NotifyClient is the client API for Notify pkg.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotifyClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ConfirmRegister(ctx context.Context, in *EmailWithCode, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RecoveryPassword(ctx context.Context, in *EmailWithCode, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type notifyClient struct {
	cc grpc.ClientConnInterface
}

func NewNotifyClient(cc grpc.ClientConnInterface) NotifyClient {
	return &notifyClient{cc}
}

func (c *notifyClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/notify.Notify/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notifyClient) ConfirmRegister(ctx context.Context, in *EmailWithCode, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/notify.Notify/ConfirmRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notifyClient) RecoveryPassword(ctx context.Context, in *EmailWithCode, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/notify.Notify/RecoveryPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotifyServer is the server API for Notify pkg.
// All implementations must embed UnimplementedNotifyServer
// for forward compatibility
type NotifyServer interface {
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	ConfirmRegister(context.Context, *EmailWithCode) (*emptypb.Empty, error)
	RecoveryPassword(context.Context, *EmailWithCode) (*emptypb.Empty, error)
	mustEmbedUnimplementedNotifyServer()
}

// UnimplementedNotifyServer must be embedded to have forward compatible implementations.
type UnimplementedNotifyServer struct {
}

func (UnimplementedNotifyServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedNotifyServer) ConfirmRegister(context.Context, *EmailWithCode) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmRegister not implemented")
}
func (UnimplementedNotifyServer) RecoveryPassword(context.Context, *EmailWithCode) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoveryPassword not implemented")
}
func (UnimplementedNotifyServer) mustEmbedUnimplementedNotifyServer() {}

// UnsafeNotifyServer may be embedded to opt out of forward compatibility for this pkg.
// Use of this interface is not recommended, as added methods to NotifyServer will
// result in compilation errors.
type UnsafeNotifyServer interface {
	mustEmbedUnimplementedNotifyServer()
}

func RegisterNotifyServer(s grpc.ServiceRegistrar, srv NotifyServer) {
	s.RegisterService(&Notify_ServiceDesc, srv)
}

func _Notify_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotifyServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notify.Notify/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotifyServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notify_ConfirmRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailWithCode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotifyServer).ConfirmRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notify.Notify/ConfirmRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotifyServer).ConfirmRegister(ctx, req.(*EmailWithCode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notify_RecoveryPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailWithCode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotifyServer).RecoveryPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notify.Notify/RecoveryPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotifyServer).RecoveryPassword(ctx, req.(*EmailWithCode))
	}
	return interceptor(ctx, in, info, handler)
}

// Notify_ServiceDesc is the grpc.ServiceDesc for Notify pkg.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Notify_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notify.Notify",
	HandlerType: (*NotifyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Notify_Ping_Handler,
		},
		{
			MethodName: "ConfirmRegister",
			Handler:    _Notify_ConfirmRegister_Handler,
		},
		{
			MethodName: "RecoveryPassword",
			Handler:    _Notify_RecoveryPassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notify.proto",
}

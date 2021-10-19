// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// BookkeeperClient is the client API for Bookkeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookkeeperClient interface {
	Load(ctx context.Context, in *Ledger, opts ...grpc.CallOption) (*Ledger, error)
	Save(ctx context.Context, in *Ledger, opts ...grpc.CallOption) (*Ledger, error)
}

type bookkeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewBookkeeperClient(cc grpc.ClientConnInterface) BookkeeperClient {
	return &bookkeeperClient{cc}
}

func (c *bookkeeperClient) Load(ctx context.Context, in *Ledger, opts ...grpc.CallOption) (*Ledger, error) {
	out := new(Ledger)
	err := c.cc.Invoke(ctx, "/pb.Bookkeeper/Load", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookkeeperClient) Save(ctx context.Context, in *Ledger, opts ...grpc.CallOption) (*Ledger, error) {
	out := new(Ledger)
	err := c.cc.Invoke(ctx, "/pb.Bookkeeper/Save", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookkeeperServer is the server API for Bookkeeper service.
// All implementations must embed UnimplementedBookkeeperServer
// for forward compatibility
type BookkeeperServer interface {
	Load(context.Context, *Ledger) (*Ledger, error)
	Save(context.Context, *Ledger) (*Ledger, error)
	mustEmbedUnimplementedBookkeeperServer()
}

// UnimplementedBookkeeperServer must be embedded to have forward compatible implementations.
type UnimplementedBookkeeperServer struct {
}

func (UnimplementedBookkeeperServer) Load(context.Context, *Ledger) (*Ledger, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Load not implemented")
}
func (UnimplementedBookkeeperServer) Save(context.Context, *Ledger) (*Ledger, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedBookkeeperServer) mustEmbedUnimplementedBookkeeperServer() {}

// UnsafeBookkeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookkeeperServer will
// result in compilation errors.
type UnsafeBookkeeperServer interface {
	mustEmbedUnimplementedBookkeeperServer()
}

func RegisterBookkeeperServer(s grpc.ServiceRegistrar, srv BookkeeperServer) {
	s.RegisterService(&Bookkeeper_ServiceDesc, srv)
}

func _Bookkeeper_Load_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ledger)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookkeeperServer).Load(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Bookkeeper/Load",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookkeeperServer).Load(ctx, req.(*Ledger))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bookkeeper_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ledger)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookkeeperServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Bookkeeper/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookkeeperServer).Save(ctx, req.(*Ledger))
	}
	return interceptor(ctx, in, info, handler)
}

// Bookkeeper_ServiceDesc is the grpc.ServiceDesc for Bookkeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bookkeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Bookkeeper",
	HandlerType: (*BookkeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Load",
			Handler:    _Bookkeeper_Load_Handler,
		},
		{
			MethodName: "Save",
			Handler:    _Bookkeeper_Save_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bookkeeper.proto",
}

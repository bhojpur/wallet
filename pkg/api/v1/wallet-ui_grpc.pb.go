// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// WalletUIClient is the client API for WalletUI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WalletUIClient interface {
	// ListEngineSpecs returns a list of Wallet Engine(s) that can be started through the UI.
	ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (WalletUI_ListEngineSpecsClient, error)
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error)
}

type walletUIClient struct {
	cc grpc.ClientConnInterface
}

func NewWalletUIClient(cc grpc.ClientConnInterface) WalletUIClient {
	return &walletUIClient{cc}
}

func (c *walletUIClient) ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (WalletUI_ListEngineSpecsClient, error) {
	stream, err := c.cc.NewStream(ctx, &WalletUI_ServiceDesc.Streams[0], "/v1.WalletUI/ListEngineSpecs", opts...)
	if err != nil {
		return nil, err
	}
	x := &walletUIListEngineSpecsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WalletUI_ListEngineSpecsClient interface {
	Recv() (*ListEngineSpecsResponse, error)
	grpc.ClientStream
}

type walletUIListEngineSpecsClient struct {
	grpc.ClientStream
}

func (x *walletUIListEngineSpecsClient) Recv() (*ListEngineSpecsResponse, error) {
	m := new(ListEngineSpecsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *walletUIClient) IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error) {
	out := new(IsReadOnlyResponse)
	err := c.cc.Invoke(ctx, "/v1.WalletUI/IsReadOnly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WalletUIServer is the server API for WalletUI service.
// All implementations must embed UnimplementedWalletUIServer
// for forward compatibility
type WalletUIServer interface {
	// ListEngineSpecs returns a list of Wallet Engine(s) that can be started through the UI.
	ListEngineSpecs(*ListEngineSpecsRequest, WalletUI_ListEngineSpecsServer) error
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error)
	mustEmbedUnimplementedWalletUIServer()
}

// UnimplementedWalletUIServer must be embedded to have forward compatible implementations.
type UnimplementedWalletUIServer struct {
}

func (UnimplementedWalletUIServer) ListEngineSpecs(*ListEngineSpecsRequest, WalletUI_ListEngineSpecsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListEngineSpecs not implemented")
}
func (UnimplementedWalletUIServer) IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReadOnly not implemented")
}
func (UnimplementedWalletUIServer) mustEmbedUnimplementedWalletUIServer() {}

// UnsafeWalletUIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WalletUIServer will
// result in compilation errors.
type UnsafeWalletUIServer interface {
	mustEmbedUnimplementedWalletUIServer()
}

func RegisterWalletUIServer(s grpc.ServiceRegistrar, srv WalletUIServer) {
	s.RegisterService(&WalletUI_ServiceDesc, srv)
}

func _WalletUI_ListEngineSpecs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListEngineSpecsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WalletUIServer).ListEngineSpecs(m, &walletUIListEngineSpecsServer{stream})
}

type WalletUI_ListEngineSpecsServer interface {
	Send(*ListEngineSpecsResponse) error
	grpc.ServerStream
}

type walletUIListEngineSpecsServer struct {
	grpc.ServerStream
}

func (x *walletUIListEngineSpecsServer) Send(m *ListEngineSpecsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _WalletUI_IsReadOnly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsReadOnlyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletUIServer).IsReadOnly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.WalletUI/IsReadOnly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletUIServer).IsReadOnly(ctx, req.(*IsReadOnlyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WalletUI_ServiceDesc is the grpc.ServiceDesc for WalletUI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WalletUI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.WalletUI",
	HandlerType: (*WalletUIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReadOnly",
			Handler:    _WalletUI_IsReadOnly_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListEngineSpecs",
			Handler:       _WalletUI_ListEngineSpecs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "wallet-ui.proto",
}
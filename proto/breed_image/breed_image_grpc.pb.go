// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: breed_image.proto

package breed_image

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

// BreedImageServiceClient is the client API for BreedImageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BreedImageServiceClient interface {
	Search(ctx context.Context, in *BreedImageSearchRequest, opts ...grpc.CallOption) (*BreedImageSearchResponse, error)
}

type breedImageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBreedImageServiceClient(cc grpc.ClientConnInterface) BreedImageServiceClient {
	return &breedImageServiceClient{cc}
}

func (c *breedImageServiceClient) Search(ctx context.Context, in *BreedImageSearchRequest, opts ...grpc.CallOption) (*BreedImageSearchResponse, error) {
	out := new(BreedImageSearchResponse)
	err := c.cc.Invoke(ctx, "/breed_image.BreedImageService/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BreedImageServiceServer is the server API for BreedImageService service.
// All implementations must embed UnimplementedBreedImageServiceServer
// for forward compatibility
type BreedImageServiceServer interface {
	Search(context.Context, *BreedImageSearchRequest) (*BreedImageSearchResponse, error)
	mustEmbedUnimplementedBreedImageServiceServer()
}

// UnimplementedBreedImageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBreedImageServiceServer struct {
}

func (UnimplementedBreedImageServiceServer) Search(context.Context, *BreedImageSearchRequest) (*BreedImageSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedBreedImageServiceServer) mustEmbedUnimplementedBreedImageServiceServer() {}

// UnsafeBreedImageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BreedImageServiceServer will
// result in compilation errors.
type UnsafeBreedImageServiceServer interface {
	mustEmbedUnimplementedBreedImageServiceServer()
}

func RegisterBreedImageServiceServer(s grpc.ServiceRegistrar, srv BreedImageServiceServer) {
	s.RegisterService(&BreedImageService_ServiceDesc, srv)
}

func _BreedImageService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BreedImageSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BreedImageServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/breed_image.BreedImageService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BreedImageServiceServer).Search(ctx, req.(*BreedImageSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BreedImageService_ServiceDesc is the grpc.ServiceDesc for BreedImageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BreedImageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "breed_image.BreedImageService",
	HandlerType: (*BreedImageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _BreedImageService_Search_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "breed_image.proto",
}

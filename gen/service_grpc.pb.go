// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	GetWorker(ctx context.Context, in *GetWorkerRequest, opts ...grpc.CallOption) (*GetWorkerResponse, error)
	PutWorker(ctx context.Context, in *PutWorkerRequest, opts ...grpc.CallOption) (*PutWorkerResponse, error)
	ListWorkers(ctx context.Context, in *ListWorkersRequest, opts ...grpc.CallOption) (*ListWorkersResponse, error)
	ShowWorker(ctx context.Context, in *ShowWorkerRequest, opts ...grpc.CallOption) (*ShowWorkerResponse, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) GetWorker(ctx context.Context, in *GetWorkerRequest, opts ...grpc.CallOption) (*GetWorkerResponse, error) {
	out := new(GetWorkerResponse)
	err := c.cc.Invoke(ctx, "/Service/GetWorker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) PutWorker(ctx context.Context, in *PutWorkerRequest, opts ...grpc.CallOption) (*PutWorkerResponse, error) {
	out := new(PutWorkerResponse)
	err := c.cc.Invoke(ctx, "/Service/PutWorker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ListWorkers(ctx context.Context, in *ListWorkersRequest, opts ...grpc.CallOption) (*ListWorkersResponse, error) {
	out := new(ListWorkersResponse)
	err := c.cc.Invoke(ctx, "/Service/ListWorkers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ShowWorker(ctx context.Context, in *ShowWorkerRequest, opts ...grpc.CallOption) (*ShowWorkerResponse, error) {
	out := new(ShowWorkerResponse)
	err := c.cc.Invoke(ctx, "/Service/ShowWorker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	GetWorker(context.Context, *GetWorkerRequest) (*GetWorkerResponse, error)
	PutWorker(context.Context, *PutWorkerRequest) (*PutWorkerResponse, error)
	ListWorkers(context.Context, *ListWorkersRequest) (*ListWorkersResponse, error)
	ShowWorker(context.Context, *ShowWorkerRequest) (*ShowWorkerResponse, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) GetWorker(context.Context, *GetWorkerRequest) (*GetWorkerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorker not implemented")
}
func (UnimplementedServiceServer) PutWorker(context.Context, *PutWorkerRequest) (*PutWorkerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutWorker not implemented")
}
func (UnimplementedServiceServer) ListWorkers(context.Context, *ListWorkersRequest) (*ListWorkersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWorkers not implemented")
}
func (UnimplementedServiceServer) ShowWorker(context.Context, *ShowWorkerRequest) (*ShowWorkerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShowWorker not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_GetWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWorkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/GetWorker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetWorker(ctx, req.(*GetWorkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_PutWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutWorkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).PutWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/PutWorker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).PutWorker(ctx, req.(*PutWorkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ListWorkers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWorkersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ListWorkers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/ListWorkers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ListWorkers(ctx, req.(*ListWorkersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ShowWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShowWorkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ShowWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/ShowWorker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ShowWorker(ctx, req.(*ShowWorkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWorker",
			Handler:    _Service_GetWorker_Handler,
		},
		{
			MethodName: "PutWorker",
			Handler:    _Service_PutWorker_Handler,
		},
		{
			MethodName: "ListWorkers",
			Handler:    _Service_ListWorkers_Handler,
		},
		{
			MethodName: "ShowWorker",
			Handler:    _Service_ShowWorker_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

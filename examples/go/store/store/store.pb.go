// Code generated by protoc-gen-go.
// source: store.proto
// DO NOT EDIT!

/*
Package store is a generated protocol buffer package.

It is generated from these files:
	store.proto

It has these top-level messages:
	Empty
	AddItemRequest
	RemoveItemRequest
	RemoveItemResponse
	QueryItemRequest
	QuantityResponse
*/
package store

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type AddItemRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *AddItemRequest) Reset()                    { *m = AddItemRequest{} }
func (m *AddItemRequest) String() string            { return proto.CompactTextString(m) }
func (*AddItemRequest) ProtoMessage()               {}
func (*AddItemRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AddItemRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RemoveItemRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *RemoveItemRequest) Reset()                    { *m = RemoveItemRequest{} }
func (m *RemoveItemRequest) String() string            { return proto.CompactTextString(m) }
func (*RemoveItemRequest) ProtoMessage()               {}
func (*RemoveItemRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RemoveItemRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RemoveItemResponse struct {
	WasSuccessful bool `protobuf:"varint,1,opt,name=was_successful,json=wasSuccessful" json:"was_successful,omitempty"`
}

func (m *RemoveItemResponse) Reset()                    { *m = RemoveItemResponse{} }
func (m *RemoveItemResponse) String() string            { return proto.CompactTextString(m) }
func (*RemoveItemResponse) ProtoMessage()               {}
func (*RemoveItemResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RemoveItemResponse) GetWasSuccessful() bool {
	if m != nil {
		return m.WasSuccessful
	}
	return false
}

type QueryItemRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *QueryItemRequest) Reset()                    { *m = QueryItemRequest{} }
func (m *QueryItemRequest) String() string            { return proto.CompactTextString(m) }
func (*QueryItemRequest) ProtoMessage()               {}
func (*QueryItemRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *QueryItemRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type QuantityResponse struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Count int32  `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
}

func (m *QuantityResponse) Reset()                    { *m = QuantityResponse{} }
func (m *QuantityResponse) String() string            { return proto.CompactTextString(m) }
func (*QuantityResponse) ProtoMessage()               {}
func (*QuantityResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *QuantityResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *QuantityResponse) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func init() {
	proto.RegisterType((*Empty)(nil), "store.Empty")
	proto.RegisterType((*AddItemRequest)(nil), "store.AddItemRequest")
	proto.RegisterType((*RemoveItemRequest)(nil), "store.RemoveItemRequest")
	proto.RegisterType((*RemoveItemResponse)(nil), "store.RemoveItemResponse")
	proto.RegisterType((*QueryItemRequest)(nil), "store.QueryItemRequest")
	proto.RegisterType((*QuantityResponse)(nil), "store.QuantityResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Store service

type StoreClient interface {
	AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*Empty, error)
	AddItems(ctx context.Context, opts ...grpc.CallOption) (Store_AddItemsClient, error)
	AddItemsError(ctx context.Context, opts ...grpc.CallOption) (Store_AddItemsErrorClient, error)
	RemoveItem(ctx context.Context, in *RemoveItemRequest, opts ...grpc.CallOption) (*RemoveItemResponse, error)
	RemoveItems(ctx context.Context, opts ...grpc.CallOption) (Store_RemoveItemsClient, error)
	ListInventory(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Store_ListInventoryClient, error)
	QueryQuantity(ctx context.Context, in *QueryItemRequest, opts ...grpc.CallOption) (*QuantityResponse, error)
	QueryQuantities(ctx context.Context, opts ...grpc.CallOption) (Store_QueryQuantitiesClient, error)
}

type storeClient struct {
	cc *grpc.ClientConn
}

func NewStoreClient(cc *grpc.ClientConn) StoreClient {
	return &storeClient{cc}
}

func (c *storeClient) AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/store.Store/AddItem", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) AddItems(ctx context.Context, opts ...grpc.CallOption) (Store_AddItemsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Store_serviceDesc.Streams[0], c.cc, "/store.Store/AddItems", opts...)
	if err != nil {
		return nil, err
	}
	x := &storeAddItemsClient{stream}
	return x, nil
}

type Store_AddItemsClient interface {
	Send(*AddItemRequest) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type storeAddItemsClient struct {
	grpc.ClientStream
}

func (x *storeAddItemsClient) Send(m *AddItemRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storeAddItemsClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storeClient) AddItemsError(ctx context.Context, opts ...grpc.CallOption) (Store_AddItemsErrorClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Store_serviceDesc.Streams[1], c.cc, "/store.Store/AddItemsError", opts...)
	if err != nil {
		return nil, err
	}
	x := &storeAddItemsErrorClient{stream}
	return x, nil
}

type Store_AddItemsErrorClient interface {
	Send(*AddItemRequest) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type storeAddItemsErrorClient struct {
	grpc.ClientStream
}

func (x *storeAddItemsErrorClient) Send(m *AddItemRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storeAddItemsErrorClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storeClient) RemoveItem(ctx context.Context, in *RemoveItemRequest, opts ...grpc.CallOption) (*RemoveItemResponse, error) {
	out := new(RemoveItemResponse)
	err := grpc.Invoke(ctx, "/store.Store/RemoveItem", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) RemoveItems(ctx context.Context, opts ...grpc.CallOption) (Store_RemoveItemsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Store_serviceDesc.Streams[2], c.cc, "/store.Store/RemoveItems", opts...)
	if err != nil {
		return nil, err
	}
	x := &storeRemoveItemsClient{stream}
	return x, nil
}

type Store_RemoveItemsClient interface {
	Send(*RemoveItemRequest) error
	CloseAndRecv() (*RemoveItemResponse, error)
	grpc.ClientStream
}

type storeRemoveItemsClient struct {
	grpc.ClientStream
}

func (x *storeRemoveItemsClient) Send(m *RemoveItemRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storeRemoveItemsClient) CloseAndRecv() (*RemoveItemResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RemoveItemResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storeClient) ListInventory(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Store_ListInventoryClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Store_serviceDesc.Streams[3], c.cc, "/store.Store/ListInventory", opts...)
	if err != nil {
		return nil, err
	}
	x := &storeListInventoryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Store_ListInventoryClient interface {
	Recv() (*QuantityResponse, error)
	grpc.ClientStream
}

type storeListInventoryClient struct {
	grpc.ClientStream
}

func (x *storeListInventoryClient) Recv() (*QuantityResponse, error) {
	m := new(QuantityResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storeClient) QueryQuantity(ctx context.Context, in *QueryItemRequest, opts ...grpc.CallOption) (*QuantityResponse, error) {
	out := new(QuantityResponse)
	err := grpc.Invoke(ctx, "/store.Store/QueryQuantity", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) QueryQuantities(ctx context.Context, opts ...grpc.CallOption) (Store_QueryQuantitiesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Store_serviceDesc.Streams[4], c.cc, "/store.Store/QueryQuantities", opts...)
	if err != nil {
		return nil, err
	}
	x := &storeQueryQuantitiesClient{stream}
	return x, nil
}

type Store_QueryQuantitiesClient interface {
	Send(*QueryItemRequest) error
	Recv() (*QuantityResponse, error)
	grpc.ClientStream
}

type storeQueryQuantitiesClient struct {
	grpc.ClientStream
}

func (x *storeQueryQuantitiesClient) Send(m *QueryItemRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storeQueryQuantitiesClient) Recv() (*QuantityResponse, error) {
	m := new(QuantityResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Store service

type StoreServer interface {
	AddItem(context.Context, *AddItemRequest) (*Empty, error)
	AddItems(Store_AddItemsServer) error
	AddItemsError(Store_AddItemsErrorServer) error
	RemoveItem(context.Context, *RemoveItemRequest) (*RemoveItemResponse, error)
	RemoveItems(Store_RemoveItemsServer) error
	ListInventory(*Empty, Store_ListInventoryServer) error
	QueryQuantity(context.Context, *QueryItemRequest) (*QuantityResponse, error)
	QueryQuantities(Store_QueryQuantitiesServer) error
}

func RegisterStoreServer(s *grpc.Server, srv StoreServer) {
	s.RegisterService(&_Store_serviceDesc, srv)
}

func _Store_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/store.Store/AddItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).AddItem(ctx, req.(*AddItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_AddItems_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StoreServer).AddItems(&storeAddItemsServer{stream})
}

type Store_AddItemsServer interface {
	SendAndClose(*Empty) error
	Recv() (*AddItemRequest, error)
	grpc.ServerStream
}

type storeAddItemsServer struct {
	grpc.ServerStream
}

func (x *storeAddItemsServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storeAddItemsServer) Recv() (*AddItemRequest, error) {
	m := new(AddItemRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Store_AddItemsError_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StoreServer).AddItemsError(&storeAddItemsErrorServer{stream})
}

type Store_AddItemsErrorServer interface {
	SendAndClose(*Empty) error
	Recv() (*AddItemRequest, error)
	grpc.ServerStream
}

type storeAddItemsErrorServer struct {
	grpc.ServerStream
}

func (x *storeAddItemsErrorServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storeAddItemsErrorServer) Recv() (*AddItemRequest, error) {
	m := new(AddItemRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Store_RemoveItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).RemoveItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/store.Store/RemoveItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).RemoveItem(ctx, req.(*RemoveItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_RemoveItems_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StoreServer).RemoveItems(&storeRemoveItemsServer{stream})
}

type Store_RemoveItemsServer interface {
	SendAndClose(*RemoveItemResponse) error
	Recv() (*RemoveItemRequest, error)
	grpc.ServerStream
}

type storeRemoveItemsServer struct {
	grpc.ServerStream
}

func (x *storeRemoveItemsServer) SendAndClose(m *RemoveItemResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storeRemoveItemsServer) Recv() (*RemoveItemRequest, error) {
	m := new(RemoveItemRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Store_ListInventory_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StoreServer).ListInventory(m, &storeListInventoryServer{stream})
}

type Store_ListInventoryServer interface {
	Send(*QuantityResponse) error
	grpc.ServerStream
}

type storeListInventoryServer struct {
	grpc.ServerStream
}

func (x *storeListInventoryServer) Send(m *QuantityResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Store_QueryQuantity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).QueryQuantity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/store.Store/QueryQuantity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).QueryQuantity(ctx, req.(*QueryItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_QueryQuantities_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StoreServer).QueryQuantities(&storeQueryQuantitiesServer{stream})
}

type Store_QueryQuantitiesServer interface {
	Send(*QuantityResponse) error
	Recv() (*QueryItemRequest, error)
	grpc.ServerStream
}

type storeQueryQuantitiesServer struct {
	grpc.ServerStream
}

func (x *storeQueryQuantitiesServer) Send(m *QuantityResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storeQueryQuantitiesServer) Recv() (*QueryItemRequest, error) {
	m := new(QueryItemRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Store_serviceDesc = grpc.ServiceDesc{
	ServiceName: "store.Store",
	HandlerType: (*StoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddItem",
			Handler:    _Store_AddItem_Handler,
		},
		{
			MethodName: "RemoveItem",
			Handler:    _Store_RemoveItem_Handler,
		},
		{
			MethodName: "QueryQuantity",
			Handler:    _Store_QueryQuantity_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AddItems",
			Handler:       _Store_AddItems_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "AddItemsError",
			Handler:       _Store_AddItemsError_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RemoveItems",
			Handler:       _Store_RemoveItems_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ListInventory",
			Handler:       _Store_ListInventory_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "QueryQuantities",
			Handler:       _Store_QueryQuantities_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "store.proto",
}

func init() { proto.RegisterFile("store.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 321 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x4f, 0x4f, 0xfa, 0x40,
	0x10, 0x65, 0x7f, 0xa1, 0x3f, 0x70, 0xb0, 0xa8, 0x13, 0x8d, 0xc8, 0x89, 0x6c, 0xfc, 0xd3, 0x13,
	0x21, 0x72, 0x31, 0xea, 0xc5, 0x18, 0x4c, 0x48, 0xbc, 0x58, 0x3e, 0x80, 0xa9, 0x30, 0x26, 0x24,
	0x76, 0x17, 0x77, 0xb6, 0x90, 0x1e, 0xfd, 0xe6, 0x86, 0xb6, 0x50, 0xaa, 0x44, 0x08, 0xb7, 0xce,
	0xcc, 0x7b, 0x6f, 0x3a, 0xef, 0x65, 0xa1, 0xc6, 0x56, 0x1b, 0x6a, 0x4f, 0x8c, 0xb6, 0x1a, 0x9d,
	0xa4, 0x90, 0x15, 0x70, 0x7a, 0xe1, 0xc4, 0xc6, 0xf2, 0x1c, 0xea, 0x0f, 0xa3, 0x51, 0xdf, 0x52,
	0xe8, 0xd3, 0x67, 0x44, 0x6c, 0x11, 0xa1, 0xac, 0x82, 0x90, 0x1a, 0xa2, 0x25, 0xbc, 0x3d, 0x3f,
	0xf9, 0x96, 0x57, 0x70, 0xe4, 0x53, 0xa8, 0xa7, 0xb4, 0x09, 0x78, 0x07, 0xb8, 0x0a, 0xe4, 0x89,
	0x56, 0x4c, 0x78, 0x01, 0xf5, 0x59, 0xc0, 0xaf, 0x1c, 0x0d, 0x87, 0xc4, 0xfc, 0x1e, 0x7d, 0x24,
	0x9c, 0xaa, 0xef, 0xce, 0x02, 0x1e, 0x2c, 0x9b, 0xf2, 0x12, 0x0e, 0x5f, 0x22, 0x32, 0xf1, 0xa6,
	0x25, 0xf7, 0x73, 0x5c, 0xa0, 0xec, 0xd8, 0xc6, 0xcb, 0x15, 0x6b, 0x70, 0x78, 0x0c, 0xce, 0x50,
	0x47, 0xca, 0x36, 0xfe, 0xb5, 0x84, 0xe7, 0xf8, 0x69, 0x71, 0xfd, 0x55, 0x06, 0x67, 0x30, 0x37,
	0x01, 0x3b, 0x50, 0xc9, 0x6e, 0xc7, 0x93, 0x76, 0x6a, 0x52, 0xd1, 0x8b, 0xe6, 0x7e, 0xd6, 0x4e,
	0xbd, 0x2a, 0x61, 0x17, 0xaa, 0x19, 0x82, 0xb7, 0xa4, 0x78, 0x02, 0x6f, 0xc0, 0x5d, 0x90, 0x7a,
	0xc6, 0x68, 0xb3, 0x3d, 0xf3, 0x11, 0x20, 0x77, 0x13, 0x1b, 0xd9, 0xfc, 0x57, 0x12, 0xcd, 0xb3,
	0x35, 0x93, 0xd4, 0x17, 0x59, 0xc2, 0x27, 0xa8, 0xe5, 0x7d, 0xde, 0x51, 0xc5, 0x13, 0x78, 0x0b,
	0xee, 0xf3, 0x98, 0x6d, 0x5f, 0x4d, 0x49, 0x59, 0x6d, 0x62, 0x2c, 0xfc, 0x6f, 0xf3, 0x34, 0xab,
	0x7e, 0x26, 0x23, 0x4b, 0x9d, 0xf9, 0x21, 0x6e, 0x92, 0xec, 0x62, 0x88, 0x39, 0xba, 0x98, 0xf7,
	0x1f, 0x32, 0xd8, 0x87, 0x83, 0x55, 0x91, 0x31, 0xf1, 0x2e, 0x32, 0x9e, 0xe8, 0x88, 0xb7, 0xff,
	0xc9, 0x63, 0xe8, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0x0d, 0xd1, 0x8b, 0x7b, 0x1b, 0x03, 0x00,
	0x00,
}

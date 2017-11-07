// Code generated by protoc-gen-go. DO NOT EDIT.
// source: alfa.proto

/*
Package alfa is a generated protocol buffer package.

Alfa Service

Alfa Service API consists of a single service which returns a message.

It is generated from these files:
	alfa.proto

It has these top-level messages:
	Message
*/
package alfa

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"
import beta "github.com/luigi-riefolo/alfa/src/beta/pb"

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

// Message represents a simple message sent to the Alfa service.
type Message struct {
	// Id represents the message identifier.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// The message to be sent.
	Msg string `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "alfa.Message")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AlfaService service

type AlfaServiceClient interface {
	// Get method receives a simple message and returns it.
	// The message posted as the id parameter will also be returned.
	Get(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*Message, error)
	// Alfa method sets a simple message.
	Set(ctx context.Context, in *Message, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	Test(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*beta.Message, error)
}

type alfaServiceClient struct {
	cc *grpc.ClientConn
}

func NewAlfaServiceClient(cc *grpc.ClientConn) AlfaServiceClient {
	return &alfaServiceClient{cc}
}

func (c *alfaServiceClient) Get(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := grpc.Invoke(ctx, "/alfa.AlfaService/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alfaServiceClient) Set(ctx context.Context, in *Message, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/alfa.AlfaService/Set", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alfaServiceClient) Test(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*beta.Message, error) {
	out := new(beta.Message)
	err := grpc.Invoke(ctx, "/alfa.AlfaService/Test", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AlfaService service

type AlfaServiceServer interface {
	// Get method receives a simple message and returns it.
	// The message posted as the id parameter will also be returned.
	Get(context.Context, *google_protobuf1.Empty) (*Message, error)
	// Alfa method sets a simple message.
	Set(context.Context, *Message) (*google_protobuf1.Empty, error)
	Test(context.Context, *google_protobuf1.Empty) (*beta.Message, error)
}

func RegisterAlfaServiceServer(s *grpc.Server, srv AlfaServiceServer) {
	s.RegisterService(&_AlfaService_serviceDesc, srv)
}

func _AlfaService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlfaServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alfa.AlfaService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlfaServiceServer).Get(ctx, req.(*google_protobuf1.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlfaService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlfaServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alfa.AlfaService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlfaServiceServer).Set(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlfaService_Test_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlfaServiceServer).Test(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alfa.AlfaService/Test",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlfaServiceServer).Test(ctx, req.(*google_protobuf1.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _AlfaService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "alfa.AlfaService",
	HandlerType: (*AlfaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _AlfaService_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _AlfaService_Set_Handler,
		},
		{
			MethodName: "Test",
			Handler:    _AlfaService_Test_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "alfa.proto",
}

func init() { proto.RegisterFile("alfa.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xbd, 0x4e, 0xc3, 0x30,
	0x10, 0x80, 0x95, 0xa4, 0x2a, 0xc2, 0xd0, 0x02, 0x16, 0x3f, 0x55, 0x60, 0x40, 0x99, 0x10, 0x08,
	0x5b, 0xfc, 0x4c, 0x6c, 0x54, 0x54, 0x4c, 0x2c, 0x94, 0x89, 0xcd, 0x49, 0x2f, 0xc6, 0x52, 0x12,
	0x47, 0xf1, 0xa5, 0x12, 0x2b, 0xaf, 0xc0, 0xa3, 0xb1, 0x33, 0xf1, 0x20, 0xc8, 0x97, 0x54, 0x88,
	0x01, 0x26, 0x9f, 0x7d, 0x77, 0x9f, 0xbe, 0xf3, 0x31, 0xa6, 0x8a, 0x5c, 0x89, 0xba, 0xb1, 0x68,
	0xf9, 0xc0, 0xc7, 0xf1, 0x91, 0xb6, 0x56, 0x17, 0x20, 0x55, 0x6d, 0xa4, 0xaa, 0x2a, 0x8b, 0x0a,
	0x8d, 0xad, 0x5c, 0x57, 0x13, 0x1f, 0xf6, 0x59, 0xba, 0xa5, 0x6d, 0x2e, 0xa1, 0xac, 0xf1, 0xb5,
	0x4f, 0x5e, 0x6b, 0x83, 0x2f, 0x6d, 0x2a, 0x32, 0x5b, 0xca, 0xa2, 0x35, 0xda, 0x9c, 0x37, 0x06,
	0x72, 0x5b, 0x58, 0xe9, 0xc9, 0xd2, 0x35, 0x99, 0x4c, 0x01, 0x95, 0xac, 0x53, 0x3a, 0xbb, 0xae,
	0xe4, 0x8c, 0xad, 0x3d, 0x80, 0x73, 0x4a, 0x03, 0x1f, 0xb3, 0xd0, 0x2c, 0x26, 0xc1, 0x71, 0x70,
	0xb2, 0xfe, 0x18, 0x9a, 0x05, 0xdf, 0x66, 0x51, 0xe9, 0xf4, 0x24, 0xa4, 0x07, 0x1f, 0x5e, 0x7e,
	0x06, 0x6c, 0xe3, 0xb6, 0xc8, 0xd5, 0x1c, 0x9a, 0xa5, 0xc9, 0x80, 0x4f, 0x59, 0x74, 0x0f, 0xc8,
	0xf7, 0x45, 0xe7, 0x25, 0x56, 0x5e, 0x62, 0xe6, 0xbd, 0xe2, 0x91, 0xa0, 0xf9, 0x7a, 0x7e, 0xb2,
	0xfb, 0xf6, 0xf1, 0xf5, 0x1e, 0x8e, 0xf9, 0xa6, 0x5c, 0x5e, 0x74, 0x4e, 0x1a, 0x90, 0xcf, 0x58,
	0x34, 0x07, 0xe4, 0xbf, 0x6b, 0xe3, 0x3f, 0x90, 0xc9, 0x01, 0x31, 0x76, 0x92, 0x1f, 0x86, 0x03,
	0xbc, 0x09, 0x4e, 0xf9, 0x1d, 0x1b, 0x3c, 0x81, 0xfb, 0xcf, 0x85, 0x86, 0x5e, 0xb9, 0xec, 0x11,
	0x67, 0x8b, 0x8f, 0x3c, 0x87, 0xbe, 0x05, 0xc1, 0xe1, 0x74, 0xf8, 0x4c, 0x6b, 0x48, 0x87, 0xd4,
	0x7d, 0xf5, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x94, 0xf8, 0x91, 0x17, 0xa1, 0x01, 0x00, 0x00,
}

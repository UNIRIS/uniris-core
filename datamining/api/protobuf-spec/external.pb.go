// Code generated by protoc-gen-go. DO NOT EDIT.
// source: external.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import any "github.com/golang/protobuf/ptypes/any"
import empty "github.com/golang/protobuf/ptypes/empty"

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

type TransactionType int32

const (
	TransactionType_CreateWallet TransactionType = 0
	TransactionType_CreateBio    TransactionType = 1
)

var TransactionType_name = map[int32]string{
	0: "CreateWallet",
	1: "CreateBio",
}
var TransactionType_value = map[string]int32{
	"CreateWallet": 0,
	"CreateBio":    1,
}

func (x TransactionType) String() string {
	return proto.EnumName(TransactionType_name, int32(x))
}
func (TransactionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_external_1f502708a57bff46, []int{0}
}

type Validation_ValidationStatus int32

const (
	Validation_OK Validation_ValidationStatus = 0
	Validation_KO Validation_ValidationStatus = 1
)

var Validation_ValidationStatus_name = map[int32]string{
	0: "OK",
	1: "KO",
}
var Validation_ValidationStatus_value = map[string]int32{
	"OK": 0,
	"KO": 1,
}

func (x Validation_ValidationStatus) String() string {
	return proto.EnumName(Validation_ValidationStatus_name, int32(x))
}
func (Validation_ValidationStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_external_1f502708a57bff46, []int{4, 0}
}

type LockRequest struct {
	TransactionHash      string   `protobuf:"bytes,1,opt,name=TransactionHash,proto3" json:"TransactionHash,omitempty"`
	MasterRobotKey       string   `protobuf:"bytes,2,opt,name=MasterRobotKey,proto3" json:"MasterRobotKey,omitempty"`
	Signature            string   `protobuf:"bytes,3,opt,name=Signature,proto3" json:"Signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LockRequest) Reset()         { *m = LockRequest{} }
func (m *LockRequest) String() string { return proto.CompactTextString(m) }
func (*LockRequest) ProtoMessage()    {}
func (*LockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_external_1f502708a57bff46, []int{0}
}
func (m *LockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LockRequest.Unmarshal(m, b)
}
func (m *LockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LockRequest.Marshal(b, m, deterministic)
}
func (dst *LockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LockRequest.Merge(dst, src)
}
func (m *LockRequest) XXX_Size() int {
	return xxx_messageInfo_LockRequest.Size(m)
}
func (m *LockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LockRequest proto.InternalMessageInfo

func (m *LockRequest) GetTransactionHash() string {
	if m != nil {
		return m.TransactionHash
	}
	return ""
}

func (m *LockRequest) GetMasterRobotKey() string {
	if m != nil {
		return m.MasterRobotKey
	}
	return ""
}

func (m *LockRequest) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

type StorageRequest struct {
	Data                 *any.Any        `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	TransactionType      TransactionType `protobuf:"varint,2,opt,name=TransactionType,proto3,enum=api.TransactionType" json:"TransactionType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *StorageRequest) Reset()         { *m = StorageRequest{} }
func (m *StorageRequest) String() string { return proto.CompactTextString(m) }
func (*StorageRequest) ProtoMessage()    {}
func (*StorageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_external_1f502708a57bff46, []int{1}
}
func (m *StorageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StorageRequest.Unmarshal(m, b)
}
func (m *StorageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StorageRequest.Marshal(b, m, deterministic)
}
func (dst *StorageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorageRequest.Merge(dst, src)
}
func (m *StorageRequest) XXX_Size() int {
	return xxx_messageInfo_StorageRequest.Size(m)
}
func (m *StorageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StorageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StorageRequest proto.InternalMessageInfo

func (m *StorageRequest) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *StorageRequest) GetTransactionType() TransactionType {
	if m != nil {
		return m.TransactionType
	}
	return TransactionType_CreateWallet
}

type ValidationRequest struct {
	Data                 *any.Any        `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	TransactionType      TransactionType `protobuf:"varint,2,opt,name=TransactionType,proto3,enum=api.TransactionType" json:"TransactionType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ValidationRequest) Reset()         { *m = ValidationRequest{} }
func (m *ValidationRequest) String() string { return proto.CompactTextString(m) }
func (*ValidationRequest) ProtoMessage()    {}
func (*ValidationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_external_1f502708a57bff46, []int{2}
}
func (m *ValidationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidationRequest.Unmarshal(m, b)
}
func (m *ValidationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidationRequest.Marshal(b, m, deterministic)
}
func (dst *ValidationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidationRequest.Merge(dst, src)
}
func (m *ValidationRequest) XXX_Size() int {
	return xxx_messageInfo_ValidationRequest.Size(m)
}
func (m *ValidationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ValidationRequest proto.InternalMessageInfo

func (m *ValidationRequest) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ValidationRequest) GetTransactionType() TransactionType {
	if m != nil {
		return m.TransactionType
	}
	return TransactionType_CreateWallet
}

type ValidationResponse struct {
	Validation           *Validation `protobuf:"bytes,1,opt,name=Validation,proto3" json:"Validation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ValidationResponse) Reset()         { *m = ValidationResponse{} }
func (m *ValidationResponse) String() string { return proto.CompactTextString(m) }
func (*ValidationResponse) ProtoMessage()    {}
func (*ValidationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_external_1f502708a57bff46, []int{3}
}
func (m *ValidationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidationResponse.Unmarshal(m, b)
}
func (m *ValidationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidationResponse.Marshal(b, m, deterministic)
}
func (dst *ValidationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidationResponse.Merge(dst, src)
}
func (m *ValidationResponse) XXX_Size() int {
	return xxx_messageInfo_ValidationResponse.Size(m)
}
func (m *ValidationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ValidationResponse proto.InternalMessageInfo

func (m *ValidationResponse) GetValidation() *Validation {
	if m != nil {
		return m.Validation
	}
	return nil
}

type Validation struct {
	Status               Validation_ValidationStatus `protobuf:"varint,1,opt,name=Status,proto3,enum=api.Validation_ValidationStatus" json:"Status,omitempty"`
	Timestamp            int64                       `protobuf:"varint,2,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	PublicKey            string                      `protobuf:"bytes,3,opt,name=PublicKey,proto3" json:"PublicKey,omitempty"`
	Signature            string                      `protobuf:"bytes,4,opt,name=Signature,proto3" json:"Signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *Validation) Reset()         { *m = Validation{} }
func (m *Validation) String() string { return proto.CompactTextString(m) }
func (*Validation) ProtoMessage()    {}
func (*Validation) Descriptor() ([]byte, []int) {
	return fileDescriptor_external_1f502708a57bff46, []int{4}
}
func (m *Validation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Validation.Unmarshal(m, b)
}
func (m *Validation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Validation.Marshal(b, m, deterministic)
}
func (dst *Validation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Validation.Merge(dst, src)
}
func (m *Validation) XXX_Size() int {
	return xxx_messageInfo_Validation.Size(m)
}
func (m *Validation) XXX_DiscardUnknown() {
	xxx_messageInfo_Validation.DiscardUnknown(m)
}

var xxx_messageInfo_Validation proto.InternalMessageInfo

func (m *Validation) GetStatus() Validation_ValidationStatus {
	if m != nil {
		return m.Status
	}
	return Validation_OK
}

func (m *Validation) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Validation) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

func (m *Validation) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func init() {
	proto.RegisterType((*LockRequest)(nil), "api.LockRequest")
	proto.RegisterType((*StorageRequest)(nil), "api.StorageRequest")
	proto.RegisterType((*ValidationRequest)(nil), "api.ValidationRequest")
	proto.RegisterType((*ValidationResponse)(nil), "api.ValidationResponse")
	proto.RegisterType((*Validation)(nil), "api.Validation")
	proto.RegisterEnum("api.TransactionType", TransactionType_name, TransactionType_value)
	proto.RegisterEnum("api.Validation_ValidationStatus", Validation_ValidationStatus_name, Validation_ValidationStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ExternalClient is the client API for External service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ExternalClient interface {
	LockTransaction(ctx context.Context, in *LockRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	UnlockTransaction(ctx context.Context, in *LockRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Validate(ctx context.Context, in *ValidationRequest, opts ...grpc.CallOption) (*ValidationResponse, error)
	Store(ctx context.Context, in *StorageRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type externalClient struct {
	cc *grpc.ClientConn
}

func NewExternalClient(cc *grpc.ClientConn) ExternalClient {
	return &externalClient{cc}
}

func (c *externalClient) LockTransaction(ctx context.Context, in *LockRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/api.External/LockTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalClient) UnlockTransaction(ctx context.Context, in *LockRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/api.External/UnlockTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalClient) Validate(ctx context.Context, in *ValidationRequest, opts ...grpc.CallOption) (*ValidationResponse, error) {
	out := new(ValidationResponse)
	err := c.cc.Invoke(ctx, "/api.External/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalClient) Store(ctx context.Context, in *StorageRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/api.External/Store", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExternalServer is the server API for External service.
type ExternalServer interface {
	LockTransaction(context.Context, *LockRequest) (*empty.Empty, error)
	UnlockTransaction(context.Context, *LockRequest) (*empty.Empty, error)
	Validate(context.Context, *ValidationRequest) (*ValidationResponse, error)
	Store(context.Context, *StorageRequest) (*empty.Empty, error)
}

func RegisterExternalServer(s *grpc.Server, srv ExternalServer) {
	s.RegisterService(&_External_serviceDesc, srv)
}

func _External_LockTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalServer).LockTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.External/LockTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalServer).LockTransaction(ctx, req.(*LockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _External_UnlockTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalServer).UnlockTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.External/UnlockTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalServer).UnlockTransaction(ctx, req.(*LockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _External_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.External/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalServer).Validate(ctx, req.(*ValidationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _External_Store_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalServer).Store(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.External/Store",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalServer).Store(ctx, req.(*StorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _External_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.External",
	HandlerType: (*ExternalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LockTransaction",
			Handler:    _External_LockTransaction_Handler,
		},
		{
			MethodName: "UnlockTransaction",
			Handler:    _External_UnlockTransaction_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _External_Validate_Handler,
		},
		{
			MethodName: "Store",
			Handler:    _External_Store_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "external.proto",
}

func init() { proto.RegisterFile("external.proto", fileDescriptor_external_1f502708a57bff46) }

var fileDescriptor_external_1f502708a57bff46 = []byte{
	// 442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x93, 0x4f, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0xe3, 0xa4, 0x44, 0xcd, 0x14, 0x1c, 0x77, 0xa8, 0x4a, 0x08, 0x1c, 0x2a, 0x1f, 0x50,
	0xc5, 0xc1, 0x91, 0x8c, 0x84, 0xb8, 0x14, 0xc4, 0x9f, 0x48, 0x48, 0x01, 0x15, 0x39, 0x01, 0xce,
	0x93, 0x30, 0x04, 0x0b, 0x67, 0xd7, 0x78, 0xd7, 0x12, 0x41, 0xea, 0x81, 0xef, 0xc6, 0x07, 0x43,
	0xbb, 0x36, 0xd8, 0xd9, 0x88, 0x13, 0x52, 0x4f, 0x89, 0x7e, 0xf3, 0x66, 0xe7, 0xad, 0xdf, 0x0e,
	0xf8, 0xfc, 0x5d, 0x73, 0x21, 0x28, 0x8b, 0xf2, 0x42, 0x6a, 0x89, 0x3d, 0xca, 0xd3, 0xf1, 0xbd,
	0xb5, 0x94, 0xeb, 0x8c, 0x27, 0x16, 0x2d, 0xcb, 0xcf, 0x13, 0xde, 0xe4, 0x7a, 0x5b, 0x29, 0xc6,
	0x77, 0xdd, 0x22, 0x89, 0xba, 0x14, 0x5e, 0xc1, 0xd1, 0x1b, 0xb9, 0xfa, 0x9a, 0xf0, 0xb7, 0x92,
	0x95, 0xc6, 0x73, 0x18, 0x2e, 0x0a, 0x12, 0x8a, 0x56, 0x3a, 0x95, 0xe2, 0x35, 0xa9, 0x2f, 0x23,
	0xef, 0xcc, 0x3b, 0x1f, 0x24, 0x2e, 0xc6, 0x07, 0xe0, 0xbf, 0x25, 0xa5, 0xb9, 0x48, 0xe4, 0x52,
	0xea, 0x19, 0x6f, 0x47, 0x5d, 0x2b, 0x74, 0x28, 0xde, 0x87, 0xc1, 0x3c, 0x5d, 0x0b, 0xd2, 0x65,
	0xc1, 0xa3, 0x9e, 0x95, 0x34, 0x20, 0xfc, 0x01, 0xfe, 0x5c, 0xcb, 0x82, 0xd6, 0xdc, 0x38, 0x38,
	0x78, 0x45, 0x9a, 0xec, 0xd8, 0xa3, 0xf8, 0x24, 0xaa, 0xac, 0x47, 0x7f, 0xac, 0x47, 0xcf, 0xc5,
	0x36, 0xb1, 0x0a, 0x7c, 0xba, 0xe3, 0x75, 0xb1, 0xcd, 0xd9, 0x5a, 0xf0, 0xe3, 0x93, 0x88, 0xf2,
	0x34, 0x72, 0x6a, 0x89, 0x2b, 0x0e, 0xaf, 0xe0, 0xf8, 0x03, 0x65, 0xe9, 0x27, 0x32, 0xe4, 0xfa,
	0xc7, 0x4f, 0x01, 0xdb, 0xe3, 0x55, 0x2e, 0x85, 0x62, 0x9c, 0x00, 0x34, 0xb4, 0x76, 0x31, 0xb4,
	0x07, 0xb6, 0xc4, 0x2d, 0x49, 0xf8, 0xcb, 0x6b, 0x77, 0xe0, 0x13, 0xe8, 0xcf, 0x35, 0xe9, 0x52,
	0xd9, 0x5e, 0x3f, 0x3e, 0x73, 0x7a, 0x5b, 0x7f, 0x2b, 0x5d, 0x52, 0xeb, 0x4d, 0x50, 0x8b, 0x74,
	0xc3, 0x4a, 0xd3, 0x26, 0xb7, 0x37, 0xe9, 0x25, 0x0d, 0x30, 0xd5, 0x77, 0xe5, 0x32, 0x4b, 0x57,
	0x26, 0xe9, 0x3a, 0xc6, 0xbf, 0x60, 0x37, 0xe4, 0x03, 0x37, 0xe4, 0x10, 0x02, 0x77, 0x2a, 0xf6,
	0xa1, 0x7b, 0x39, 0x0b, 0x3a, 0xe6, 0x77, 0x76, 0x19, 0x78, 0x0f, 0xe3, 0xbd, 0xaf, 0x89, 0x01,
	0xdc, 0x7c, 0x59, 0x30, 0x69, 0xfe, 0x48, 0x59, 0xc6, 0x3a, 0xe8, 0xe0, 0x2d, 0x18, 0x54, 0xe4,
	0x45, 0x2a, 0x03, 0x2f, 0xfe, 0xd9, 0x85, 0xc3, 0x69, 0xbd, 0x0b, 0x78, 0x01, 0x43, 0xf3, 0x90,
	0x5b, 0x87, 0x60, 0x60, 0xef, 0xde, 0x7a, 0xde, 0xe3, 0xd3, 0xbd, 0x3c, 0xa7, 0x66, 0x4d, 0xc2,
	0x0e, 0x3e, 0x83, 0xe3, 0xf7, 0x22, 0xfb, 0x8f, 0x03, 0x2e, 0xe0, 0xb0, 0xbe, 0x24, 0xe3, 0xa9,
	0x1b, 0x58, 0xdd, 0x7d, 0x67, 0x8f, 0x57, 0xa9, 0x87, 0x1d, 0x7c, 0x0c, 0x37, 0xcc, 0x22, 0x30,
	0xde, 0xb6, 0x9a, 0xdd, 0xa5, 0xf8, 0xf7, 0xd8, 0x65, 0xdf, 0x92, 0x47, 0xbf, 0x03, 0x00, 0x00,
	0xff, 0xff, 0xa9, 0xc0, 0x68, 0xc8, 0x15, 0x04, 0x00, 0x00,
}

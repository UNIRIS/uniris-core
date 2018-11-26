// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal.proto

package api

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

type AccountSearchRequest struct {
	EncryptedIDHash      string   `protobuf:"bytes,1,opt,name=EncryptedIDHash,proto3" json:"EncryptedIDHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountSearchRequest) Reset()         { *m = AccountSearchRequest{} }
func (m *AccountSearchRequest) String() string { return proto.CompactTextString(m) }
func (*AccountSearchRequest) ProtoMessage()    {}
func (*AccountSearchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_f4536b8c2a9c2411, []int{0}
}
func (m *AccountSearchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountSearchRequest.Unmarshal(m, b)
}
func (m *AccountSearchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountSearchRequest.Marshal(b, m, deterministic)
}
func (dst *AccountSearchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountSearchRequest.Merge(dst, src)
}
func (m *AccountSearchRequest) XXX_Size() int {
	return xxx_messageInfo_AccountSearchRequest.Size(m)
}
func (m *AccountSearchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountSearchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AccountSearchRequest proto.InternalMessageInfo

func (m *AccountSearchRequest) GetEncryptedIDHash() string {
	if m != nil {
		return m.EncryptedIDHash
	}
	return ""
}

type AccountSearchResult struct {
	EncryptedWallet      string   `protobuf:"bytes,1,opt,name=EncryptedWallet,proto3" json:"EncryptedWallet,omitempty"`
	EncryptedAESkey      string   `protobuf:"bytes,2,opt,name=EncryptedAESkey,proto3" json:"EncryptedAESkey,omitempty"`
	EncryptedAddress     string   `protobuf:"bytes,3,opt,name=EncryptedAddress,proto3" json:"EncryptedAddress,omitempty"`
	Signature            string   `protobuf:"bytes,4,opt,name=Signature,proto3" json:"Signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountSearchResult) Reset()         { *m = AccountSearchResult{} }
func (m *AccountSearchResult) String() string { return proto.CompactTextString(m) }
func (*AccountSearchResult) ProtoMessage()    {}
func (*AccountSearchResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_f4536b8c2a9c2411, []int{1}
}
func (m *AccountSearchResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountSearchResult.Unmarshal(m, b)
}
func (m *AccountSearchResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountSearchResult.Marshal(b, m, deterministic)
}
func (dst *AccountSearchResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountSearchResult.Merge(dst, src)
}
func (m *AccountSearchResult) XXX_Size() int {
	return xxx_messageInfo_AccountSearchResult.Size(m)
}
func (m *AccountSearchResult) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountSearchResult.DiscardUnknown(m)
}

var xxx_messageInfo_AccountSearchResult proto.InternalMessageInfo

func (m *AccountSearchResult) GetEncryptedWallet() string {
	if m != nil {
		return m.EncryptedWallet
	}
	return ""
}

func (m *AccountSearchResult) GetEncryptedAESkey() string {
	if m != nil {
		return m.EncryptedAESkey
	}
	return ""
}

func (m *AccountSearchResult) GetEncryptedAddress() string {
	if m != nil {
		return m.EncryptedAddress
	}
	return ""
}

func (m *AccountSearchResult) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

type KeychainCreationRequest struct {
	EncryptedKeychain    string   `protobuf:"bytes,1,opt,name=EncryptedKeychain,proto3" json:"EncryptedKeychain,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeychainCreationRequest) Reset()         { *m = KeychainCreationRequest{} }
func (m *KeychainCreationRequest) String() string { return proto.CompactTextString(m) }
func (*KeychainCreationRequest) ProtoMessage()    {}
func (*KeychainCreationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_f4536b8c2a9c2411, []int{2}
}
func (m *KeychainCreationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeychainCreationRequest.Unmarshal(m, b)
}
func (m *KeychainCreationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeychainCreationRequest.Marshal(b, m, deterministic)
}
func (dst *KeychainCreationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeychainCreationRequest.Merge(dst, src)
}
func (m *KeychainCreationRequest) XXX_Size() int {
	return xxx_messageInfo_KeychainCreationRequest.Size(m)
}
func (m *KeychainCreationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_KeychainCreationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_KeychainCreationRequest proto.InternalMessageInfo

func (m *KeychainCreationRequest) GetEncryptedKeychain() string {
	if m != nil {
		return m.EncryptedKeychain
	}
	return ""
}

type IDCreationRequest struct {
	EncryptedID          string   `protobuf:"bytes,1,opt,name=EncryptedID,proto3" json:"EncryptedID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDCreationRequest) Reset()         { *m = IDCreationRequest{} }
func (m *IDCreationRequest) String() string { return proto.CompactTextString(m) }
func (*IDCreationRequest) ProtoMessage()    {}
func (*IDCreationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_f4536b8c2a9c2411, []int{3}
}
func (m *IDCreationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDCreationRequest.Unmarshal(m, b)
}
func (m *IDCreationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDCreationRequest.Marshal(b, m, deterministic)
}
func (dst *IDCreationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDCreationRequest.Merge(dst, src)
}
func (m *IDCreationRequest) XXX_Size() int {
	return xxx_messageInfo_IDCreationRequest.Size(m)
}
func (m *IDCreationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IDCreationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IDCreationRequest proto.InternalMessageInfo

func (m *IDCreationRequest) GetEncryptedID() string {
	if m != nil {
		return m.EncryptedID
	}
	return ""
}

type CreationResult struct {
	TransactionHash      string   `protobuf:"bytes,1,opt,name=TransactionHash,proto3" json:"TransactionHash,omitempty"`
	MasterPeerIP         string   `protobuf:"bytes,2,opt,name=MasterPeerIP,proto3" json:"MasterPeerIP,omitempty"`
	Signature            string   `protobuf:"bytes,3,opt,name=Signature,proto3" json:"Signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreationResult) Reset()         { *m = CreationResult{} }
func (m *CreationResult) String() string { return proto.CompactTextString(m) }
func (*CreationResult) ProtoMessage()    {}
func (*CreationResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_f4536b8c2a9c2411, []int{4}
}
func (m *CreationResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreationResult.Unmarshal(m, b)
}
func (m *CreationResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreationResult.Marshal(b, m, deterministic)
}
func (dst *CreationResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreationResult.Merge(dst, src)
}
func (m *CreationResult) XXX_Size() int {
	return xxx_messageInfo_CreationResult.Size(m)
}
func (m *CreationResult) XXX_DiscardUnknown() {
	xxx_messageInfo_CreationResult.DiscardUnknown(m)
}

var xxx_messageInfo_CreationResult proto.InternalMessageInfo

func (m *CreationResult) GetTransactionHash() string {
	if m != nil {
		return m.TransactionHash
	}
	return ""
}

func (m *CreationResult) GetMasterPeerIP() string {
	if m != nil {
		return m.MasterPeerIP
	}
	return ""
}

func (m *CreationResult) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func init() {
	proto.RegisterType((*AccountSearchRequest)(nil), "api.AccountSearchRequest")
	proto.RegisterType((*AccountSearchResult)(nil), "api.AccountSearchResult")
	proto.RegisterType((*KeychainCreationRequest)(nil), "api.KeychainCreationRequest")
	proto.RegisterType((*IDCreationRequest)(nil), "api.IDCreationRequest")
	proto.RegisterType((*CreationResult)(nil), "api.CreationResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// InternalClient is the client API for Internal service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InternalClient interface {
	GetAccount(ctx context.Context, in *AccountSearchRequest, opts ...grpc.CallOption) (*AccountSearchResult, error)
	CreateKeychain(ctx context.Context, in *KeychainCreationRequest, opts ...grpc.CallOption) (*CreationResult, error)
	CreateID(ctx context.Context, in *IDCreationRequest, opts ...grpc.CallOption) (*CreationResult, error)
}

type internalClient struct {
	cc *grpc.ClientConn
}

func NewInternalClient(cc *grpc.ClientConn) InternalClient {
	return &internalClient{cc}
}

func (c *internalClient) GetAccount(ctx context.Context, in *AccountSearchRequest, opts ...grpc.CallOption) (*AccountSearchResult, error) {
	out := new(AccountSearchResult)
	err := c.cc.Invoke(ctx, "/api.Internal/GetAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) CreateKeychain(ctx context.Context, in *KeychainCreationRequest, opts ...grpc.CallOption) (*CreationResult, error) {
	out := new(CreationResult)
	err := c.cc.Invoke(ctx, "/api.Internal/CreateKeychain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) CreateID(ctx context.Context, in *IDCreationRequest, opts ...grpc.CallOption) (*CreationResult, error) {
	out := new(CreationResult)
	err := c.cc.Invoke(ctx, "/api.Internal/CreateID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InternalServer is the server API for Internal service.
type InternalServer interface {
	GetAccount(context.Context, *AccountSearchRequest) (*AccountSearchResult, error)
	CreateKeychain(context.Context, *KeychainCreationRequest) (*CreationResult, error)
	CreateID(context.Context, *IDCreationRequest) (*CreationResult, error)
}

func RegisterInternalServer(s *grpc.Server, srv InternalServer) {
	s.RegisterService(&_Internal_serviceDesc, srv)
}

func _Internal_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Internal/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).GetAccount(ctx, req.(*AccountSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_CreateKeychain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeychainCreationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).CreateKeychain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Internal/CreateKeychain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).CreateKeychain(ctx, req.(*KeychainCreationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_CreateID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDCreationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).CreateID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Internal/CreateID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).CreateID(ctx, req.(*IDCreationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Internal_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Internal",
	HandlerType: (*InternalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAccount",
			Handler:    _Internal_GetAccount_Handler,
		},
		{
			MethodName: "CreateKeychain",
			Handler:    _Internal_CreateKeychain_Handler,
		},
		{
			MethodName: "CreateID",
			Handler:    _Internal_CreateID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal.proto",
}

func init() { proto.RegisterFile("internal.proto", fileDescriptor_internal_f4536b8c2a9c2411) }

var fileDescriptor_internal_f4536b8c2a9c2411 = []byte{
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x41, 0x4f, 0xc2, 0x40,
	0x10, 0x85, 0xa9, 0x18, 0x03, 0xa3, 0x41, 0x59, 0x8c, 0x56, 0xc2, 0x81, 0xec, 0x89, 0x18, 0xc3,
	0x41, 0xe3, 0xc1, 0x9b, 0x04, 0x08, 0x36, 0xc6, 0x84, 0x80, 0x89, 0xe7, 0xb1, 0x4c, 0xa4, 0xb1,
	0xd9, 0xd6, 0xdd, 0xed, 0x81, 0xc4, 0x7f, 0xe5, 0x2f, 0xf1, 0x1f, 0x19, 0xca, 0x56, 0xda, 0x2d,
	0x5c, 0xbf, 0x79, 0xf3, 0xb2, 0xef, 0xed, 0x40, 0x23, 0x10, 0x9a, 0xa4, 0xc0, 0xb0, 0x1f, 0xcb,
	0x48, 0x47, 0xac, 0x8a, 0x71, 0xc0, 0x1f, 0xe1, 0x7c, 0xe0, 0xfb, 0x51, 0x22, 0xf4, 0x9c, 0x50,
	0xfa, 0xcb, 0x19, 0x7d, 0x25, 0xa4, 0x34, 0xeb, 0xc1, 0xe9, 0x58, 0xf8, 0x72, 0x15, 0x6b, 0x5a,
	0x78, 0xa3, 0x27, 0x54, 0x4b, 0xd7, 0xe9, 0x3a, 0xbd, 0xfa, 0xcc, 0xc6, 0xfc, 0xc7, 0x81, 0x96,
	0x65, 0xa1, 0x92, 0xb0, 0xe8, 0xf0, 0x86, 0x61, 0x48, 0xba, 0xe4, 0xb0, 0xc1, 0x05, 0xe5, 0x60,
	0x3c, 0xff, 0xa4, 0x95, 0x7b, 0x60, 0x29, 0x37, 0x98, 0x5d, 0xc3, 0xd9, 0x16, 0x2d, 0x16, 0x92,
	0x94, 0x72, 0xab, 0xa9, 0xb4, 0xc4, 0x59, 0x07, 0xea, 0xf3, 0xe0, 0x43, 0xa0, 0x4e, 0x24, 0xb9,
	0x87, 0xa9, 0x68, 0x0b, 0xf8, 0x04, 0x2e, 0x9f, 0x69, 0xe5, 0x2f, 0x31, 0x10, 0x43, 0x49, 0xa8,
	0x83, 0x48, 0x64, 0xd1, 0x6f, 0xa0, 0xf9, 0x6f, 0x96, 0x69, 0xcc, 0xd3, 0xcb, 0x03, 0x7e, 0x0f,
	0x4d, 0x6f, 0x64, 0x5b, 0x74, 0xe1, 0x38, 0x57, 0x93, 0x59, 0xce, 0x23, 0xfe, 0x0d, 0x8d, 0xed,
	0x52, 0xd6, 0xd7, 0xab, 0x44, 0xa1, 0xd0, 0x5f, 0xc3, 0x7c, 0xe3, 0x16, 0x66, 0x1c, 0x4e, 0x5e,
	0x50, 0x69, 0x92, 0x53, 0x22, 0xe9, 0x4d, 0x4d, 0x59, 0x05, 0x56, 0x4c, 0x5f, 0xb5, 0xd2, 0xdf,
	0xfe, 0x3a, 0x50, 0xf3, 0xcc, 0x35, 0xb0, 0x21, 0xc0, 0x84, 0xb4, 0xf9, 0x42, 0x76, 0xd5, 0xc7,
	0x38, 0xe8, 0xef, 0xba, 0x89, 0xb6, 0xbb, 0x6b, 0xb4, 0x7e, 0x3b, 0xaf, 0xb0, 0xb1, 0xc9, 0x43,
	0x59, 0x31, 0xac, 0x93, 0xaa, 0xf7, 0x94, 0xdc, 0x6e, 0xa5, 0xd3, 0x62, 0x05, 0xbc, 0xc2, 0x1e,
	0xa0, 0xb6, 0xb1, 0xf1, 0x46, 0xec, 0x22, 0x95, 0x94, 0xca, 0xdd, 0xb3, 0xfa, 0x7e, 0x94, 0x5e,
	0xf5, 0xdd, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x42, 0xb5, 0x1b, 0xff, 0xe7, 0x02, 0x00, 0x00,
}

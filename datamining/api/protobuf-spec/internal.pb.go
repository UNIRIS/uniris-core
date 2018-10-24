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

type WalletSearchRequest struct {
	EncryptedHashPerson  string   `protobuf:"bytes,1,opt,name=EncryptedHashPerson,proto3" json:"EncryptedHashPerson,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WalletSearchRequest) Reset()         { *m = WalletSearchRequest{} }
func (m *WalletSearchRequest) String() string { return proto.CompactTextString(m) }
func (*WalletSearchRequest) ProtoMessage()    {}
func (*WalletSearchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_e6e79e54ec80d1ae, []int{0}
}
func (m *WalletSearchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WalletSearchRequest.Unmarshal(m, b)
}
func (m *WalletSearchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WalletSearchRequest.Marshal(b, m, deterministic)
}
func (dst *WalletSearchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WalletSearchRequest.Merge(dst, src)
}
func (m *WalletSearchRequest) XXX_Size() int {
	return xxx_messageInfo_WalletSearchRequest.Size(m)
}
func (m *WalletSearchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WalletSearchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WalletSearchRequest proto.InternalMessageInfo

func (m *WalletSearchRequest) GetEncryptedHashPerson() string {
	if m != nil {
		return m.EncryptedHashPerson
	}
	return ""
}

type WalletStorageRequest struct {
	EncryptedBioData     string     `protobuf:"bytes,1,opt,name=EncryptedBioData,proto3" json:"EncryptedBioData,omitempty"`
	SignatureBioData     *Signature `protobuf:"bytes,2,opt,name=SignatureBioData,proto3" json:"SignatureBioData,omitempty"`
	EncryptedWalletData  string     `protobuf:"bytes,3,opt,name=EncryptedWalletData,proto3" json:"EncryptedWalletData,omitempty"`
	SignatureWalletData  *Signature `protobuf:"bytes,4,opt,name=SignatureWalletData,proto3" json:"SignatureWalletData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *WalletStorageRequest) Reset()         { *m = WalletStorageRequest{} }
func (m *WalletStorageRequest) String() string { return proto.CompactTextString(m) }
func (*WalletStorageRequest) ProtoMessage()    {}
func (*WalletStorageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_e6e79e54ec80d1ae, []int{1}
}
func (m *WalletStorageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WalletStorageRequest.Unmarshal(m, b)
}
func (m *WalletStorageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WalletStorageRequest.Marshal(b, m, deterministic)
}
func (dst *WalletStorageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WalletStorageRequest.Merge(dst, src)
}
func (m *WalletStorageRequest) XXX_Size() int {
	return xxx_messageInfo_WalletStorageRequest.Size(m)
}
func (m *WalletStorageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WalletStorageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WalletStorageRequest proto.InternalMessageInfo

func (m *WalletStorageRequest) GetEncryptedBioData() string {
	if m != nil {
		return m.EncryptedBioData
	}
	return ""
}

func (m *WalletStorageRequest) GetSignatureBioData() *Signature {
	if m != nil {
		return m.SignatureBioData
	}
	return nil
}

func (m *WalletStorageRequest) GetEncryptedWalletData() string {
	if m != nil {
		return m.EncryptedWalletData
	}
	return ""
}

func (m *WalletStorageRequest) GetSignatureWalletData() *Signature {
	if m != nil {
		return m.SignatureWalletData
	}
	return nil
}

type WalletStorageResult struct {
	BioTransactionHash   string   `protobuf:"bytes,1,opt,name=BioTransactionHash,proto3" json:"BioTransactionHash,omitempty"`
	DataTransactionHash  string   `protobuf:"bytes,2,opt,name=DataTransactionHash,proto3" json:"DataTransactionHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WalletStorageResult) Reset()         { *m = WalletStorageResult{} }
func (m *WalletStorageResult) String() string { return proto.CompactTextString(m) }
func (*WalletStorageResult) ProtoMessage()    {}
func (*WalletStorageResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_e6e79e54ec80d1ae, []int{2}
}
func (m *WalletStorageResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WalletStorageResult.Unmarshal(m, b)
}
func (m *WalletStorageResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WalletStorageResult.Marshal(b, m, deterministic)
}
func (dst *WalletStorageResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WalletStorageResult.Merge(dst, src)
}
func (m *WalletStorageResult) XXX_Size() int {
	return xxx_messageInfo_WalletStorageResult.Size(m)
}
func (m *WalletStorageResult) XXX_DiscardUnknown() {
	xxx_messageInfo_WalletStorageResult.DiscardUnknown(m)
}

var xxx_messageInfo_WalletStorageResult proto.InternalMessageInfo

func (m *WalletStorageResult) GetBioTransactionHash() string {
	if m != nil {
		return m.BioTransactionHash
	}
	return ""
}

func (m *WalletStorageResult) GetDataTransactionHash() string {
	if m != nil {
		return m.DataTransactionHash
	}
	return ""
}

type WalletSearchResult struct {
	EncryptedWallet        string   `protobuf:"bytes,1,opt,name=EncryptedWallet,proto3" json:"EncryptedWallet,omitempty"`
	EncryptedAESkey        string   `protobuf:"bytes,2,opt,name=EncryptedAESkey,proto3" json:"EncryptedAESkey,omitempty"`
	EncryptedWalletAddress string   `protobuf:"bytes,3,opt,name=EncryptedWalletAddress,proto3" json:"EncryptedWalletAddress,omitempty"`
	XXX_NoUnkeyedLiteral   struct{} `json:"-"`
	XXX_unrecognized       []byte   `json:"-"`
	XXX_sizecache          int32    `json:"-"`
}

func (m *WalletSearchResult) Reset()         { *m = WalletSearchResult{} }
func (m *WalletSearchResult) String() string { return proto.CompactTextString(m) }
func (*WalletSearchResult) ProtoMessage()    {}
func (*WalletSearchResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_e6e79e54ec80d1ae, []int{3}
}
func (m *WalletSearchResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WalletSearchResult.Unmarshal(m, b)
}
func (m *WalletSearchResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WalletSearchResult.Marshal(b, m, deterministic)
}
func (dst *WalletSearchResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WalletSearchResult.Merge(dst, src)
}
func (m *WalletSearchResult) XXX_Size() int {
	return xxx_messageInfo_WalletSearchResult.Size(m)
}
func (m *WalletSearchResult) XXX_DiscardUnknown() {
	xxx_messageInfo_WalletSearchResult.DiscardUnknown(m)
}

var xxx_messageInfo_WalletSearchResult proto.InternalMessageInfo

func (m *WalletSearchResult) GetEncryptedWallet() string {
	if m != nil {
		return m.EncryptedWallet
	}
	return ""
}

func (m *WalletSearchResult) GetEncryptedAESkey() string {
	if m != nil {
		return m.EncryptedAESkey
	}
	return ""
}

func (m *WalletSearchResult) GetEncryptedWalletAddress() string {
	if m != nil {
		return m.EncryptedWalletAddress
	}
	return ""
}

type Signature struct {
	Biod                 string   `protobuf:"bytes,1,opt,name=Biod,proto3" json:"Biod,omitempty"`
	Person               string   `protobuf:"bytes,2,opt,name=Person,proto3" json:"Person,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Signature) Reset()         { *m = Signature{} }
func (m *Signature) String() string { return proto.CompactTextString(m) }
func (*Signature) ProtoMessage()    {}
func (*Signature) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_e6e79e54ec80d1ae, []int{4}
}
func (m *Signature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Signature.Unmarshal(m, b)
}
func (m *Signature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Signature.Marshal(b, m, deterministic)
}
func (dst *Signature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Signature.Merge(dst, src)
}
func (m *Signature) XXX_Size() int {
	return xxx_messageInfo_Signature.Size(m)
}
func (m *Signature) XXX_DiscardUnknown() {
	xxx_messageInfo_Signature.DiscardUnknown(m)
}

var xxx_messageInfo_Signature proto.InternalMessageInfo

func (m *Signature) GetBiod() string {
	if m != nil {
		return m.Biod
	}
	return ""
}

func (m *Signature) GetPerson() string {
	if m != nil {
		return m.Person
	}
	return ""
}

func init() {
	proto.RegisterType((*WalletSearchRequest)(nil), "api.WalletSearchRequest")
	proto.RegisterType((*WalletStorageRequest)(nil), "api.WalletStorageRequest")
	proto.RegisterType((*WalletStorageResult)(nil), "api.WalletStorageResult")
	proto.RegisterType((*WalletSearchResult)(nil), "api.WalletSearchResult")
	proto.RegisterType((*Signature)(nil), "api.Signature")
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
	GetWallet(ctx context.Context, in *WalletSearchRequest, opts ...grpc.CallOption) (*WalletSearchResult, error)
	StoreWallet(ctx context.Context, in *WalletStorageRequest, opts ...grpc.CallOption) (*WalletStorageResult, error)
}

type internalClient struct {
	cc *grpc.ClientConn
}

func NewInternalClient(cc *grpc.ClientConn) InternalClient {
	return &internalClient{cc}
}

func (c *internalClient) GetWallet(ctx context.Context, in *WalletSearchRequest, opts ...grpc.CallOption) (*WalletSearchResult, error) {
	out := new(WalletSearchResult)
	err := c.cc.Invoke(ctx, "/api.Internal/GetWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) StoreWallet(ctx context.Context, in *WalletStorageRequest, opts ...grpc.CallOption) (*WalletStorageResult, error) {
	out := new(WalletStorageResult)
	err := c.cc.Invoke(ctx, "/api.Internal/StoreWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InternalServer is the server API for Internal service.
type InternalServer interface {
	GetWallet(context.Context, *WalletSearchRequest) (*WalletSearchResult, error)
	StoreWallet(context.Context, *WalletStorageRequest) (*WalletStorageResult, error)
}

func RegisterInternalServer(s *grpc.Server, srv InternalServer) {
	s.RegisterService(&_Internal_serviceDesc, srv)
}

func _Internal_GetWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WalletSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).GetWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Internal/GetWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).GetWallet(ctx, req.(*WalletSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_StoreWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WalletStorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).StoreWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Internal/StoreWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).StoreWallet(ctx, req.(*WalletStorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Internal_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Internal",
	HandlerType: (*InternalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWallet",
			Handler:    _Internal_GetWallet_Handler,
		},
		{
			MethodName: "StoreWallet",
			Handler:    _Internal_StoreWallet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal.proto",
}

func init() { proto.RegisterFile("internal.proto", fileDescriptor_internal_e6e79e54ec80d1ae) }

var fileDescriptor_internal_e6e79e54ec80d1ae = []byte{
	// 358 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0xcf, 0x4e, 0xfa, 0x40,
	0x10, 0xfe, 0x15, 0x08, 0xf9, 0x31, 0x24, 0x48, 0x06, 0x83, 0xc8, 0x89, 0xf4, 0x44, 0x3c, 0x34,
	0x06, 0x13, 0x4d, 0x3c, 0x01, 0x81, 0xa0, 0x37, 0x03, 0x26, 0x9e, 0x57, 0xba, 0x81, 0x8d, 0xcd,
	0x6e, 0xdd, 0x5d, 0x62, 0x78, 0x0d, 0x9f, 0xc0, 0x67, 0xf4, 0x09, 0x4c, 0xdb, 0xa1, 0xd2, 0xb2,
	0xde, 0xda, 0xf9, 0xbe, 0xf9, 0xe6, 0x9b, 0x3f, 0x0b, 0x2d, 0x21, 0x2d, 0xd7, 0x92, 0x45, 0x41,
	0xac, 0x95, 0x55, 0x58, 0x65, 0xb1, 0xf0, 0x17, 0xd0, 0x79, 0x61, 0x51, 0xc4, 0xed, 0x8a, 0x33,
	0xbd, 0xde, 0x2e, 0xf9, 0xfb, 0x8e, 0x1b, 0x8b, 0xd7, 0xd0, 0x99, 0xcb, 0xb5, 0xde, 0xc7, 0x96,
	0x87, 0x0f, 0xcc, 0x6c, 0x9f, 0xb8, 0x36, 0x4a, 0xf6, 0xbc, 0x81, 0x37, 0x6c, 0x2c, 0x5d, 0x90,
	0xff, 0xed, 0xc1, 0x39, 0x29, 0x59, 0xa5, 0xd9, 0x86, 0x1f, 0xa4, 0xae, 0xa0, 0x9d, 0xf3, 0xa7,
	0x42, 0xcd, 0x98, 0x65, 0xa4, 0x73, 0x12, 0xc7, 0x7b, 0x68, 0xaf, 0xc4, 0x46, 0x32, 0xbb, 0xd3,
	0xfc, 0xc0, 0xad, 0x0c, 0xbc, 0x61, 0x73, 0xd4, 0x0a, 0x58, 0x2c, 0x82, 0x1c, 0x5c, 0x9e, 0xf0,
	0x0a, 0x96, 0x33, 0x23, 0x69, 0x7a, 0xb5, 0x64, 0xf9, 0x17, 0xc2, 0x31, 0x74, 0x72, 0x95, 0xa3,
	0x8c, 0x9a, 0xb3, 0xa0, 0x8b, 0xea, 0x7f, 0xe4, 0xd3, 0x3b, 0xf4, 0x6c, 0x76, 0x91, 0xc5, 0x00,
	0x70, 0x2a, 0xd4, 0xb3, 0x66, 0xd2, 0xb0, 0xb5, 0x15, 0x4a, 0x26, 0x73, 0xa2, 0xa6, 0x1d, 0x48,
	0x62, 0x3d, 0x91, 0x2b, 0x27, 0x54, 0x32, 0xeb, 0x0e, 0xc8, 0xff, 0xf2, 0x00, 0x8b, 0x7b, 0x4b,
	0x0b, 0x0f, 0xe1, 0xac, 0xd4, 0x28, 0x55, 0x2d, 0x87, 0x0b, 0xcc, 0xc9, 0x7c, 0xf5, 0xc6, 0xf7,
	0x54, 0xae, 0x1c, 0xc6, 0x5b, 0xe8, 0x96, 0x92, 0x27, 0x61, 0xa8, 0xb9, 0x31, 0x34, 0xda, 0x3f,
	0x50, 0xff, 0x0e, 0x1a, 0xf9, 0xc8, 0x10, 0xa1, 0x36, 0x15, 0x2a, 0x24, 0x37, 0xe9, 0x37, 0x76,
	0xa1, 0x4e, 0x67, 0x95, 0x55, 0xa6, 0xbf, 0xd1, 0xa7, 0x07, 0xff, 0x1f, 0xe9, 0x54, 0x71, 0x0c,
	0x8d, 0x05, 0xb7, 0x64, 0xba, 0x97, 0xee, 0xc4, 0x71, 0xaf, 0xfd, 0x0b, 0x07, 0x92, 0x4c, 0xc4,
	0xff, 0x87, 0x33, 0x68, 0x26, 0xdb, 0xa1, 0xb5, 0xe1, 0xe5, 0x31, 0xb3, 0x70, 0xa9, 0xfd, 0x9e,
	0x0b, 0xca, 0x54, 0x5e, 0xeb, 0xe9, 0x9b, 0xb9, 0xf9, 0x09, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x51,
	0x33, 0x83, 0x45, 0x03, 0x00, 0x00,
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wallet.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

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

type GetWalletRequest struct {
	EncHashPerson        []byte   `protobuf:"bytes,1,opt,name=EncHashPerson,proto3" json:"EncHashPerson,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWalletRequest) Reset()         { *m = GetWalletRequest{} }
func (m *GetWalletRequest) String() string { return proto.CompactTextString(m) }
func (*GetWalletRequest) ProtoMessage()    {}
func (*GetWalletRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{0}
}

func (m *GetWalletRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWalletRequest.Unmarshal(m, b)
}
func (m *GetWalletRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWalletRequest.Marshal(b, m, deterministic)
}
func (m *GetWalletRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWalletRequest.Merge(m, src)
}
func (m *GetWalletRequest) XXX_Size() int {
	return xxx_messageInfo_GetWalletRequest.Size(m)
}
func (m *GetWalletRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWalletRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetWalletRequest proto.InternalMessageInfo

func (m *GetWalletRequest) GetEncHashPerson() []byte {
	if m != nil {
		return m.EncHashPerson
	}
	return nil
}

type GetWalletAck struct {
	Wd                   *WalletDetails `protobuf:"bytes,1,opt,name=Wd,proto3" json:"Wd,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetWalletAck) Reset()         { *m = GetWalletAck{} }
func (m *GetWalletAck) String() string { return proto.CompactTextString(m) }
func (*GetWalletAck) ProtoMessage()    {}
func (*GetWalletAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{1}
}

func (m *GetWalletAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWalletAck.Unmarshal(m, b)
}
func (m *GetWalletAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWalletAck.Marshal(b, m, deterministic)
}
func (m *GetWalletAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWalletAck.Merge(m, src)
}
func (m *GetWalletAck) XXX_Size() int {
	return xxx_messageInfo_GetWalletAck.Size(m)
}
func (m *GetWalletAck) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWalletAck.DiscardUnknown(m)
}

var xxx_messageInfo_GetWalletAck proto.InternalMessageInfo

func (m *GetWalletAck) GetWd() *WalletDetails {
	if m != nil {
		return m.Wd
	}
	return nil
}

type SetWalletRequest struct {
	EncryptedBioData     []byte         `protobuf:"bytes,1,opt,name=EncryptedBioData,proto3" json:"EncryptedBioData,omitempty"`
	SignatureBioData     *SignatureData `protobuf:"bytes,2,opt,name=SignatureBioData,proto3" json:"SignatureBioData,omitempty"`
	EncryptedWalletData  []byte         `protobuf:"bytes,3,opt,name=EncryptedWalletData,proto3" json:"EncryptedWalletData,omitempty"`
	SignatureWalletData  *SignatureData `protobuf:"bytes,4,opt,name=SignatureWalletData,proto3" json:"SignatureWalletData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SetWalletRequest) Reset()         { *m = SetWalletRequest{} }
func (m *SetWalletRequest) String() string { return proto.CompactTextString(m) }
func (*SetWalletRequest) ProtoMessage()    {}
func (*SetWalletRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{2}
}

func (m *SetWalletRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetWalletRequest.Unmarshal(m, b)
}
func (m *SetWalletRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetWalletRequest.Marshal(b, m, deterministic)
}
func (m *SetWalletRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetWalletRequest.Merge(m, src)
}
func (m *SetWalletRequest) XXX_Size() int {
	return xxx_messageInfo_SetWalletRequest.Size(m)
}
func (m *SetWalletRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetWalletRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetWalletRequest proto.InternalMessageInfo

func (m *SetWalletRequest) GetEncryptedBioData() []byte {
	if m != nil {
		return m.EncryptedBioData
	}
	return nil
}

func (m *SetWalletRequest) GetSignatureBioData() *SignatureData {
	if m != nil {
		return m.SignatureBioData
	}
	return nil
}

func (m *SetWalletRequest) GetEncryptedWalletData() []byte {
	if m != nil {
		return m.EncryptedWalletData
	}
	return nil
}

func (m *SetWalletRequest) GetSignatureWalletData() *SignatureData {
	if m != nil {
		return m.SignatureWalletData
	}
	return nil
}

type SetWalletAck struct {
	HashUpdatedWallet    []byte   `protobuf:"bytes,1,opt,name=HashUpdatedWallet,proto3" json:"HashUpdatedWallet,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetWalletAck) Reset()         { *m = SetWalletAck{} }
func (m *SetWalletAck) String() string { return proto.CompactTextString(m) }
func (*SetWalletAck) ProtoMessage()    {}
func (*SetWalletAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{3}
}

func (m *SetWalletAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetWalletAck.Unmarshal(m, b)
}
func (m *SetWalletAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetWalletAck.Marshal(b, m, deterministic)
}
func (m *SetWalletAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetWalletAck.Merge(m, src)
}
func (m *SetWalletAck) XXX_Size() int {
	return xxx_messageInfo_SetWalletAck.Size(m)
}
func (m *SetWalletAck) XXX_DiscardUnknown() {
	xxx_messageInfo_SetWalletAck.DiscardUnknown(m)
}

var xxx_messageInfo_SetWalletAck proto.InternalMessageInfo

func (m *SetWalletAck) GetHashUpdatedWallet() []byte {
	if m != nil {
		return m.HashUpdatedWallet
	}
	return nil
}

type WalletDetails struct {
	EncWallet            []byte   `protobuf:"bytes,1,opt,name=EncWallet,proto3" json:"EncWallet,omitempty"`
	EncAeskey            []byte   `protobuf:"bytes,2,opt,name=EncAeskey,proto3" json:"EncAeskey,omitempty"`
	EncWalletaddr        []byte   `protobuf:"bytes,3,opt,name=EncWalletaddr,proto3" json:"EncWalletaddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WalletDetails) Reset()         { *m = WalletDetails{} }
func (m *WalletDetails) String() string { return proto.CompactTextString(m) }
func (*WalletDetails) ProtoMessage()    {}
func (*WalletDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{4}
}

func (m *WalletDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WalletDetails.Unmarshal(m, b)
}
func (m *WalletDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WalletDetails.Marshal(b, m, deterministic)
}
func (m *WalletDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WalletDetails.Merge(m, src)
}
func (m *WalletDetails) XXX_Size() int {
	return xxx_messageInfo_WalletDetails.Size(m)
}
func (m *WalletDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_WalletDetails.DiscardUnknown(m)
}

var xxx_messageInfo_WalletDetails proto.InternalMessageInfo

func (m *WalletDetails) GetEncWallet() []byte {
	if m != nil {
		return m.EncWallet
	}
	return nil
}

func (m *WalletDetails) GetEncAeskey() []byte {
	if m != nil {
		return m.EncAeskey
	}
	return nil
}

func (m *WalletDetails) GetEncWalletaddr() []byte {
	if m != nil {
		return m.EncWalletaddr
	}
	return nil
}

type SignatureData struct {
	SigBiod              []byte   `protobuf:"bytes,1,opt,name=SigBiod,proto3" json:"SigBiod,omitempty"`
	SigPerson            []byte   `protobuf:"bytes,2,opt,name=SigPerson,proto3" json:"SigPerson,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignatureData) Reset()         { *m = SignatureData{} }
func (m *SignatureData) String() string { return proto.CompactTextString(m) }
func (*SignatureData) ProtoMessage()    {}
func (*SignatureData) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{5}
}

func (m *SignatureData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignatureData.Unmarshal(m, b)
}
func (m *SignatureData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignatureData.Marshal(b, m, deterministic)
}
func (m *SignatureData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignatureData.Merge(m, src)
}
func (m *SignatureData) XXX_Size() int {
	return xxx_messageInfo_SignatureData.Size(m)
}
func (m *SignatureData) XXX_DiscardUnknown() {
	xxx_messageInfo_SignatureData.DiscardUnknown(m)
}

var xxx_messageInfo_SignatureData proto.InternalMessageInfo

func (m *SignatureData) GetSigBiod() []byte {
	if m != nil {
		return m.SigBiod
	}
	return nil
}

func (m *SignatureData) GetSigPerson() []byte {
	if m != nil {
		return m.SigPerson
	}
	return nil
}

func init() {
	proto.RegisterType((*GetWalletRequest)(nil), "api.GetWalletRequest")
	proto.RegisterType((*GetWalletAck)(nil), "api.GetWalletAck")
	proto.RegisterType((*SetWalletRequest)(nil), "api.SetWalletRequest")
	proto.RegisterType((*SetWalletAck)(nil), "api.SetWalletAck")
	proto.RegisterType((*WalletDetails)(nil), "api.WalletDetails")
	proto.RegisterType((*SignatureData)(nil), "api.SignatureData")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WalletClient is the client API for Wallet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WalletClient interface {
	GetWallet(ctx context.Context, in *GetWalletRequest, opts ...grpc.CallOption) (*GetWalletAck, error)
	SetWallet(ctx context.Context, in *SetWalletRequest, opts ...grpc.CallOption) (*SetWalletAck, error)
}

type walletClient struct {
	cc *grpc.ClientConn
}

func NewWalletClient(cc *grpc.ClientConn) WalletClient {
	return &walletClient{cc}
}

func (c *walletClient) GetWallet(ctx context.Context, in *GetWalletRequest, opts ...grpc.CallOption) (*GetWalletAck, error) {
	out := new(GetWalletAck)
	err := c.cc.Invoke(ctx, "/api.Wallet/GetWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) SetWallet(ctx context.Context, in *SetWalletRequest, opts ...grpc.CallOption) (*SetWalletAck, error) {
	out := new(SetWalletAck)
	err := c.cc.Invoke(ctx, "/api.Wallet/SetWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WalletServer is the server API for Wallet service.
type WalletServer interface {
	GetWallet(context.Context, *GetWalletRequest) (*GetWalletAck, error)
	SetWallet(context.Context, *SetWalletRequest) (*SetWalletAck, error)
}

func RegisterWalletServer(s *grpc.Server, srv WalletServer) {
	s.RegisterService(&_Wallet_serviceDesc, srv)
}

func _Wallet_GetWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Wallet/GetWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetWallet(ctx, req.(*GetWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_SetWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).SetWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Wallet/SetWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).SetWallet(ctx, req.(*SetWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Wallet_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Wallet",
	HandlerType: (*WalletServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWallet",
			Handler:    _Wallet_GetWallet_Handler,
		},
		{
			MethodName: "SetWallet",
			Handler:    _Wallet_SetWallet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wallet.proto",
}

func init() { proto.RegisterFile("wallet.proto", fileDescriptor_b88fd140af4deb6f) }

var fileDescriptor_b88fd140af4deb6f = []byte{
	// 338 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0xcf, 0x4b, 0xf3, 0x30,
	0x18, 0xfe, 0xb6, 0x7d, 0x4c, 0xf6, 0xda, 0x41, 0x97, 0x21, 0x0c, 0xf1, 0x20, 0xc1, 0x83, 0x88,
	0x0c, 0x99, 0x07, 0x3d, 0x88, 0xb0, 0xb1, 0x31, 0x8f, 0xd2, 0x20, 0x3d, 0xc7, 0x26, 0xd4, 0xd0,
	0xd2, 0x76, 0x6d, 0x86, 0xcc, 0x3f, 0xdc, 0xb3, 0x34, 0xcd, 0x62, 0xd3, 0xd6, 0x63, 0x9f, 0x1f,
	0xef, 0x93, 0xf7, 0xe9, 0x0b, 0xce, 0x27, 0x8d, 0x63, 0x2e, 0xe7, 0x59, 0x9e, 0xca, 0x14, 0x0d,
	0x68, 0x26, 0xf0, 0x23, 0xb8, 0x5b, 0x2e, 0x7d, 0x85, 0x7b, 0x7c, 0xb7, 0xe7, 0x85, 0x44, 0x57,
	0x30, 0xde, 0x24, 0xc1, 0x0b, 0x2d, 0x3e, 0x5e, 0x79, 0x5e, 0xa4, 0xc9, 0xac, 0x77, 0xd9, 0xbb,
	0x76, 0x3c, 0x1b, 0xc4, 0x0b, 0x70, 0x8c, 0x73, 0x19, 0x44, 0x08, 0x43, 0xdf, 0x67, 0x4a, 0x7a,
	0xba, 0x40, 0x73, 0x9a, 0x89, 0x79, 0xc5, 0xad, 0xb9, 0xa4, 0x22, 0x2e, 0xbc, 0xbe, 0xcf, 0xf0,
	0x77, 0x0f, 0x5c, 0xd2, 0x8c, 0xbb, 0x01, 0x77, 0x93, 0x04, 0xf9, 0x21, 0x93, 0x9c, 0xad, 0x44,
	0xba, 0xa6, 0x92, 0xea, 0xc4, 0x16, 0x8e, 0x9e, 0xc1, 0x25, 0x22, 0x4c, 0xa8, 0xdc, 0xe7, 0xfc,
	0xa8, 0xed, 0xd7, 0x22, 0x0d, 0x59, 0x32, 0x5e, 0x4b, 0x8b, 0xee, 0x60, 0x6a, 0x66, 0xea, 0xe7,
	0x95, 0x23, 0x06, 0x2a, 0xae, 0x8b, 0x42, 0x6b, 0x98, 0x9a, 0x29, 0x35, 0xc7, 0xff, 0x3f, 0x43,
	0xbb, 0xe4, 0xf8, 0x09, 0x1c, 0x52, 0x2f, 0xeb, 0x16, 0x26, 0x65, 0x95, 0x6f, 0x19, 0xa3, 0x26,
	0x4e, 0x2f, 0xdd, 0x26, 0xf0, 0x0e, 0xc6, 0x56, 0x97, 0xe8, 0x02, 0x46, 0x9b, 0x24, 0xb0, 0x6c,
	0xbf, 0x80, 0x66, 0x97, 0xbc, 0x88, 0xf8, 0x41, 0xb5, 0x53, 0xb1, 0x15, 0xa0, 0xff, 0x6e, 0x25,
	0xa5, 0x8c, 0xe5, 0x7a, 0x79, 0x1b, 0xc4, 0x5b, 0x18, 0x5b, 0x6b, 0xa1, 0x19, 0x9c, 0x10, 0x11,
	0xae, 0x44, 0xca, 0x74, 0xe0, 0xf1, 0xb3, 0x8c, 0x23, 0x22, 0xd4, 0xa7, 0xa2, 0xe3, 0x0c, 0xb0,
	0xf8, 0x82, 0xa1, 0x7e, 0xd6, 0x03, 0x8c, 0xcc, 0xc1, 0xa0, 0x33, 0xd5, 0x5c, 0xf3, 0xf4, 0xce,
	0x27, 0x36, 0xbc, 0x0c, 0x22, 0xfc, 0xaf, 0x34, 0x92, 0x86, 0x91, 0x74, 0x1b, 0x89, 0x65, 0x7c,
	0x1f, 0xaa, 0x43, 0xbf, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0x33, 0x29, 0x86, 0x5e, 0xf8, 0x02,
	0x00, 0x00,
}

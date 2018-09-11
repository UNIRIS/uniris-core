// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gossip.proto

package gossip

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/golang/protobuf/ptypes/empty"

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

type Peer_PeerAppState_PeerState int32

const (
	Peer_PeerAppState_Fault        Peer_PeerAppState_PeerState = 0
	Peer_PeerAppState_Bootstraping Peer_PeerAppState_PeerState = 1
	Peer_PeerAppState_Ok           Peer_PeerAppState_PeerState = 2
	Peer_PeerAppState_StorageOnly  Peer_PeerAppState_PeerState = 3
)

var Peer_PeerAppState_PeerState_name = map[int32]string{
	0: "Fault",
	1: "Bootstraping",
	2: "Ok",
	3: "StorageOnly",
}
var Peer_PeerAppState_PeerState_value = map[string]int32{
	"Fault":        0,
	"Bootstraping": 1,
	"Ok":           2,
	"StorageOnly":  3,
}

func (x Peer_PeerAppState_PeerState) String() string {
	return proto.EnumName(Peer_PeerAppState_PeerState_name, int32(x))
}
func (Peer_PeerAppState_PeerState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_gossip_af78c85477b8f11a, []int{2, 1, 0}
}

type SynchronizeRequest struct {
	KnownPeers           []*Peer  `protobuf:"bytes,1,rep,name=KnownPeers,proto3" json:"KnownPeers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SynchronizeRequest) Reset()         { *m = SynchronizeRequest{} }
func (m *SynchronizeRequest) String() string { return proto.CompactTextString(m) }
func (*SynchronizeRequest) ProtoMessage()    {}
func (*SynchronizeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_gossip_af78c85477b8f11a, []int{0}
}
func (m *SynchronizeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SynchronizeRequest.Unmarshal(m, b)
}
func (m *SynchronizeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SynchronizeRequest.Marshal(b, m, deterministic)
}
func (dst *SynchronizeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SynchronizeRequest.Merge(dst, src)
}
func (m *SynchronizeRequest) XXX_Size() int {
	return xxx_messageInfo_SynchronizeRequest.Size(m)
}
func (m *SynchronizeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SynchronizeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SynchronizeRequest proto.InternalMessageInfo

func (m *SynchronizeRequest) GetKnownPeers() []*Peer {
	if m != nil {
		return m.KnownPeers
	}
	return nil
}

type AcknowledgeResponse struct {
	UnknownInitiatorPeers []*Peer  `protobuf:"bytes,1,rep,name=UnknownInitiatorPeers,proto3" json:"UnknownInitiatorPeers,omitempty"`
	WishedUnknownPeers    []*Peer  `protobuf:"bytes,2,rep,name=WishedUnknownPeers,proto3" json:"WishedUnknownPeers,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *AcknowledgeResponse) Reset()         { *m = AcknowledgeResponse{} }
func (m *AcknowledgeResponse) String() string { return proto.CompactTextString(m) }
func (*AcknowledgeResponse) ProtoMessage()    {}
func (*AcknowledgeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_gossip_af78c85477b8f11a, []int{1}
}
func (m *AcknowledgeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcknowledgeResponse.Unmarshal(m, b)
}
func (m *AcknowledgeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcknowledgeResponse.Marshal(b, m, deterministic)
}
func (dst *AcknowledgeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcknowledgeResponse.Merge(dst, src)
}
func (m *AcknowledgeResponse) XXX_Size() int {
	return xxx_messageInfo_AcknowledgeResponse.Size(m)
}
func (m *AcknowledgeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AcknowledgeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AcknowledgeResponse proto.InternalMessageInfo

func (m *AcknowledgeResponse) GetUnknownInitiatorPeers() []*Peer {
	if m != nil {
		return m.UnknownInitiatorPeers
	}
	return nil
}

func (m *AcknowledgeResponse) GetWishedUnknownPeers() []*Peer {
	if m != nil {
		return m.WishedUnknownPeers
	}
	return nil
}

type Peer struct {
	PublicKey            []byte             `protobuf:"bytes,1,opt,name=PublicKey,proto3" json:"PublicKey,omitempty"`
	IP                   string             `protobuf:"bytes,2,opt,name=IP,proto3" json:"IP,omitempty"`
	Port                 int32              `protobuf:"varint,3,opt,name=Port,proto3" json:"Port,omitempty"`
	HeartbeatState       *Peer_Heartbeat    `protobuf:"bytes,4,opt,name=HeartbeatState,proto3" json:"HeartbeatState,omitempty"`
	AppState             *Peer_PeerAppState `protobuf:"bytes,5,opt,name=AppState,proto3" json:"AppState,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Peer) Reset()         { *m = Peer{} }
func (m *Peer) String() string { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()    {}
func (*Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_gossip_af78c85477b8f11a, []int{2}
}
func (m *Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer.Unmarshal(m, b)
}
func (m *Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer.Marshal(b, m, deterministic)
}
func (dst *Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer.Merge(dst, src)
}
func (m *Peer) XXX_Size() int {
	return xxx_messageInfo_Peer.Size(m)
}
func (m *Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_Peer proto.InternalMessageInfo

func (m *Peer) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *Peer) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

func (m *Peer) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Peer) GetHeartbeatState() *Peer_Heartbeat {
	if m != nil {
		return m.HeartbeatState
	}
	return nil
}

func (m *Peer) GetAppState() *Peer_PeerAppState {
	if m != nil {
		return m.AppState
	}
	return nil
}

type Peer_Heartbeat struct {
	GenerationTime       int64    `protobuf:"varint,1,opt,name=GenerationTime,proto3" json:"GenerationTime,omitempty"`
	ElapsedBeats         int64    `protobuf:"varint,2,opt,name=ElapsedBeats,proto3" json:"ElapsedBeats,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peer_Heartbeat) Reset()         { *m = Peer_Heartbeat{} }
func (m *Peer_Heartbeat) String() string { return proto.CompactTextString(m) }
func (*Peer_Heartbeat) ProtoMessage()    {}
func (*Peer_Heartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_gossip_af78c85477b8f11a, []int{2, 0}
}
func (m *Peer_Heartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer_Heartbeat.Unmarshal(m, b)
}
func (m *Peer_Heartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer_Heartbeat.Marshal(b, m, deterministic)
}
func (dst *Peer_Heartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer_Heartbeat.Merge(dst, src)
}
func (m *Peer_Heartbeat) XXX_Size() int {
	return xxx_messageInfo_Peer_Heartbeat.Size(m)
}
func (m *Peer_Heartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer_Heartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_Peer_Heartbeat proto.InternalMessageInfo

func (m *Peer_Heartbeat) GetGenerationTime() int64 {
	if m != nil {
		return m.GenerationTime
	}
	return 0
}

func (m *Peer_Heartbeat) GetElapsedBeats() int64 {
	if m != nil {
		return m.ElapsedBeats
	}
	return 0
}

type Peer_PeerAppState struct {
	State                Peer_PeerAppState_PeerState    `protobuf:"varint,1,opt,name=State,proto3,enum=Peer_PeerAppState_PeerState" json:"State,omitempty"`
	CPULoad              string                         `protobuf:"bytes,2,opt,name=CPULoad,proto3" json:"CPULoad,omitempty"`
	IOWaitRate           float32                        `protobuf:"fixed32,3,opt,name=IOWaitRate,proto3" json:"IOWaitRate,omitempty"`
	FreeDiskSpace        float32                        `protobuf:"fixed32,4,opt,name=FreeDiskSpace,proto3" json:"FreeDiskSpace,omitempty"`
	Version              string                         `protobuf:"bytes,5,opt,name=Version,proto3" json:"Version,omitempty"`
	GeoCoordinates       *Peer_PeerAppState_Coordinates `protobuf:"bytes,6,opt,name=GeoCoordinates,proto3" json:"GeoCoordinates,omitempty"`
	P2PFactor            int32                          `protobuf:"varint,7,opt,name=P2PFactor,proto3" json:"P2PFactor,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *Peer_PeerAppState) Reset()         { *m = Peer_PeerAppState{} }
func (m *Peer_PeerAppState) String() string { return proto.CompactTextString(m) }
func (*Peer_PeerAppState) ProtoMessage()    {}
func (*Peer_PeerAppState) Descriptor() ([]byte, []int) {
	return fileDescriptor_gossip_af78c85477b8f11a, []int{2, 1}
}
func (m *Peer_PeerAppState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer_PeerAppState.Unmarshal(m, b)
}
func (m *Peer_PeerAppState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer_PeerAppState.Marshal(b, m, deterministic)
}
func (dst *Peer_PeerAppState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer_PeerAppState.Merge(dst, src)
}
func (m *Peer_PeerAppState) XXX_Size() int {
	return xxx_messageInfo_Peer_PeerAppState.Size(m)
}
func (m *Peer_PeerAppState) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer_PeerAppState.DiscardUnknown(m)
}

var xxx_messageInfo_Peer_PeerAppState proto.InternalMessageInfo

func (m *Peer_PeerAppState) GetState() Peer_PeerAppState_PeerState {
	if m != nil {
		return m.State
	}
	return Peer_PeerAppState_Fault
}

func (m *Peer_PeerAppState) GetCPULoad() string {
	if m != nil {
		return m.CPULoad
	}
	return ""
}

func (m *Peer_PeerAppState) GetIOWaitRate() float32 {
	if m != nil {
		return m.IOWaitRate
	}
	return 0
}

func (m *Peer_PeerAppState) GetFreeDiskSpace() float32 {
	if m != nil {
		return m.FreeDiskSpace
	}
	return 0
}

func (m *Peer_PeerAppState) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Peer_PeerAppState) GetGeoCoordinates() *Peer_PeerAppState_Coordinates {
	if m != nil {
		return m.GeoCoordinates
	}
	return nil
}

func (m *Peer_PeerAppState) GetP2PFactor() int32 {
	if m != nil {
		return m.P2PFactor
	}
	return 0
}

type Peer_PeerAppState_Coordinates struct {
	Lat                  float32  `protobuf:"fixed32,1,opt,name=Lat,proto3" json:"Lat,omitempty"`
	Lon                  float32  `protobuf:"fixed32,2,opt,name=Lon,proto3" json:"Lon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peer_PeerAppState_Coordinates) Reset()         { *m = Peer_PeerAppState_Coordinates{} }
func (m *Peer_PeerAppState_Coordinates) String() string { return proto.CompactTextString(m) }
func (*Peer_PeerAppState_Coordinates) ProtoMessage()    {}
func (*Peer_PeerAppState_Coordinates) Descriptor() ([]byte, []int) {
	return fileDescriptor_gossip_af78c85477b8f11a, []int{2, 1, 0}
}
func (m *Peer_PeerAppState_Coordinates) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer_PeerAppState_Coordinates.Unmarshal(m, b)
}
func (m *Peer_PeerAppState_Coordinates) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer_PeerAppState_Coordinates.Marshal(b, m, deterministic)
}
func (dst *Peer_PeerAppState_Coordinates) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer_PeerAppState_Coordinates.Merge(dst, src)
}
func (m *Peer_PeerAppState_Coordinates) XXX_Size() int {
	return xxx_messageInfo_Peer_PeerAppState_Coordinates.Size(m)
}
func (m *Peer_PeerAppState_Coordinates) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer_PeerAppState_Coordinates.DiscardUnknown(m)
}

var xxx_messageInfo_Peer_PeerAppState_Coordinates proto.InternalMessageInfo

func (m *Peer_PeerAppState_Coordinates) GetLat() float32 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Peer_PeerAppState_Coordinates) GetLon() float32 {
	if m != nil {
		return m.Lon
	}
	return 0
}

func init() {
	proto.RegisterType((*SynchronizeRequest)(nil), "SynchronizeRequest")
	proto.RegisterType((*AcknowledgeResponse)(nil), "AcknowledgeResponse")
	proto.RegisterType((*Peer)(nil), "Peer")
	proto.RegisterType((*Peer_Heartbeat)(nil), "Peer.Heartbeat")
	proto.RegisterType((*Peer_PeerAppState)(nil), "Peer.PeerAppState")
	proto.RegisterType((*Peer_PeerAppState_Coordinates)(nil), "Peer.PeerAppState.Coordinates")
	proto.RegisterEnum("Peer_PeerAppState_PeerState", Peer_PeerAppState_PeerState_name, Peer_PeerAppState_PeerState_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GossipServiceClient is the client API for GossipService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GossipServiceClient interface {
	Synchronize(ctx context.Context, in *SynchronizeRequest, opts ...grpc.CallOption) (*AcknowledgeResponse, error)
}

type gossipServiceClient struct {
	cc *grpc.ClientConn
}

func NewGossipServiceClient(cc *grpc.ClientConn) GossipServiceClient {
	return &gossipServiceClient{cc}
}

func (c *gossipServiceClient) Synchronize(ctx context.Context, in *SynchronizeRequest, opts ...grpc.CallOption) (*AcknowledgeResponse, error) {
	out := new(AcknowledgeResponse)
	err := c.cc.Invoke(ctx, "/GossipService/Synchronize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GossipServiceServer is the server API for GossipService service.
type GossipServiceServer interface {
	Synchronize(context.Context, *SynchronizeRequest) (*AcknowledgeResponse, error)
}

func RegisterGossipServiceServer(s *grpc.Server, srv GossipServiceServer) {
	s.RegisterService(&_GossipService_serviceDesc, srv)
}

func _GossipService_Synchronize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SynchronizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipServiceServer).Synchronize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GossipService/Synchronize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipServiceServer).Synchronize(ctx, req.(*SynchronizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GossipService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "GossipService",
	HandlerType: (*GossipServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Synchronize",
			Handler:    _GossipService_Synchronize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gossip.proto",
}

func init() { proto.RegisterFile("gossip.proto", fileDescriptor_gossip_af78c85477b8f11a) }

var fileDescriptor_gossip_af78c85477b8f11a = []byte{
	// 553 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x53, 0xd1, 0x4e, 0xdb, 0x4a,
	0x10, 0xc5, 0x36, 0x81, 0x9b, 0x49, 0x08, 0xd1, 0x70, 0x2b, 0x59, 0x29, 0x42, 0x51, 0xd4, 0x56,
	0x79, 0x32, 0xaa, 0xab, 0xaa, 0x52, 0x79, 0x02, 0xda, 0xd0, 0x08, 0x24, 0xac, 0x4d, 0x29, 0xcf,
	0x1b, 0x67, 0x6a, 0x56, 0x98, 0x5d, 0x77, 0x77, 0x53, 0x94, 0xfe, 0x41, 0xff, 0xa0, 0x3f, 0xd0,
	0xff, 0xac, 0xbc, 0x86, 0x60, 0x20, 0x2f, 0xd6, 0xcc, 0x99, 0x73, 0x66, 0x67, 0x8f, 0x77, 0xa0,
	0x9d, 0x29, 0x63, 0x44, 0x11, 0x15, 0x5a, 0x59, 0xd5, 0x7b, 0x99, 0x29, 0x95, 0xe5, 0xb4, 0xef,
	0xb2, 0xe9, 0xfc, 0xfb, 0x3e, 0xdd, 0x14, 0x76, 0x51, 0x15, 0x07, 0x07, 0x80, 0x93, 0x85, 0x4c,
	0xaf, 0xb4, 0x92, 0xe2, 0x17, 0x31, 0xfa, 0x31, 0x27, 0x63, 0xf1, 0x35, 0xc0, 0xa9, 0x54, 0xb7,
	0x32, 0x21, 0xd2, 0x26, 0xf4, 0xfa, 0xc1, 0xb0, 0x15, 0x37, 0xa2, 0x32, 0x63, 0xb5, 0xc2, 0xe0,
	0xb7, 0x07, 0x3b, 0x87, 0xe9, 0xb5, 0x54, 0xb7, 0x39, 0xcd, 0x32, 0x62, 0x64, 0x0a, 0x25, 0x0d,
	0xe1, 0x01, 0xbc, 0xb8, 0x90, 0x25, 0x2c, 0xc7, 0x52, 0x58, 0xc1, 0xad, 0xd2, 0x2b, 0x3a, 0xad,
	0xe6, 0xe0, 0x7b, 0xc0, 0x4b, 0x61, 0xae, 0x68, 0x76, 0x57, 0xae, 0x94, 0x7e, 0x5d, 0xb9, 0x82,
	0x30, 0xf8, 0xdb, 0x80, 0xf5, 0x32, 0xc2, 0x5d, 0x68, 0x26, 0xf3, 0x69, 0x2e, 0xd2, 0x53, 0x5a,
	0x84, 0x5e, 0xdf, 0x1b, 0xb6, 0xd9, 0x03, 0x80, 0x1d, 0xf0, 0xc7, 0x49, 0xe8, 0xf7, 0xbd, 0x61,
	0x93, 0xf9, 0xe3, 0x04, 0x11, 0xd6, 0x13, 0xa5, 0x6d, 0x18, 0xf4, 0xbd, 0x61, 0x83, 0xb9, 0x18,
	0x3f, 0x40, 0xe7, 0x0b, 0x71, 0x6d, 0xa7, 0xc4, 0xed, 0xc4, 0x72, 0x4b, 0xe1, 0x7a, 0xdf, 0x1b,
	0xb6, 0xe2, 0x6d, 0x77, 0x7a, 0xb4, 0xac, 0xb1, 0x27, 0x34, 0x8c, 0xe0, 0xbf, 0xc3, 0xa2, 0xa8,
	0x24, 0x0d, 0x27, 0xc1, 0x4a, 0x52, 0x7e, 0xee, 0x2b, 0x6c, 0xc9, 0xe9, 0x5d, 0x42, 0x73, 0xd9,
	0x01, 0xdf, 0x40, 0xe7, 0x84, 0x24, 0x69, 0x6e, 0x85, 0x92, 0x5f, 0xc5, 0x0d, 0xb9, 0xe1, 0x03,
	0xf6, 0x04, 0xc5, 0x01, 0xb4, 0x3f, 0xe7, 0xbc, 0x30, 0x34, 0x3b, 0x22, 0x6e, 0x8d, 0xbb, 0x4b,
	0xc0, 0x1e, 0x61, 0xbd, 0x3f, 0x01, 0xb4, 0xeb, 0x67, 0x62, 0x0c, 0x8d, 0x6a, 0xac, 0xb2, 0x67,
	0x27, 0xde, 0x7d, 0x3e, 0x96, 0x4b, 0xaa, 0x01, 0x2b, 0x2a, 0x86, 0xb0, 0x79, 0x9c, 0x5c, 0x9c,
	0x29, 0x3e, 0xbb, 0xf3, 0xeb, 0x3e, 0xc5, 0x3d, 0x80, 0xf1, 0xf9, 0x25, 0x17, 0x96, 0x95, 0x2d,
	0x4b, 0xeb, 0x7c, 0x56, 0x43, 0xf0, 0x15, 0x6c, 0x8d, 0x34, 0xd1, 0x27, 0x61, 0xae, 0x27, 0x05,
	0x4f, 0x2b, 0xff, 0x7c, 0xf6, 0x18, 0x2c, 0xfb, 0x7f, 0x23, 0x6d, 0x84, 0x92, 0xce, 0xac, 0x26,
	0xbb, 0x4f, 0x71, 0x54, 0x5a, 0xa1, 0x8e, 0x95, 0xd2, 0x33, 0x21, 0xb9, 0x25, 0x13, 0x6e, 0x38,
	0x37, 0xf7, 0x56, 0x8c, 0x5d, 0x63, 0xb1, 0x27, 0x2a, 0xf7, 0x14, 0xe2, 0x64, 0xc4, 0x53, 0xab,
	0x74, 0xb8, 0xe9, 0xfe, 0xf0, 0x03, 0xd0, 0x7b, 0x0b, 0xad, 0x3a, 0xb9, 0x0b, 0xc1, 0x19, 0xb7,
	0xce, 0x20, 0x9f, 0x95, 0xa1, 0x43, 0x94, 0x74, 0x97, 0x2f, 0x11, 0x25, 0x07, 0x87, 0xd0, 0x5c,
	0xda, 0x84, 0x4d, 0x68, 0x8c, 0xf8, 0x3c, 0xb7, 0xdd, 0x35, 0xec, 0x42, 0xfb, 0x48, 0x29, 0x6b,
	0xac, 0xe6, 0x85, 0x90, 0x59, 0xd7, 0xc3, 0x0d, 0xf0, 0xcf, 0xaf, 0xbb, 0x3e, 0x6e, 0x43, 0x6b,
	0x62, 0x95, 0xe6, 0x19, 0x9d, 0xcb, 0x7c, 0xd1, 0x0d, 0xe2, 0x53, 0xd8, 0x3a, 0x71, 0xdb, 0x39,
	0x21, 0xfd, 0x53, 0xa4, 0x84, 0x1f, 0xa1, 0x55, 0xdb, 0x40, 0xdc, 0x89, 0x9e, 0xef, 0x63, 0xef,
	0xff, 0x68, 0xc5, 0x9a, 0x0d, 0xd6, 0xa6, 0x1b, 0x6e, 0x89, 0xdf, 0xfd, 0x0b, 0x00, 0x00, 0xff,
	0xff, 0xaa, 0x93, 0x5b, 0x75, 0xf1, 0x03, 0x00, 0x00,
}

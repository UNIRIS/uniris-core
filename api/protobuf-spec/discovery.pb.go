// Code generated by protoc-gen-go. DO NOT EDIT.
// source: discovery.proto

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

type PeerAppState_PeerStatus int32

const (
	PeerAppState_BOOTSTRAPING PeerAppState_PeerStatus = 0
	PeerAppState_OK           PeerAppState_PeerStatus = 1
	PeerAppState_FAULTY       PeerAppState_PeerStatus = 2
	PeerAppState_STORAGE_ONLY PeerAppState_PeerStatus = 3
)

var PeerAppState_PeerStatus_name = map[int32]string{
	0: "BOOTSTRAPING",
	1: "OK",
	2: "FAULTY",
	3: "STORAGE_ONLY",
}
var PeerAppState_PeerStatus_value = map[string]int32{
	"BOOTSTRAPING": 0,
	"OK":           1,
	"FAULTY":       2,
	"STORAGE_ONLY": 3,
}

func (x PeerAppState_PeerStatus) String() string {
	return proto.EnumName(PeerAppState_PeerStatus_name, int32(x))
}
func (PeerAppState_PeerStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{8, 0}
}

type SynRequest struct {
	KnownPeers           []*PeerDigest `protobuf:"bytes,1,rep,name=known_peers,json=knownPeers,proto3" json:"known_peers,omitempty"`
	Timestamp            int64         `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SynRequest) Reset()         { *m = SynRequest{} }
func (m *SynRequest) String() string { return proto.CompactTextString(m) }
func (*SynRequest) ProtoMessage()    {}
func (*SynRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{0}
}
func (m *SynRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SynRequest.Unmarshal(m, b)
}
func (m *SynRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SynRequest.Marshal(b, m, deterministic)
}
func (dst *SynRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SynRequest.Merge(dst, src)
}
func (m *SynRequest) XXX_Size() int {
	return xxx_messageInfo_SynRequest.Size(m)
}
func (m *SynRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SynRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SynRequest proto.InternalMessageInfo

func (m *SynRequest) GetKnownPeers() []*PeerDigest {
	if m != nil {
		return m.KnownPeers
	}
	return nil
}

func (m *SynRequest) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type SynResponse struct {
	RemoteDiscoveris     []*PeerDiscovered `protobuf:"bytes,1,rep,name=remote_discoveris,json=remoteDiscoveris,proto3" json:"remote_discoveris,omitempty"`
	LocalDiscoveries     []*PeerDigest     `protobuf:"bytes,2,rep,name=local_discoveries,json=localDiscoveries,proto3" json:"local_discoveries,omitempty"`
	Timestamp            int64             `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SynResponse) Reset()         { *m = SynResponse{} }
func (m *SynResponse) String() string { return proto.CompactTextString(m) }
func (*SynResponse) ProtoMessage()    {}
func (*SynResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{1}
}
func (m *SynResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SynResponse.Unmarshal(m, b)
}
func (m *SynResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SynResponse.Marshal(b, m, deterministic)
}
func (dst *SynResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SynResponse.Merge(dst, src)
}
func (m *SynResponse) XXX_Size() int {
	return xxx_messageInfo_SynResponse.Size(m)
}
func (m *SynResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SynResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SynResponse proto.InternalMessageInfo

func (m *SynResponse) GetRemoteDiscoveris() []*PeerDiscovered {
	if m != nil {
		return m.RemoteDiscoveris
	}
	return nil
}

func (m *SynResponse) GetLocalDiscoveries() []*PeerDigest {
	if m != nil {
		return m.LocalDiscoveries
	}
	return nil
}

func (m *SynResponse) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type AckRequest struct {
	RequestedPeers       []*PeerDiscovered `protobuf:"bytes,1,rep,name=requested_peers,json=requestedPeers,proto3" json:"requested_peers,omitempty"`
	Timestamp            int64             `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AckRequest) Reset()         { *m = AckRequest{} }
func (m *AckRequest) String() string { return proto.CompactTextString(m) }
func (*AckRequest) ProtoMessage()    {}
func (*AckRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{2}
}
func (m *AckRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckRequest.Unmarshal(m, b)
}
func (m *AckRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckRequest.Marshal(b, m, deterministic)
}
func (dst *AckRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckRequest.Merge(dst, src)
}
func (m *AckRequest) XXX_Size() int {
	return xxx_messageInfo_AckRequest.Size(m)
}
func (m *AckRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AckRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AckRequest proto.InternalMessageInfo

func (m *AckRequest) GetRequestedPeers() []*PeerDiscovered {
	if m != nil {
		return m.RequestedPeers
	}
	return nil
}

func (m *AckRequest) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type AckResponse struct {
	Timestamp            int64    `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AckResponse) Reset()         { *m = AckResponse{} }
func (m *AckResponse) String() string { return proto.CompactTextString(m) }
func (*AckResponse) ProtoMessage()    {}
func (*AckResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{3}
}
func (m *AckResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckResponse.Unmarshal(m, b)
}
func (m *AckResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckResponse.Marshal(b, m, deterministic)
}
func (dst *AckResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckResponse.Merge(dst, src)
}
func (m *AckResponse) XXX_Size() int {
	return xxx_messageInfo_AckResponse.Size(m)
}
func (m *AckResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AckResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AckResponse proto.InternalMessageInfo

func (m *AckResponse) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type PeerDigest struct {
	Identity             *PeerIdentity       `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
	HeartbeatState       *PeerHeartbeatState `protobuf:"bytes,2,opt,name=heartbeat_state,json=heartbeatState,proto3" json:"heartbeat_state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *PeerDigest) Reset()         { *m = PeerDigest{} }
func (m *PeerDigest) String() string { return proto.CompactTextString(m) }
func (*PeerDigest) ProtoMessage()    {}
func (*PeerDigest) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{4}
}
func (m *PeerDigest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerDigest.Unmarshal(m, b)
}
func (m *PeerDigest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerDigest.Marshal(b, m, deterministic)
}
func (dst *PeerDigest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerDigest.Merge(dst, src)
}
func (m *PeerDigest) XXX_Size() int {
	return xxx_messageInfo_PeerDigest.Size(m)
}
func (m *PeerDigest) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerDigest.DiscardUnknown(m)
}

var xxx_messageInfo_PeerDigest proto.InternalMessageInfo

func (m *PeerDigest) GetIdentity() *PeerIdentity {
	if m != nil {
		return m.Identity
	}
	return nil
}

func (m *PeerDigest) GetHeartbeatState() *PeerHeartbeatState {
	if m != nil {
		return m.HeartbeatState
	}
	return nil
}

type PeerDiscovered struct {
	Identity             *PeerIdentity       `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
	HeartbeatState       *PeerHeartbeatState `protobuf:"bytes,2,opt,name=heartbeat_state,json=heartbeatState,proto3" json:"heartbeat_state,omitempty"`
	AppState             *PeerAppState       `protobuf:"bytes,3,opt,name=app_state,json=appState,proto3" json:"app_state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *PeerDiscovered) Reset()         { *m = PeerDiscovered{} }
func (m *PeerDiscovered) String() string { return proto.CompactTextString(m) }
func (*PeerDiscovered) ProtoMessage()    {}
func (*PeerDiscovered) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{5}
}
func (m *PeerDiscovered) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerDiscovered.Unmarshal(m, b)
}
func (m *PeerDiscovered) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerDiscovered.Marshal(b, m, deterministic)
}
func (dst *PeerDiscovered) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerDiscovered.Merge(dst, src)
}
func (m *PeerDiscovered) XXX_Size() int {
	return xxx_messageInfo_PeerDiscovered.Size(m)
}
func (m *PeerDiscovered) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerDiscovered.DiscardUnknown(m)
}

var xxx_messageInfo_PeerDiscovered proto.InternalMessageInfo

func (m *PeerDiscovered) GetIdentity() *PeerIdentity {
	if m != nil {
		return m.Identity
	}
	return nil
}

func (m *PeerDiscovered) GetHeartbeatState() *PeerHeartbeatState {
	if m != nil {
		return m.HeartbeatState
	}
	return nil
}

func (m *PeerDiscovered) GetAppState() *PeerAppState {
	if m != nil {
		return m.AppState
	}
	return nil
}

type PeerIdentity struct {
	PublicKey            string   `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Ip                   string   `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Port                 int32    `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeerIdentity) Reset()         { *m = PeerIdentity{} }
func (m *PeerIdentity) String() string { return proto.CompactTextString(m) }
func (*PeerIdentity) ProtoMessage()    {}
func (*PeerIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{6}
}
func (m *PeerIdentity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerIdentity.Unmarshal(m, b)
}
func (m *PeerIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerIdentity.Marshal(b, m, deterministic)
}
func (dst *PeerIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerIdentity.Merge(dst, src)
}
func (m *PeerIdentity) XXX_Size() int {
	return xxx_messageInfo_PeerIdentity.Size(m)
}
func (m *PeerIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_PeerIdentity proto.InternalMessageInfo

func (m *PeerIdentity) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

func (m *PeerIdentity) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *PeerIdentity) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type PeerHeartbeatState struct {
	GenerationTime       int64    `protobuf:"varint,1,opt,name=generation_time,json=generationTime,proto3" json:"generation_time,omitempty"`
	ElapsedHeartbeats    int64    `protobuf:"varint,2,opt,name=elapsed_heartbeats,json=elapsedHeartbeats,proto3" json:"elapsed_heartbeats,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeerHeartbeatState) Reset()         { *m = PeerHeartbeatState{} }
func (m *PeerHeartbeatState) String() string { return proto.CompactTextString(m) }
func (*PeerHeartbeatState) ProtoMessage()    {}
func (*PeerHeartbeatState) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{7}
}
func (m *PeerHeartbeatState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerHeartbeatState.Unmarshal(m, b)
}
func (m *PeerHeartbeatState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerHeartbeatState.Marshal(b, m, deterministic)
}
func (dst *PeerHeartbeatState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerHeartbeatState.Merge(dst, src)
}
func (m *PeerHeartbeatState) XXX_Size() int {
	return xxx_messageInfo_PeerHeartbeatState.Size(m)
}
func (m *PeerHeartbeatState) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerHeartbeatState.DiscardUnknown(m)
}

var xxx_messageInfo_PeerHeartbeatState proto.InternalMessageInfo

func (m *PeerHeartbeatState) GetGenerationTime() int64 {
	if m != nil {
		return m.GenerationTime
	}
	return 0
}

func (m *PeerHeartbeatState) GetElapsedHeartbeats() int64 {
	if m != nil {
		return m.ElapsedHeartbeats
	}
	return 0
}

type PeerAppState struct {
	Status                PeerAppState_PeerStatus      `protobuf:"varint,1,opt,name=status,proto3,enum=api.PeerAppState_PeerStatus" json:"status,omitempty"`
	CpuLoad               string                       `protobuf:"bytes,2,opt,name=cpu_load,json=cpuLoad,proto3" json:"cpu_load,omitempty"`
	FreeDiskSpace         float32                      `protobuf:"fixed32,3,opt,name=free_disk_space,json=freeDiskSpace,proto3" json:"free_disk_space,omitempty"`
	Version               string                       `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
	GeoPosition           *PeerAppState_GeoCoordinates `protobuf:"bytes,5,opt,name=geo_position,json=geoPosition,proto3" json:"geo_position,omitempty"`
	P2PFactor             int32                        `protobuf:"varint,6,opt,name=p2p_factor,json=p2pFactor,proto3" json:"p2p_factor,omitempty"`
	DiscoveredPeersNumber int32                        `protobuf:"varint,7,opt,name=discovered_peers_number,json=discoveredPeersNumber,proto3" json:"discovered_peers_number,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}                     `json:"-"`
	XXX_unrecognized      []byte                       `json:"-"`
	XXX_sizecache         int32                        `json:"-"`
}

func (m *PeerAppState) Reset()         { *m = PeerAppState{} }
func (m *PeerAppState) String() string { return proto.CompactTextString(m) }
func (*PeerAppState) ProtoMessage()    {}
func (*PeerAppState) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{8}
}
func (m *PeerAppState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerAppState.Unmarshal(m, b)
}
func (m *PeerAppState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerAppState.Marshal(b, m, deterministic)
}
func (dst *PeerAppState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerAppState.Merge(dst, src)
}
func (m *PeerAppState) XXX_Size() int {
	return xxx_messageInfo_PeerAppState.Size(m)
}
func (m *PeerAppState) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerAppState.DiscardUnknown(m)
}

var xxx_messageInfo_PeerAppState proto.InternalMessageInfo

func (m *PeerAppState) GetStatus() PeerAppState_PeerStatus {
	if m != nil {
		return m.Status
	}
	return PeerAppState_BOOTSTRAPING
}

func (m *PeerAppState) GetCpuLoad() string {
	if m != nil {
		return m.CpuLoad
	}
	return ""
}

func (m *PeerAppState) GetFreeDiskSpace() float32 {
	if m != nil {
		return m.FreeDiskSpace
	}
	return 0
}

func (m *PeerAppState) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *PeerAppState) GetGeoPosition() *PeerAppState_GeoCoordinates {
	if m != nil {
		return m.GeoPosition
	}
	return nil
}

func (m *PeerAppState) GetP2PFactor() int32 {
	if m != nil {
		return m.P2PFactor
	}
	return 0
}

func (m *PeerAppState) GetDiscoveredPeersNumber() int32 {
	if m != nil {
		return m.DiscoveredPeersNumber
	}
	return 0
}

type PeerAppState_GeoCoordinates struct {
	Latitude             float32  `protobuf:"fixed32,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude            float32  `protobuf:"fixed32,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeerAppState_GeoCoordinates) Reset()         { *m = PeerAppState_GeoCoordinates{} }
func (m *PeerAppState_GeoCoordinates) String() string { return proto.CompactTextString(m) }
func (*PeerAppState_GeoCoordinates) ProtoMessage()    {}
func (*PeerAppState_GeoCoordinates) Descriptor() ([]byte, []int) {
	return fileDescriptor_discovery_dc9a66bfee8fbc39, []int{8, 0}
}
func (m *PeerAppState_GeoCoordinates) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerAppState_GeoCoordinates.Unmarshal(m, b)
}
func (m *PeerAppState_GeoCoordinates) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerAppState_GeoCoordinates.Marshal(b, m, deterministic)
}
func (dst *PeerAppState_GeoCoordinates) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerAppState_GeoCoordinates.Merge(dst, src)
}
func (m *PeerAppState_GeoCoordinates) XXX_Size() int {
	return xxx_messageInfo_PeerAppState_GeoCoordinates.Size(m)
}
func (m *PeerAppState_GeoCoordinates) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerAppState_GeoCoordinates.DiscardUnknown(m)
}

var xxx_messageInfo_PeerAppState_GeoCoordinates proto.InternalMessageInfo

func (m *PeerAppState_GeoCoordinates) GetLatitude() float32 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *PeerAppState_GeoCoordinates) GetLongitude() float32 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

func init() {
	proto.RegisterType((*SynRequest)(nil), "api.SynRequest")
	proto.RegisterType((*SynResponse)(nil), "api.SynResponse")
	proto.RegisterType((*AckRequest)(nil), "api.AckRequest")
	proto.RegisterType((*AckResponse)(nil), "api.AckResponse")
	proto.RegisterType((*PeerDigest)(nil), "api.PeerDigest")
	proto.RegisterType((*PeerDiscovered)(nil), "api.PeerDiscovered")
	proto.RegisterType((*PeerIdentity)(nil), "api.PeerIdentity")
	proto.RegisterType((*PeerHeartbeatState)(nil), "api.PeerHeartbeatState")
	proto.RegisterType((*PeerAppState)(nil), "api.PeerAppState")
	proto.RegisterType((*PeerAppState_GeoCoordinates)(nil), "api.PeerAppState.GeoCoordinates")
	proto.RegisterEnum("api.PeerAppState_PeerStatus", PeerAppState_PeerStatus_name, PeerAppState_PeerStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DiscoveryServiceClient is the client API for DiscoveryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DiscoveryServiceClient interface {
	Synchronize(ctx context.Context, in *SynRequest, opts ...grpc.CallOption) (*SynResponse, error)
	Acknowledge(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*AckResponse, error)
}

type discoveryServiceClient struct {
	cc *grpc.ClientConn
}

func NewDiscoveryServiceClient(cc *grpc.ClientConn) DiscoveryServiceClient {
	return &discoveryServiceClient{cc}
}

func (c *discoveryServiceClient) Synchronize(ctx context.Context, in *SynRequest, opts ...grpc.CallOption) (*SynResponse, error) {
	out := new(SynResponse)
	err := c.cc.Invoke(ctx, "/api.DiscoveryService/Synchronize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discoveryServiceClient) Acknowledge(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*AckResponse, error) {
	out := new(AckResponse)
	err := c.cc.Invoke(ctx, "/api.DiscoveryService/Acknowledge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscoveryServiceServer is the server API for DiscoveryService service.
type DiscoveryServiceServer interface {
	Synchronize(context.Context, *SynRequest) (*SynResponse, error)
	Acknowledge(context.Context, *AckRequest) (*AckResponse, error)
}

func RegisterDiscoveryServiceServer(s *grpc.Server, srv DiscoveryServiceServer) {
	s.RegisterService(&_DiscoveryService_serviceDesc, srv)
}

func _DiscoveryService_Synchronize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SynRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoveryServiceServer).Synchronize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.DiscoveryService/Synchronize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoveryServiceServer).Synchronize(ctx, req.(*SynRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscoveryService_Acknowledge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoveryServiceServer).Acknowledge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.DiscoveryService/Acknowledge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoveryServiceServer).Acknowledge(ctx, req.(*AckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DiscoveryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.DiscoveryService",
	HandlerType: (*DiscoveryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Synchronize",
			Handler:    _DiscoveryService_Synchronize_Handler,
		},
		{
			MethodName: "Acknowledge",
			Handler:    _DiscoveryService_Acknowledge_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "discovery.proto",
}

func init() { proto.RegisterFile("discovery.proto", fileDescriptor_discovery_dc9a66bfee8fbc39) }

var fileDescriptor_discovery_dc9a66bfee8fbc39 = []byte{
	// 722 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x54, 0xc1, 0x4e, 0xc3, 0x46,
	0x10, 0x25, 0x0e, 0x84, 0x64, 0x42, 0x13, 0x67, 0xab, 0x0a, 0x17, 0x51, 0x09, 0xf9, 0xd0, 0x22,
	0x55, 0x44, 0x95, 0x5b, 0xf5, 0xc4, 0x81, 0x94, 0x14, 0x4a, 0x41, 0x84, 0x3a, 0xe9, 0x01, 0xa9,
	0x92, 0xb5, 0xb1, 0x87, 0x64, 0x15, 0xc7, 0xbb, 0xdd, 0xdd, 0x50, 0x05, 0xa9, 0x5f, 0xd4, 0x53,
	0xbf, 0xa4, 0xbf, 0x54, 0x79, 0xd7, 0xb1, 0x09, 0x70, 0xe8, 0xa9, 0x37, 0xef, 0xcc, 0x9b, 0x99,
	0xf7, 0x66, 0x3c, 0x03, 0xdd, 0x84, 0xa9, 0x98, 0x3f, 0xa3, 0x5c, 0xf7, 0x85, 0xe4, 0x9a, 0x93,
	0x3a, 0x15, 0xcc, 0xff, 0x0d, 0x60, 0xbc, 0xce, 0x42, 0xfc, 0x7d, 0x85, 0x4a, 0x93, 0x6f, 0xa0,
	0xbd, 0xc8, 0xf8, 0x1f, 0x59, 0x24, 0x10, 0xa5, 0xf2, 0x6a, 0x27, 0xf5, 0xd3, 0x76, 0xd0, 0xed,
	0x53, 0xc1, 0xfa, 0x0f, 0x88, 0x72, 0xc8, 0x66, 0xa8, 0x74, 0x08, 0x06, 0x93, 0x1b, 0x14, 0x39,
	0x86, 0x96, 0x66, 0x4b, 0x54, 0x9a, 0x2e, 0x85, 0xe7, 0x9c, 0xd4, 0x4e, 0xeb, 0x61, 0x65, 0xf0,
	0xff, 0xaa, 0x41, 0xdb, 0xa4, 0x57, 0x82, 0x67, 0x0a, 0xc9, 0x05, 0xf4, 0x24, 0x2e, 0xb9, 0xc6,
	0x68, 0x43, 0x86, 0x6d, 0xaa, 0x7c, 0xfa, 0xaa, 0x8a, 0x75, 0x61, 0x12, 0xba, 0x16, 0x3d, 0x2c,
	0xc1, 0xe4, 0x1c, 0x7a, 0x29, 0x8f, 0x69, 0x5a, 0x25, 0x40, 0xe5, 0x39, 0x1f, 0xf3, 0x74, 0x0d,
	0x72, 0x58, 0x01, 0xb7, 0xd9, 0xd6, 0xdf, 0xb2, 0x9d, 0x03, 0x0c, 0xe2, 0xc5, 0xa6, 0x17, 0xe7,
	0xd0, 0x95, 0xf6, 0x13, 0x93, 0xad, 0x7e, 0x7c, 0xc8, 0xb4, 0x53, 0x62, 0xff, 0x4b, 0x5f, 0xbe,
	0x86, 0xb6, 0xa9, 0x54, 0xb4, 0x65, 0x0b, 0x5c, 0x7b, 0x0b, 0xfe, 0x13, 0xa0, 0x12, 0x45, 0xce,
	0xa0, 0xc9, 0x12, 0xcc, 0x34, 0xd3, 0x6b, 0x03, 0x6d, 0x07, 0xbd, 0x92, 0xcf, 0x4d, 0xe1, 0x08,
	0x4b, 0x08, 0xb9, 0x80, 0xee, 0x1c, 0xa9, 0xd4, 0x53, 0xa4, 0x3a, 0x52, 0x9a, 0x6a, 0x34, 0x6c,
	0xda, 0xc1, 0x61, 0x19, 0xf5, 0xd3, 0xc6, 0x3f, 0xce, 0xdd, 0x61, 0x67, 0xbe, 0xf5, 0xf6, 0xff,
	0xae, 0x41, 0x67, 0x5b, 0xec, 0xff, 0xce, 0x81, 0xf4, 0xa1, 0x45, 0x85, 0x28, 0x62, 0xeb, 0x6f,
	0x2a, 0x0e, 0x84, 0xb0, 0x51, 0x4d, 0x5a, 0x7c, 0xf9, 0xbf, 0xc0, 0xc1, 0x6b, 0x2e, 0xe4, 0x0b,
	0x00, 0xb1, 0x9a, 0xa6, 0x2c, 0x8e, 0x16, 0x68, 0x29, 0xb7, 0xc2, 0x96, 0xb5, 0xdc, 0xe2, 0x9a,
	0x74, 0xc0, 0x61, 0x76, 0x4a, 0xad, 0xd0, 0x61, 0x82, 0x10, 0xd8, 0x15, 0x5c, 0x6a, 0x53, 0x69,
	0x2f, 0x34, 0xdf, 0x7e, 0x0a, 0xe4, 0x3d, 0x51, 0xf2, 0x15, 0x74, 0x67, 0x98, 0xa1, 0xa4, 0x9a,
	0xf1, 0x2c, 0xca, 0x67, 0x56, 0xcc, 0xaf, 0x53, 0x99, 0x27, 0x6c, 0x89, 0xe4, 0x0c, 0x08, 0xa6,
	0x54, 0x28, 0x4c, 0xa2, 0x52, 0x9b, 0x2a, 0x7e, 0x8c, 0x5e, 0xe1, 0x29, 0x73, 0x2b, 0xff, 0x9f,
	0xba, 0x55, 0xb0, 0xd1, 0x46, 0xbe, 0x83, 0x46, 0xae, 0x7e, 0xa5, 0x4c, 0xfe, 0x4e, 0x70, 0xfc,
	0x4e, 0xbe, 0x79, 0x8c, 0x0d, 0x26, 0x2c, 0xb0, 0xe4, 0x73, 0x68, 0xc6, 0x62, 0x15, 0xa5, 0x9c,
	0x26, 0x85, 0xbc, 0xfd, 0x58, 0xac, 0xee, 0x38, 0x4d, 0xc8, 0x97, 0xd0, 0x7d, 0x92, 0x68, 0x16,
	0x71, 0x11, 0x29, 0x41, 0x63, 0xdb, 0x58, 0x27, 0xfc, 0x24, 0x37, 0x0f, 0x99, 0x5a, 0x8c, 0x73,
	0x23, 0xf1, 0x60, 0xff, 0x19, 0xa5, 0x62, 0x3c, 0xf3, 0x76, 0x6d, 0x86, 0xe2, 0x49, 0x2e, 0xe1,
	0x60, 0x86, 0x3c, 0x12, 0x5c, 0xb1, 0x5c, 0xa6, 0xb7, 0x67, 0xe6, 0x72, 0xf2, 0x9e, 0xd8, 0x35,
	0xf2, 0x4b, 0xce, 0x65, 0xc2, 0x32, 0xaa, 0x51, 0x85, 0xed, 0x19, 0xf2, 0x87, 0x22, 0xc8, 0x4c,
	0x26, 0x10, 0xd1, 0x13, 0x8d, 0x35, 0x97, 0x5e, 0xc3, 0x34, 0xbc, 0x25, 0x02, 0x71, 0x65, 0x0c,
	0xe4, 0x7b, 0x38, 0x4c, 0xca, 0xff, 0xce, 0x6e, 0x61, 0x94, 0xad, 0x96, 0x53, 0x94, 0xde, 0xbe,
	0xc1, 0x7e, 0x56, 0xb9, 0xcd, 0xe2, 0xdd, 0x1b, 0xe7, 0xd1, 0xcf, 0xd0, 0xd9, 0xae, 0x4a, 0x8e,
	0xa0, 0x99, 0x52, 0xcd, 0xf4, 0x2a, 0xb1, 0x23, 0x72, 0xc2, 0xf2, 0x9d, 0xef, 0x5f, 0xca, 0xb3,
	0x99, 0x75, 0x3a, 0xc6, 0x59, 0x19, 0xfc, 0xa1, 0xdd, 0x3f, 0xdb, 0x5a, 0xe2, 0xc2, 0xc1, 0x0f,
	0xa3, 0xd1, 0x64, 0x3c, 0x09, 0x07, 0x0f, 0x37, 0xf7, 0xd7, 0xee, 0x0e, 0x69, 0x80, 0x33, 0xba,
	0x75, 0x6b, 0x04, 0xa0, 0x71, 0x35, 0xf8, 0xf5, 0x6e, 0xf2, 0xe8, 0x3a, 0x39, 0x6a, 0x3c, 0x19,
	0x85, 0x83, 0xeb, 0x1f, 0xa3, 0xd1, 0xfd, 0xdd, 0xa3, 0x5b, 0x0f, 0x5e, 0xc0, 0xdd, 0x6c, 0xd0,
	0x7a, 0x8c, 0xf2, 0x99, 0xc5, 0x48, 0x02, 0x73, 0x1d, 0xe3, 0xb9, 0xe4, 0x19, 0x7b, 0x41, 0x62,
	0x0f, 0x58, 0x75, 0x8e, 0x8f, 0xdc, 0xca, 0x60, 0x2f, 0x85, 0xbf, 0x93, 0xc7, 0x0c, 0xe2, 0xfc,
	0x00, 0xa7, 0x98, 0xcc, 0x36, 0x31, 0xd5, 0xd9, 0x2a, 0x62, 0x5e, 0x5d, 0x17, 0x7f, 0x67, 0xda,
	0x30, 0x07, 0xff, 0xdb, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x9f, 0xb6, 0x6a, 0x03, 0x06,
	0x00, 0x00,
}

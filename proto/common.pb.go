// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type HealthCheckRes_ServingStatus int32

const (
	HealthCheckRes_UNKNOWN     HealthCheckRes_ServingStatus = 0
	HealthCheckRes_SERVING     HealthCheckRes_ServingStatus = 1
	HealthCheckRes_NOT_SERVING HealthCheckRes_ServingStatus = 2
)

var HealthCheckRes_ServingStatus_name = map[int32]string{
	0: "UNKNOWN",
	1: "SERVING",
	2: "NOT_SERVING",
}

var HealthCheckRes_ServingStatus_value = map[string]int32{
	"UNKNOWN":     0,
	"SERVING":     1,
	"NOT_SERVING": 2,
}

func (x HealthCheckRes_ServingStatus) String() string {
	return proto.EnumName(HealthCheckRes_ServingStatus_name, int32(x))
}

func (HealthCheckRes_ServingStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{1, 0}
}

type HealthCheckReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthCheckReq) Reset()         { *m = HealthCheckReq{} }
func (m *HealthCheckReq) String() string { return proto.CompactTextString(m) }
func (*HealthCheckReq) ProtoMessage()    {}
func (*HealthCheckReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

func (m *HealthCheckReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthCheckReq.Unmarshal(m, b)
}
func (m *HealthCheckReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthCheckReq.Marshal(b, m, deterministic)
}
func (m *HealthCheckReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckReq.Merge(m, src)
}
func (m *HealthCheckReq) XXX_Size() int {
	return xxx_messageInfo_HealthCheckReq.Size(m)
}
func (m *HealthCheckReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckReq.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckReq proto.InternalMessageInfo

type HealthCheckRes struct {
	Status               HealthCheckRes_ServingStatus `protobuf:"varint,1,opt,name=status,proto3,enum=pb.HealthCheckRes_ServingStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *HealthCheckRes) Reset()         { *m = HealthCheckRes{} }
func (m *HealthCheckRes) String() string { return proto.CompactTextString(m) }
func (*HealthCheckRes) ProtoMessage()    {}
func (*HealthCheckRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{1}
}

func (m *HealthCheckRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthCheckRes.Unmarshal(m, b)
}
func (m *HealthCheckRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthCheckRes.Marshal(b, m, deterministic)
}
func (m *HealthCheckRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckRes.Merge(m, src)
}
func (m *HealthCheckRes) XXX_Size() int {
	return xxx_messageInfo_HealthCheckRes.Size(m)
}
func (m *HealthCheckRes) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckRes.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckRes proto.InternalMessageInfo

func (m *HealthCheckRes) GetStatus() HealthCheckRes_ServingStatus {
	if m != nil {
		return m.Status
	}
	return HealthCheckRes_UNKNOWN
}

func init() {
	proto.RegisterEnum("pb.HealthCheckRes_ServingStatus", HealthCheckRes_ServingStatus_name, HealthCheckRes_ServingStatus_value)
	proto.RegisterType((*HealthCheckReq)(nil), "pb.HealthCheckReq")
	proto.RegisterType((*HealthCheckRes)(nil), "pb.HealthCheckRes")
}

func init() { proto.RegisterFile("common.proto", fileDescriptor_555bd8c177793206) }

var fileDescriptor_555bd8c177793206 = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x12, 0xe0, 0xe2,
	0xf3, 0x48, 0x4d, 0xcc, 0x29, 0xc9, 0x70, 0xce, 0x48, 0x4d, 0xce, 0x0e, 0x4a, 0x2d, 0x54, 0x6a,
	0x63, 0x44, 0x13, 0x2a, 0x16, 0xb2, 0xe0, 0x62, 0x2b, 0x2e, 0x49, 0x2c, 0x29, 0x2d, 0x96, 0x60,
	0x54, 0x60, 0xd4, 0xe0, 0x33, 0x52, 0xd0, 0x2b, 0x48, 0xd2, 0x43, 0x55, 0xa3, 0x17, 0x9c, 0x5a,
	0x54, 0x96, 0x99, 0x97, 0x1e, 0x0c, 0x56, 0x17, 0x04, 0x55, 0xaf, 0x64, 0xc5, 0xc5, 0x8b, 0x22,
	0x21, 0xc4, 0xcd, 0xc5, 0x1e, 0xea, 0xe7, 0xed, 0xe7, 0x1f, 0xee, 0x27, 0xc0, 0x00, 0xe2, 0x04,
	0xbb, 0x06, 0x85, 0x79, 0xfa, 0xb9, 0x0b, 0x30, 0x0a, 0xf1, 0x73, 0x71, 0xfb, 0xf9, 0x87, 0xc4,
	0xc3, 0x04, 0x98, 0x92, 0xd8, 0xc0, 0xae, 0x34, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x0b, 0x3c,
	0xfa, 0x70, 0xb5, 0x00, 0x00, 0x00,
}
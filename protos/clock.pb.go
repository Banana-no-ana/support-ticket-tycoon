// Code generated by protoc-gen-go. DO NOT EDIT.
// source: clock.proto

package protos // import "github.com/Banana-no-ana/support-ticket-tycoon/protos"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type WorkerRegister struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WorkerRegister) Reset()         { *m = WorkerRegister{} }
func (m *WorkerRegister) String() string { return proto.CompactTextString(m) }
func (*WorkerRegister) ProtoMessage()    {}
func (*WorkerRegister) Descriptor() ([]byte, []int) {
	return fileDescriptor_clock_ec56df741038c7d7, []int{0}
}
func (m *WorkerRegister) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkerRegister.Unmarshal(m, b)
}
func (m *WorkerRegister) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkerRegister.Marshal(b, m, deterministic)
}
func (dst *WorkerRegister) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkerRegister.Merge(dst, src)
}
func (m *WorkerRegister) XXX_Size() int {
	return xxx_messageInfo_WorkerRegister.Size(m)
}
func (m *WorkerRegister) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkerRegister.DiscardUnknown(m)
}

var xxx_messageInfo_WorkerRegister proto.InternalMessageInfo

func (m *WorkerRegister) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

type Tick struct {
	TickNum              int32    `protobuf:"varint,1,opt,name=TickNum,proto3" json:"TickNum,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tick) Reset()         { *m = Tick{} }
func (m *Tick) String() string { return proto.CompactTextString(m) }
func (*Tick) ProtoMessage()    {}
func (*Tick) Descriptor() ([]byte, []int) {
	return fileDescriptor_clock_ec56df741038c7d7, []int{1}
}
func (m *Tick) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tick.Unmarshal(m, b)
}
func (m *Tick) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tick.Marshal(b, m, deterministic)
}
func (dst *Tick) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tick.Merge(dst, src)
}
func (m *Tick) XXX_Size() int {
	return xxx_messageInfo_Tick.Size(m)
}
func (m *Tick) XXX_DiscardUnknown() {
	xxx_messageInfo_Tick.DiscardUnknown(m)
}

var xxx_messageInfo_Tick proto.InternalMessageInfo

func (m *Tick) GetTickNum() int32 {
	if m != nil {
		return m.TickNum
	}
	return 0
}

func init() {
	proto.RegisterType((*WorkerRegister)(nil), "clock.WorkerRegister")
	proto.RegisterType((*Tick)(nil), "clock.Tick")
}

func init() { proto.RegisterFile("clock.proto", fileDescriptor_clock_ec56df741038c7d7) }

var fileDescriptor_clock_ec56df741038c7d7 = []byte{
	// 183 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0xce, 0xc9, 0x4f,
	0xce, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x14, 0xb8, 0xf8, 0xc2,
	0xf3, 0x8b, 0xb2, 0x53, 0x8b, 0x82, 0x52, 0xd3, 0x33, 0x8b, 0x4b, 0x52, 0x8b, 0x84, 0xf8, 0xb8,
	0x98, 0x3c, 0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x98, 0x3c, 0x5d, 0x94, 0x14, 0xb8,
	0x58, 0x42, 0x32, 0x93, 0xb3, 0x85, 0x24, 0xb8, 0xd8, 0x41, 0xb4, 0x5f, 0x69, 0x2e, 0x58, 0x92,
	0x35, 0x08, 0xc6, 0x35, 0xb2, 0xe4, 0x62, 0x75, 0x06, 0x19, 0x26, 0x64, 0xc0, 0xc5, 0x01, 0x37,
	0x46, 0x54, 0x0f, 0x62, 0x1b, 0xaa, 0xe9, 0x52, 0xdc, 0x50, 0x61, 0x90, 0x5e, 0x03, 0x46, 0x27,
	0xf3, 0x28, 0xd3, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0xa7, 0xc4,
	0xbc, 0xc4, 0xbc, 0x44, 0xdd, 0xbc, 0x7c, 0xdd, 0xc4, 0xbc, 0x44, 0xfd, 0xe2, 0xd2, 0x82, 0x82,
	0xfc, 0xa2, 0x12, 0xdd, 0x92, 0xcc, 0xe4, 0xec, 0xd4, 0x12, 0xdd, 0x92, 0xca, 0xe4, 0xfc, 0xfc,
	0x3c, 0x7d, 0xb0, 0xeb, 0x8b, 0x93, 0xd8, 0xc0, 0xb4, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xc4,
	0xc4, 0x50, 0x41, 0xd4, 0x00, 0x00, 0x00,
}

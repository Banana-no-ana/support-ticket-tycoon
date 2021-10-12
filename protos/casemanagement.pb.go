// Code generated by protoc-gen-go. DO NOT EDIT.
// source: casemanagement.proto

package casemanagement // import "github.com/Banana-no-ana/support-ticket-tycoon/protos/casemanagement"

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

type Case struct {
	CaseID               int32    `protobuf:"varint,1,opt,name=caseID,proto3" json:"caseID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Case) Reset()         { *m = Case{} }
func (m *Case) String() string { return proto.CompactTextString(m) }
func (*Case) ProtoMessage()    {}
func (*Case) Descriptor() ([]byte, []int) {
	return fileDescriptor_casemanagement_eefb5d921c748e9f, []int{0}
}
func (m *Case) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Case.Unmarshal(m, b)
}
func (m *Case) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Case.Marshal(b, m, deterministic)
}
func (dst *Case) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Case.Merge(dst, src)
}
func (m *Case) XXX_Size() int {
	return xxx_messageInfo_Case.Size(m)
}
func (m *Case) XXX_DiscardUnknown() {
	xxx_messageInfo_Case.DiscardUnknown(m)
}

var xxx_messageInfo_Case proto.InternalMessageInfo

func (m *Case) GetCaseID() int32 {
	if m != nil {
		return m.CaseID
	}
	return 0
}

type Response struct {
	Status               bool     `protobuf:"varint,1,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_casemanagement_eefb5d921c748e9f, []int{1}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func init() {
	proto.RegisterType((*Case)(nil), "casemanagement.Case")
	proto.RegisterType((*Response)(nil), "casemanagement.Response")
}

func init() {
	proto.RegisterFile("casemanagement.proto", fileDescriptor_casemanagement_eefb5d921c748e9f)
}

var fileDescriptor_casemanagement_eefb5d921c748e9f = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x49, 0x4e, 0x2c, 0x4e,
	0xcd, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f, 0xcd, 0x4d, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x43, 0x15, 0x55, 0x92, 0xe3, 0x62, 0x71, 0x4e, 0x2c, 0x4e, 0x15, 0x12, 0xe3, 0x62,
	0x03, 0xc9, 0x78, 0xba, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x06, 0x41, 0x79, 0x4a, 0x4a, 0x5c,
	0x1c, 0x41, 0xa9, 0xc5, 0x05, 0xf9, 0x79, 0x10, 0x35, 0xc1, 0x25, 0x89, 0x25, 0xa5, 0xc5, 0x60,
	0x35, 0x1c, 0x41, 0x50, 0x9e, 0x51, 0x1d, 0x17, 0x9b, 0x63, 0x71, 0x71, 0x66, 0x7a, 0x9e, 0x90,
	0x05, 0x9c, 0x25, 0xa2, 0x87, 0x66, 0x3d, 0xc8, 0x16, 0x29, 0x09, 0x74, 0x51, 0xb8, 0xd9, 0x56,
	0x5c, 0x1c, 0xa1, 0x79, 0xe4, 0xe9, 0x75, 0x72, 0x8b, 0x72, 0x49, 0xcf, 0x2c, 0xc9, 0x28, 0x4d,
	0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x77, 0x4a, 0xcc, 0x4b, 0xcc, 0x4b, 0xd4, 0xcd, 0xcb, 0xd7, 0x4d,
	0xcc, 0x4b, 0xd4, 0x2f, 0x2e, 0x2d, 0x28, 0xc8, 0x2f, 0x2a, 0xd1, 0x2d, 0xc9, 0x4c, 0xce, 0x4e,
	0x2d, 0xd1, 0x2d, 0xa9, 0x4c, 0xce, 0xcf, 0xcf, 0xd3, 0x07, 0x87, 0x45, 0xb1, 0x3e, 0xaa, 0x81,
	0x49, 0x6c, 0x60, 0x61, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xdf, 0x4e, 0x76, 0x26, 0x3a,
	0x01, 0x00, 0x00,
}

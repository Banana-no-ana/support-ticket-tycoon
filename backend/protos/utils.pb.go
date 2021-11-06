// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/utils.proto

package protos // import "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"

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

type SkillEnum int32

const (
	SkillEnum_Troubleshoot SkillEnum = 0
	SkillEnum_Build        SkillEnum = 1
	SkillEnum_Tech         SkillEnum = 2
	SkillEnum_Usage        SkillEnum = 3
	SkillEnum_Architecture SkillEnum = 4
	SkillEnum_Environment  SkillEnum = 5
	SkillEnum_Explain      SkillEnum = 6
	SkillEnum_Empathy      SkillEnum = 7
	SkillEnum_Relationship SkillEnum = 8
)

var SkillEnum_name = map[int32]string{
	0: "Troubleshoot",
	1: "Build",
	2: "Tech",
	3: "Usage",
	4: "Architecture",
	5: "Environment",
	6: "Explain",
	7: "Empathy",
	8: "Relationship",
}
var SkillEnum_value = map[string]int32{
	"Troubleshoot": 0,
	"Build":        1,
	"Tech":         2,
	"Usage":        3,
	"Architecture": 4,
	"Environment":  5,
	"Explain":      6,
	"Empathy":      7,
	"Relationship": 8,
}

func (x SkillEnum) String() string {
	return proto.EnumName(SkillEnum_name, int32(x))
}
func (SkillEnum) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_utils_4306bbb4a25bade7, []int{0}
}

type Response struct {
	Success              bool     `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_utils_4306bbb4a25bade7, []int{0}
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

func (m *Response) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

// Set scenario difficulty
type Difficulty struct {
	MinDifficulty        int32    `protobuf:"varint,1,opt,name=minDifficulty,proto3" json:"minDifficulty,omitempty"`
	MaxDifficulty        int32    `protobuf:"varint,2,opt,name=maxDifficulty,proto3" json:"maxDifficulty,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Difficulty) Reset()         { *m = Difficulty{} }
func (m *Difficulty) String() string { return proto.CompactTextString(m) }
func (*Difficulty) ProtoMessage()    {}
func (*Difficulty) Descriptor() ([]byte, []int) {
	return fileDescriptor_utils_4306bbb4a25bade7, []int{1}
}
func (m *Difficulty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Difficulty.Unmarshal(m, b)
}
func (m *Difficulty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Difficulty.Marshal(b, m, deterministic)
}
func (dst *Difficulty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Difficulty.Merge(dst, src)
}
func (m *Difficulty) XXX_Size() int {
	return xxx_messageInfo_Difficulty.Size(m)
}
func (m *Difficulty) XXX_DiscardUnknown() {
	xxx_messageInfo_Difficulty.DiscardUnknown(m)
}

var xxx_messageInfo_Difficulty proto.InternalMessageInfo

func (m *Difficulty) GetMinDifficulty() int32 {
	if m != nil {
		return m.MinDifficulty
	}
	return 0
}

func (m *Difficulty) GetMaxDifficulty() int32 {
	if m != nil {
		return m.MaxDifficulty
	}
	return 0
}

func init() {
	proto.RegisterType((*Response)(nil), "utils.Response")
	proto.RegisterType((*Difficulty)(nil), "utils.Difficulty")
	proto.RegisterEnum("utils.SkillEnum", SkillEnum_name, SkillEnum_value)
}

func init() { proto.RegisterFile("protos/utils.proto", fileDescriptor_utils_4306bbb4a25bade7) }

var fileDescriptor_utils_4306bbb4a25bade7 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0xc1, 0x4e, 0xc2, 0x40,
	0x10, 0x86, 0x2d, 0x52, 0x28, 0x8b, 0xc6, 0xcd, 0x9e, 0x38, 0x1a, 0xc2, 0xc1, 0x98, 0x94, 0x1e,
	0x3c, 0x1b, 0x23, 0x91, 0x17, 0x28, 0x98, 0x18, 0x6f, 0xdb, 0x65, 0xa1, 0x13, 0xb6, 0x33, 0x9b,
	0xee, 0xac, 0x81, 0x77, 0xf0, 0xa1, 0x0d, 0x55, 0x23, 0x9e, 0x26, 0xdf, 0x97, 0xef, 0xbf, 0x8c,
	0x50, 0xbe, 0x25, 0xa6, 0x50, 0x44, 0x06, 0x17, 0xe6, 0x1d, 0xa8, 0xb4, 0x83, 0xe9, 0x4c, 0x64,
	0xa5, 0x0d, 0x9e, 0x30, 0x58, 0x35, 0x11, 0xc3, 0x55, 0x34, 0xc6, 0x86, 0x30, 0x49, 0x6e, 0x93,
	0xbb, 0xac, 0xfc, 0xc5, 0xe9, 0x9b, 0x10, 0x2f, 0xb0, 0xdd, 0x82, 0x89, 0x8e, 0x8f, 0x6a, 0x26,
	0xae, 0x1b, 0xc0, 0x3f, 0xd1, 0xd5, 0x69, 0xf9, 0x5f, 0x76, 0x95, 0x3e, 0x9c, 0x55, 0xbd, 0x9f,
	0xea, 0x5c, 0xde, 0x7f, 0x26, 0x62, 0xb4, 0xda, 0x83, 0x73, 0x4b, 0x8c, 0x8d, 0x92, 0xe2, 0x6a,
	0xdd, 0x52, 0xac, 0x9c, 0x0d, 0x35, 0x11, 0xcb, 0x0b, 0x35, 0x12, 0xe9, 0x22, 0x82, 0xdb, 0xc8,
	0x44, 0x65, 0xa2, 0xbf, 0xb6, 0xa6, 0x96, 0xbd, 0x93, 0x7c, 0x0d, 0x7a, 0x67, 0xe5, 0xe5, 0x69,
	0xf1, 0xdc, 0x9a, 0x1a, 0xd8, 0x1a, 0x8e, 0xad, 0x95, 0x7d, 0x75, 0x23, 0xc6, 0x4b, 0xfc, 0x80,
	0x96, 0xb0, 0xb1, 0xc8, 0x32, 0x55, 0x63, 0x31, 0x5c, 0x1e, 0xbc, 0xd3, 0x80, 0x72, 0xd0, 0x41,
	0xe3, 0x35, 0xd7, 0x47, 0x39, 0x3c, 0x8d, 0x4b, 0xeb, 0x34, 0x03, 0x61, 0xa8, 0xc1, 0xcb, 0x6c,
	0xf1, 0xf4, 0xfe, 0xb8, 0x03, 0xae, 0x63, 0x35, 0x37, 0xd4, 0x14, 0x0b, 0x8d, 0x1a, 0x75, 0x8e,
	0x94, 0x6b, 0xd4, 0x45, 0x88, 0xde, 0x53, 0xcb, 0x39, 0x83, 0xd9, 0x5b, 0xce, 0xf9, 0x68, 0x88,
	0xb0, 0xa8, 0xb4, 0xd9, 0x5b, 0xdc, 0x14, 0xdf, 0x2f, 0xae, 0x06, 0xdd, 0x7d, 0xf8, 0x0a, 0x00,
	0x00, 0xff, 0xff, 0x6e, 0x0b, 0xa2, 0xd6, 0x73, 0x01, 0x00, 0x00,
}

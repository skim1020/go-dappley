// Code generated by protoc-gen-go. DO NOT EDIT.
// source: network/pb/dapmsg.proto

package networkpb

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

type Dapmsg struct {
	Cmd                  string   `protobuf:"bytes,1,opt,name=Cmd,proto3" json:"Cmd,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Dapmsg) Reset()         { *m = Dapmsg{} }
func (m *Dapmsg) String() string { return proto.CompactTextString(m) }
func (*Dapmsg) ProtoMessage()    {}
func (*Dapmsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_dapmsg_848b211b2fd09e9f, []int{0}
}
func (m *Dapmsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Dapmsg.Unmarshal(m, b)
}
func (m *Dapmsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Dapmsg.Marshal(b, m, deterministic)
}
func (dst *Dapmsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dapmsg.Merge(dst, src)
}
func (m *Dapmsg) XXX_Size() int {
	return xxx_messageInfo_Dapmsg.Size(m)
}
func (m *Dapmsg) XXX_DiscardUnknown() {
	xxx_messageInfo_Dapmsg.DiscardUnknown(m)
}

var xxx_messageInfo_Dapmsg proto.InternalMessageInfo

func (m *Dapmsg) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *Dapmsg) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Dapmsg)(nil), "networkpb.Dapmsg")
}

func init() { proto.RegisterFile("network/pb/dapmsg.proto", fileDescriptor_dapmsg_848b211b2fd09e9f) }

var fileDescriptor_dapmsg_848b211b2fd09e9f = []byte{
	// 99 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcf, 0x4b, 0x2d, 0x29,
	0xcf, 0x2f, 0xca, 0xd6, 0x2f, 0x48, 0xd2, 0x4f, 0x49, 0x2c, 0xc8, 0x2d, 0x4e, 0xd7, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x84, 0x4a, 0x14, 0x24, 0x29, 0xe9, 0x71, 0xb1, 0xb9, 0x80, 0xa5,
	0x84, 0x04, 0xb8, 0x98, 0x9d, 0x73, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x40, 0x4c,
	0x21, 0x21, 0x2e, 0x96, 0x94, 0xc4, 0x92, 0x44, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x30,
	0x3b, 0x89, 0x0d, 0x6c, 0x82, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xcb, 0x54, 0xf7, 0x22, 0x5c,
	0x00, 0x00, 0x00,
}

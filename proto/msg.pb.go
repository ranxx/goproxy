// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg.proto

package proto

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

type NetworkType int32

const (
	NetworkType_HTTP NetworkType = 0
	NetworkType_TCP  NetworkType = 1
)

var NetworkType_name = map[int32]string{
	0: "HTTP",
	1: "TCP",
}

var NetworkType_value = map[string]int32{
	"HTTP": 0,
	"TCP":  1,
}

func (x NetworkType) String() string {
	return proto.EnumName(NetworkType_name, int32(x))
}

func (NetworkType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

type Msg struct {
	Network              string   `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	MsgId                int64    `protobuf:"varint,2,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	Body                 []byte   `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Msg) Reset()         { *m = Msg{} }
func (m *Msg) String() string { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()    {}
func (*Msg) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

func (m *Msg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Msg.Unmarshal(m, b)
}
func (m *Msg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Msg.Marshal(b, m, deterministic)
}
func (m *Msg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Msg.Merge(m, src)
}
func (m *Msg) XXX_Size() int {
	return xxx_messageInfo_Msg.Size(m)
}
func (m *Msg) XXX_DiscardUnknown() {
	xxx_messageInfo_Msg.DiscardUnknown(m)
}

var xxx_messageInfo_Msg proto.InternalMessageInfo

func (m *Msg) GetNetwork() string {
	if m != nil {
		return m.Network
	}
	return ""
}

func (m *Msg) GetMsgId() int64 {
	if m != nil {
		return m.MsgId
	}
	return 0
}

func (m *Msg) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type Addr struct {
	Ip                   string   `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port                 int32    `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Addr) Reset()         { *m = Addr{} }
func (m *Addr) String() string { return proto.CompactTextString(m) }
func (*Addr) ProtoMessage()    {}
func (*Addr) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{1}
}

func (m *Addr) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Addr.Unmarshal(m, b)
}
func (m *Addr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Addr.Marshal(b, m, deterministic)
}
func (m *Addr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Addr.Merge(m, src)
}
func (m *Addr) XXX_Size() int {
	return xxx_messageInfo_Addr.Size(m)
}
func (m *Addr) XXX_DiscardUnknown() {
	xxx_messageInfo_Addr.DiscardUnknown(m)
}

var xxx_messageInfo_Addr proto.InternalMessageInfo

func (m *Addr) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *Addr) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type Header struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []string `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{2}
}

func (m *Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Header.Unmarshal(m, b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Header.Marshal(b, m, deterministic)
}
func (m *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(m, src)
}
func (m *Header) XXX_Size() int {
	return xxx_messageInfo_Header.Size(m)
}
func (m *Header) XXX_DiscardUnknown() {
	xxx_messageInfo_Header.DiscardUnknown(m)
}

var xxx_messageInfo_Header proto.InternalMessageInfo

func (m *Header) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Header) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

type HTTPBody struct {
	MsgId                int64     `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	Laddr                *Addr     `protobuf:"bytes,2,opt,name=Laddr,proto3" json:"Laddr,omitempty"`
	Url                  string    `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Method               string    `protobuf:"bytes,4,opt,name=Method,proto3" json:"Method,omitempty"`
	Body                 []byte    `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	Header               []*Header `protobuf:"bytes,6,rep,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *HTTPBody) Reset()         { *m = HTTPBody{} }
func (m *HTTPBody) String() string { return proto.CompactTextString(m) }
func (*HTTPBody) ProtoMessage()    {}
func (*HTTPBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{3}
}

func (m *HTTPBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTTPBody.Unmarshal(m, b)
}
func (m *HTTPBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTTPBody.Marshal(b, m, deterministic)
}
func (m *HTTPBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTTPBody.Merge(m, src)
}
func (m *HTTPBody) XXX_Size() int {
	return xxx_messageInfo_HTTPBody.Size(m)
}
func (m *HTTPBody) XXX_DiscardUnknown() {
	xxx_messageInfo_HTTPBody.DiscardUnknown(m)
}

var xxx_messageInfo_HTTPBody proto.InternalMessageInfo

func (m *HTTPBody) GetMsgId() int64 {
	if m != nil {
		return m.MsgId
	}
	return 0
}

func (m *HTTPBody) GetLaddr() *Addr {
	if m != nil {
		return m.Laddr
	}
	return nil
}

func (m *HTTPBody) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *HTTPBody) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *HTTPBody) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *HTTPBody) GetHeader() []*Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type TCPBody struct {
	MsgId                int64    `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	Laddr                *Addr    `protobuf:"bytes,2,opt,name=Laddr,proto3" json:"Laddr,omitempty"`
	Raddr                *Addr    `protobuf:"bytes,3,opt,name=Raddr,proto3" json:"Raddr,omitempty"`
	Body                 []byte   `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TCPBody) Reset()         { *m = TCPBody{} }
func (m *TCPBody) String() string { return proto.CompactTextString(m) }
func (*TCPBody) ProtoMessage()    {}
func (*TCPBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{4}
}

func (m *TCPBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TCPBody.Unmarshal(m, b)
}
func (m *TCPBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TCPBody.Marshal(b, m, deterministic)
}
func (m *TCPBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TCPBody.Merge(m, src)
}
func (m *TCPBody) XXX_Size() int {
	return xxx_messageInfo_TCPBody.Size(m)
}
func (m *TCPBody) XXX_DiscardUnknown() {
	xxx_messageInfo_TCPBody.DiscardUnknown(m)
}

var xxx_messageInfo_TCPBody proto.InternalMessageInfo

func (m *TCPBody) GetMsgId() int64 {
	if m != nil {
		return m.MsgId
	}
	return 0
}

func (m *TCPBody) GetLaddr() *Addr {
	if m != nil {
		return m.Laddr
	}
	return nil
}

func (m *TCPBody) GetRaddr() *Addr {
	if m != nil {
		return m.Raddr
	}
	return nil
}

func (m *TCPBody) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto.RegisterEnum("proto.NetworkType", NetworkType_name, NetworkType_value)
	proto.RegisterType((*Msg)(nil), "proto.Msg")
	proto.RegisterType((*Addr)(nil), "proto.Addr")
	proto.RegisterType((*Header)(nil), "proto.Header")
	proto.RegisterType((*HTTPBody)(nil), "proto.HTTPBody")
	proto.RegisterType((*TCPBody)(nil), "proto.TCPBody")
}

func init() {
	proto.RegisterFile("msg.proto", fileDescriptor_c06e4cca6c2cc899)
}

var fileDescriptor_c06e4cca6c2cc899 = []byte{
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0xc7, 0x4d, 0xd3, 0x74, 0xeb, 0xab, 0xca, 0x78, 0xa8, 0xe4, 0xd8, 0x15, 0x84, 0xb2, 0xc3,
	0x90, 0xf9, 0x09, 0x74, 0x97, 0x29, 0x4e, 0x46, 0xe8, 0x5d, 0x36, 0x12, 0xba, 0xb1, 0xd5, 0x94,
	0xb4, 0x53, 0xea, 0x37, 0xf2, 0x5b, 0x4a, 0xd2, 0xba, 0xed, 0xe0, 0xcd, 0x53, 0xfe, 0x2f, 0xef,
	0xcf, 0x7b, 0xbf, 0x7f, 0x02, 0x61, 0x51, 0xe5, 0xe3, 0xd2, 0xe8, 0x5a, 0x23, 0x73, 0x47, 0xf2,
	0x0c, 0x74, 0x5e, 0xe5, 0xc8, 0xa1, 0xf7, 0xae, 0xea, 0x4f, 0x6d, 0xb6, 0x9c, 0xc4, 0x24, 0x0d,
	0xc5, 0x6f, 0x89, 0xd7, 0x10, 0x14, 0x55, 0xfe, 0xb6, 0x91, 0xdc, 0x8b, 0x49, 0x4a, 0x05, 0x2b,
	0xaa, 0xfc, 0x49, 0x22, 0x82, 0xbf, 0xd2, 0xb2, 0xe1, 0x34, 0x26, 0xe9, 0xb9, 0x70, 0x3a, 0x19,
	0x81, 0xff, 0x20, 0xa5, 0xc1, 0x4b, 0xf0, 0x36, 0x65, 0x37, 0xc7, 0xdb, 0x94, 0xd6, 0x5b, 0x6a,
	0x53, 0xbb, 0x01, 0x4c, 0x38, 0x9d, 0xdc, 0x41, 0x30, 0x53, 0x4b, 0xa9, 0x0c, 0x0e, 0x80, 0x6e,
	0x55, 0xd3, 0xd9, 0xad, 0xc4, 0x2b, 0x60, 0x1f, 0xcb, 0xdd, 0x5e, 0x71, 0x2f, 0xa6, 0x69, 0x28,
	0xda, 0x22, 0xf9, 0x26, 0xd0, 0x9f, 0x65, 0xd9, 0xe2, 0x51, 0xcb, 0xe6, 0x84, 0x8a, 0x9c, 0x52,
	0x0d, 0x81, 0xbd, 0x2c, 0xa5, 0x34, 0x6e, 0x55, 0x34, 0x89, 0xda, 0xac, 0x63, 0x4b, 0x25, 0xda,
	0x8e, 0x5d, 0xb7, 0x37, 0x3b, 0xc7, 0x1d, 0x0a, 0x2b, 0xf1, 0x06, 0x82, 0xb9, 0xaa, 0xd7, 0x5a,
	0x72, 0xdf, 0x5d, 0x76, 0xd5, 0x21, 0x22, 0x3b, 0x46, 0xc4, 0x5b, 0x08, 0xd6, 0x0e, 0x9b, 0x07,
	0x31, 0x4d, 0xa3, 0xc9, 0x45, 0xb7, 0xa1, 0xcd, 0x22, 0xba, 0x66, 0xf2, 0x05, 0xbd, 0x6c, 0xfa,
	0x5f, 0xd2, 0x21, 0x30, 0xe1, 0x2c, 0xf4, 0x0f, 0x8b, 0xeb, 0x1c, 0x10, 0xfd, 0x23, 0xe2, 0x28,
	0x86, 0xe8, 0xb5, 0xfd, 0xbb, 0xac, 0x29, 0x15, 0xf6, 0xc1, 0xb7, 0xaf, 0x36, 0x38, 0xc3, 0x1e,
	0xd0, 0x6c, 0xba, 0x18, 0x90, 0x55, 0xe0, 0x06, 0xdd, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x7d,
	0x3c, 0x3e, 0x2c, 0x0e, 0x02, 0x00, 0x00,
}

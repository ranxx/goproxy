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

type MsgType int32

const (
	// 消息
	MsgType_None MsgType = 0
	// 错误
	MsgType_Error MsgType = 1
	// 心跳
	MsgType_Heartbeat MsgType = 3
)

var MsgType_name = map[int32]string{
	0: "None",
	1: "Error",
	3: "Heartbeat",
}

var MsgType_value = map[string]int32{
	"None":      0,
	"Error":     1,
	"Heartbeat": 3,
}

func (x MsgType) String() string {
	return proto.EnumName(MsgType_name, int32(x))
}

func (MsgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{1}
}

type Options struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Options) Reset()         { *m = Options{} }
func (m *Options) String() string { return proto.CompactTextString(m) }
func (*Options) ProtoMessage()    {}
func (*Options) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

func (m *Options) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Options.Unmarshal(m, b)
}
func (m *Options) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Options.Marshal(b, m, deterministic)
}
func (m *Options) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Options.Merge(m, src)
}
func (m *Options) XXX_Size() int {
	return xxx_messageInfo_Options.Size(m)
}
func (m *Options) XXX_DiscardUnknown() {
	xxx_messageInfo_Options.DiscardUnknown(m)
}

var xxx_messageInfo_Options proto.InternalMessageInfo

type Msg struct {
	Network              string   `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	MsgId                int64    `protobuf:"varint,2,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	Type                 MsgType  `protobuf:"varint,3,opt,name=type,proto3,enum=proto.MsgType" json:"type,omitempty"`
	Body                 []byte   `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Msg) Reset()         { *m = Msg{} }
func (m *Msg) String() string { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()    {}
func (*Msg) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{1}
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

func (m *Msg) GetType() MsgType {
	if m != nil {
		return m.Type
	}
	return MsgType_None
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
	return fileDescriptor_c06e4cca6c2cc899, []int{2}
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
	return fileDescriptor_c06e4cca6c2cc899, []int{3}
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
	Method               string    `protobuf:"bytes,4,opt,name=method,proto3" json:"method,omitempty"`
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
	return fileDescriptor_c06e4cca6c2cc899, []int{4}
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
	MsgId int64  `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	Laddr *Addr  `protobuf:"bytes,2,opt,name=Laddr,proto3" json:"Laddr,omitempty"`
	Raddr *Addr  `protobuf:"bytes,3,opt,name=Raddr,proto3" json:"Raddr,omitempty"`
	Body  []byte `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	// 0:走bind方式；1:走代理方式
	Type                 int32    `protobuf:"varint,5,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TCPBody) Reset()         { *m = TCPBody{} }
func (m *TCPBody) String() string { return proto.CompactTextString(m) }
func (*TCPBody) ProtoMessage()    {}
func (*TCPBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{5}
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

func (m *TCPBody) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

// 绑定请求
type Bind struct {
	MsgId                int64    `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bind) Reset()         { *m = Bind{} }
func (m *Bind) String() string { return proto.CompactTextString(m) }
func (*Bind) ProtoMessage()    {}
func (*Bind) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{6}
}

func (m *Bind) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bind.Unmarshal(m, b)
}
func (m *Bind) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bind.Marshal(b, m, deterministic)
}
func (m *Bind) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bind.Merge(m, src)
}
func (m *Bind) XXX_Size() int {
	return xxx_messageInfo_Bind.Size(m)
}
func (m *Bind) XXX_DiscardUnknown() {
	xxx_messageInfo_Bind.DiscardUnknown(m)
}

var xxx_messageInfo_Bind proto.InternalMessageInfo

func (m *Bind) GetMsgId() int64 {
	if m != nil {
		return m.MsgId
	}
	return 0
}

// 连接错误请求
type ErrorBody struct {
	PMsgId               int64    `protobuf:"varint,1,opt,name=p_msg_id,json=pMsgId,proto3" json:"p_msg_id,omitempty"`
	MsgId                int64    `protobuf:"varint,2,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	Err                  string   `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorBody) Reset()         { *m = ErrorBody{} }
func (m *ErrorBody) String() string { return proto.CompactTextString(m) }
func (*ErrorBody) ProtoMessage()    {}
func (*ErrorBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{7}
}

func (m *ErrorBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorBody.Unmarshal(m, b)
}
func (m *ErrorBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorBody.Marshal(b, m, deterministic)
}
func (m *ErrorBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorBody.Merge(m, src)
}
func (m *ErrorBody) XXX_Size() int {
	return xxx_messageInfo_ErrorBody.Size(m)
}
func (m *ErrorBody) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorBody.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorBody proto.InternalMessageInfo

func (m *ErrorBody) GetPMsgId() int64 {
	if m != nil {
		return m.PMsgId
	}
	return 0
}

func (m *ErrorBody) GetMsgId() int64 {
	if m != nil {
		return m.MsgId
	}
	return 0
}

func (m *ErrorBody) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type HeartBeat struct {
	Now                  int64    `protobuf:"varint,1,opt,name=now,proto3" json:"now,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartBeat) Reset()         { *m = HeartBeat{} }
func (m *HeartBeat) String() string { return proto.CompactTextString(m) }
func (*HeartBeat) ProtoMessage()    {}
func (*HeartBeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{8}
}

func (m *HeartBeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeat.Unmarshal(m, b)
}
func (m *HeartBeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeat.Marshal(b, m, deterministic)
}
func (m *HeartBeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeat.Merge(m, src)
}
func (m *HeartBeat) XXX_Size() int {
	return xxx_messageInfo_HeartBeat.Size(m)
}
func (m *HeartBeat) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeat.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeat proto.InternalMessageInfo

func (m *HeartBeat) GetNow() int64 {
	if m != nil {
		return m.Now
	}
	return 0
}

func init() {
	proto.RegisterEnum("proto.NetworkType", NetworkType_name, NetworkType_value)
	proto.RegisterEnum("proto.MsgType", MsgType_name, MsgType_value)
	proto.RegisterType((*Options)(nil), "proto.Options")
	proto.RegisterType((*Msg)(nil), "proto.Msg")
	proto.RegisterType((*Addr)(nil), "proto.Addr")
	proto.RegisterType((*Header)(nil), "proto.Header")
	proto.RegisterType((*HTTPBody)(nil), "proto.HTTPBody")
	proto.RegisterType((*TCPBody)(nil), "proto.TCPBody")
	proto.RegisterType((*Bind)(nil), "proto.Bind")
	proto.RegisterType((*ErrorBody)(nil), "proto.ErrorBody")
	proto.RegisterType((*HeartBeat)(nil), "proto.HeartBeat")
}

func init() {
	proto.RegisterFile("msg.proto", fileDescriptor_c06e4cca6c2cc899)
}

var fileDescriptor_c06e4cca6c2cc899 = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0xd1, 0x6a, 0xdb, 0x30,
	0x14, 0xad, 0x22, 0xcb, 0x8e, 0x6e, 0xd6, 0x60, 0xc4, 0x36, 0xfc, 0x52, 0x70, 0x05, 0x03, 0x13,
	0x58, 0x19, 0xd9, 0x17, 0x2c, 0x65, 0x90, 0xc1, 0x92, 0x15, 0xe1, 0xf7, 0xe2, 0x20, 0xe1, 0x9a,
	0x36, 0x96, 0x90, 0xd5, 0x15, 0xff, 0xc4, 0xfe, 0x63, 0x7f, 0x39, 0x24, 0xab, 0x5b, 0x07, 0xed,
	0x53, 0x9f, 0x7c, 0xae, 0xcf, 0x45, 0xe7, 0x9c, 0x7b, 0x2f, 0xd0, 0xe3, 0xd0, 0x5e, 0x18, 0xab,
	0x9d, 0x66, 0x24, 0x7c, 0x38, 0x85, 0xec, 0x87, 0x71, 0x9d, 0xee, 0x07, 0xde, 0x03, 0xde, 0x0d,
	0x2d, 0x2b, 0x20, 0xeb, 0x95, 0x7b, 0xd0, 0xf6, 0xb6, 0x40, 0x25, 0xaa, 0xa8, 0x78, 0x2c, 0xd9,
	0x3b, 0x48, 0x8f, 0x43, 0x7b, 0xdd, 0xc9, 0x62, 0x56, 0xa2, 0x0a, 0x0b, 0x72, 0x1c, 0xda, 0x6f,
	0x92, 0x71, 0x48, 0xdc, 0x68, 0x54, 0x81, 0x4b, 0x54, 0x2d, 0xd7, 0xcb, 0xe9, 0xfd, 0x8b, 0xdd,
	0xd0, 0xd6, 0xa3, 0x51, 0x22, 0x70, 0x8c, 0x41, 0x72, 0xd0, 0x72, 0x2c, 0x92, 0x12, 0x55, 0x6f,
	0x44, 0xc0, 0x7c, 0x05, 0xc9, 0x17, 0x29, 0x2d, 0x5b, 0xc2, 0xac, 0x33, 0x51, 0x6b, 0xd6, 0x19,
	0xdf, 0x6b, 0xb4, 0x75, 0x41, 0x84, 0x88, 0x80, 0xf9, 0x27, 0x48, 0xb7, 0xaa, 0x91, 0xca, 0xb2,
	0x1c, 0xf0, 0xad, 0x1a, 0x63, 0xbb, 0x87, 0xec, 0x2d, 0x90, 0x9f, 0xcd, 0xdd, 0xbd, 0x2a, 0x66,
	0x25, 0xae, 0xa8, 0x98, 0x0a, 0xfe, 0x1b, 0xc1, 0x7c, 0x5b, 0xd7, 0x57, 0x1b, 0x2d, 0xc7, 0x27,
	0xce, 0xd1, 0x53, 0xe7, 0xe7, 0x40, 0xbe, 0x37, 0x52, 0xda, 0x20, 0xb5, 0x58, 0x2f, 0xa2, 0x75,
	0xef, 0x4a, 0x4c, 0x8c, 0x97, 0xbb, 0xb7, 0x77, 0x21, 0x1b, 0x15, 0x1e, 0xb2, 0xf7, 0x90, 0x1e,
	0x95, 0xbb, 0xd1, 0x32, 0x84, 0xa1, 0x22, 0x56, 0x7f, 0x23, 0x92, 0x7f, 0x11, 0xd9, 0x07, 0x48,
	0x6f, 0x82, 0xed, 0x22, 0x2d, 0x71, 0xb5, 0x58, 0x9f, 0x46, 0x85, 0x29, 0x8b, 0x88, 0x24, 0xff,
	0x85, 0x20, 0xab, 0x2f, 0x5f, 0x6b, 0xf5, 0x1c, 0x88, 0x08, 0x2d, 0xf8, 0x99, 0x96, 0xc0, 0x3c,
	0xb7, 0x06, 0xff, 0x2f, 0xac, 0x8f, 0x4c, 0xe3, 0xf6, 0x98, 0x9f, 0x41, 0xb2, 0xe9, 0x7a, 0xf9,
	0x82, 0x19, 0xbe, 0x07, 0xfa, 0xd5, 0x5a, 0x6d, 0x83, 0xe1, 0x02, 0xe6, 0xe6, 0xfa, 0xbf, 0xae,
	0xd4, 0xec, 0x82, 0xe7, 0x17, 0xee, 0x25, 0x07, 0xac, 0xac, 0x7d, 0x1c, 0xa9, 0xb2, 0x96, 0x9f,
	0x01, 0xdd, 0xaa, 0xc6, 0xba, 0x8d, 0x6a, 0x9c, 0xa7, 0x7b, 0xfd, 0x10, 0x9f, 0xf2, 0x70, 0x55,
	0xc2, 0x62, 0x3f, 0x9d, 0xa0, 0xbf, 0x28, 0x36, 0x87, 0xc4, 0x2f, 0x36, 0x3f, 0x61, 0x19, 0xe0,
	0xfa, 0xf2, 0x2a, 0x47, 0xab, 0x8f, 0x90, 0xc5, 0x7b, 0xf3, 0xec, 0x5e, 0xf7, 0x2a, 0x3f, 0x61,
	0x14, 0x48, 0x70, 0x99, 0x23, 0x76, 0x1a, 0x05, 0x0e, 0xaa, 0x71, 0x39, 0x3e, 0xa4, 0x61, 0x32,
	0x9f, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0xb5, 0xed, 0xcf, 0x0d, 0x0f, 0x03, 0x00, 0x00,
}

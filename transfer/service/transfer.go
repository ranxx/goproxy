package service

import (
	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	"github.com/ranxx/goproxy/utils"
)

// var ...
var (
	Manage = &manage{indexMange: utils.NewIndexI64(), data: make([]Transfer, 0, 1024)}
)

func init() {
	go Manage.customer()
}

// Transfer ...
type Transfer interface {
	Start() error

	Close()

	NetWork() proto.NetworkType

	Receive(body *[]byte)

	Info() (int64, proto.Addr, proto.Addr)
}

// Manage ...
type manage struct {
	indexMange *utils.IndexI64
	data       []Transfer
}

// NewIndex ...
func (m *manage) NewIndex() int64 {
	return m.indexMange.NewIndex()
}

// Close 关闭
func (m *manage) Close() {
	for _, v := range m.data {
		v.Close()
	}
}

// Remove ...
func (m *manage) Remove(id int64) {
	if id >= int64(len(m.data)) || id < 0 {
		return
	}
	if m.data[id] == nil {
		return
	}
	m.data[id].Close()
	m.data[id] = nil
}

// Range ...
func (m *manage) Range(fc func(v Transfer)) {
	for _, v := range m.data {
		if v == nil {
			continue
		}
		fc(v)
	}
}

// RemoveByPort ...
func (m *manage) RemoveByPort(port ...int) {
	exists := map[int]bool{}
	for _, v := range port {
		exists[v] = true
	}
	for _, v := range m.data {
		id, l, _ := v.Info()
		if exists[int(l.Port)] {
			m.data[id].Close()
		}
	}
}

func (m *manage) customer() {
	for msg := range service.ReadingMsgChannel {
		trans := m.data[msg.MsgId]
		if trans == nil {
			continue
		}
		trans.Receive(&msg.Body)
	}
}

// NewAddr ...
func NewAddr(ip string, port int) *proto.Addr {
	return &proto.Addr{Ip: ip, Port: int32(port)}
}

// NewTransfer ...
func NewTransfer(localAddr, remoteAddr proto.Addr, network proto.NetworkType) Transfer {
	index := Manage.NewIndex()
	var trans Transfer
	switch network {
	case proto.NetworkType_HTTP:
		trans = newHTTP("transfer.http", network, index, localAddr, remoteAddr)
	case proto.NetworkType_TCP:
		trans = newTunnelTCP("transfer.tcp", network, index, localAddr, remoteAddr)
	default:
		trans = newHTTP("transfer.http", network, index, localAddr, remoteAddr)
	}
	Manage.data = append(Manage.data, trans)
	Manage.data[index] = trans
	return trans
}

// NewTransferWithIPPort ...
func NewTransferWithIPPort(localIP string, localPort int, remoteIP string, remotePort int, network proto.NetworkType) Transfer {
	localAddr, remoteAddr := NewAddr(localIP, localPort), NewAddr(remoteIP, remotePort)
	index := Manage.NewIndex()
	var trans Transfer
	switch network {
	case proto.NetworkType_HTTP:
		trans = newHTTP("transfer.http", network, index, *localAddr, *remoteAddr)
	case proto.NetworkType_TCP:
		trans = newTunnelTCP("transfer.tcp", network, index, *localAddr, *remoteAddr)
	default:
		trans = newHTTP("transfer.http", network, index, *localAddr, *remoteAddr)
	}
	Manage.data = append(Manage.data, trans)
	Manage.data[index] = trans
	return trans
}

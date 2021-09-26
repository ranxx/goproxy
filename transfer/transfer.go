package transfer

import (
	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	"github.com/ranxx/goproxy/utils"
)

// var ...
var (
	Manage = &manage{indexMange: utils.NewIndexI64(), data: make([]Transfer, 0, 1024*8)}
)

func init() {
	go Manage.customer()
}

// Transfer ...
type Transfer interface {
	Start()

	Receive(body []byte)
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

func (m *manage) customer() {
	for msg := range service.ReadingMsgChannel {
		trans := m.data[int(msg.MsgId)]
		if trans == nil {
			continue
		}
		trans.Receive(msg.Body)
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
		trans = newHTTP(index, localAddr, remoteAddr)
		// return &HTTP{MsgID: index, receive: make(chan []byte), localAddr: localAddr, RemoteAddr: remoteAddr}
	case proto.NetworkType_TCP:
		// return &TCP{MsgID: index, localAddr: localAddr, RemoteAddr: remoteAddr}
	default:
		trans = newHTTP(index, localAddr, remoteAddr)
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
		trans = newHTTP(index, *localAddr, *remoteAddr)
		// return &HTTP{MsgID: index, receive: make(chan []byte), localAddr: localAddr, RemoteAddr: remoteAddr}
	case proto.NetworkType_TCP:
		// return &TCP{MsgID: index, localAddr: localAddr, RemoteAddr: remoteAddr}
	default:
		trans = newHTTP(index, *localAddr, *remoteAddr)
	}
	Manage.data = append(Manage.data, trans)
	Manage.data[index] = trans
	return trans
}

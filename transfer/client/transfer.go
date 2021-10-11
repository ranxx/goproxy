package client

import (
	"log"
	"sync"
	"time"

	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	"github.com/ranxx/goproxy/utils"
)

// var ...
var (
	Manage = &manage{indexMange: utils.NewIndexI64(), data: make([][]Transfer, 1024)}
)

func init() {
	go Manage.customer()
}

// Transfer ...
type Transfer interface {
	Start()

	Close()

	Receive(body *[]byte)
}

// Manage ...
type manage struct {
	indexMange *utils.IndexI64
	data       [][]Transfer
	mutex      sync.Mutex
}

// NewIndex ...
func (m *manage) NewIndex() int64 {
	return m.indexMange.NewIndex()
}

// Close 关闭
func (m *manage) Close() {
	for _, trs := range m.data {
		for _, v := range trs {
			v.Close()
		}
	}
}

func (m *manage) customer() {
	for {
		time.Sleep(time.Second * 2)
		if service.ReadingMsgChannel == nil {
			continue
		}
		for msg := range service.ReadingMsgChannel {
			if int64(len(m.data)) <= msg.MsgId {
				m.data = append(m.data, make([]Transfer, 0, 512))
			}
			if service.CheckHTTP(msg) {
				go m.httpHandler(msg)
				continue
			}
			go m.notHTTPHandler(msg)
		}
	}
}

func (m *manage) httpHandler(msg *proto.Msg) {
	body := new(proto.HTTPBody)
	if err := body.XXX_Unmarshal(msg.Body); err != nil {
		log.Printf("transfer.%s 解码httpBody失败\n", msg.Network)
		return
	}

	if int64(len(m.data[msg.MsgId])) > body.MsgId {
		m.data[msg.MsgId][body.MsgId].Receive(&body.Body)
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	hhttp := newHTTP("transfer.http", msg, body)

	m.data[msg.MsgId] = append(m.data[msg.MsgId], hhttp)
	m.data[msg.MsgId][body.MsgId] = hhttp
	hhttp.Start()
}

func (m *manage) notHTTPHandler(msg *proto.Msg) {
	// 解包
	body := new(proto.TCPBody)
	if err := body.XXX_Unmarshal(msg.Body); err != nil {
		log.Printf("transfer.%s 解码tcpbody失败\n", msg.Network)
		return
	}

	if body.Type == 0 {
		m.tunnel(msg, body)
		return
	}

	// 处理 not tunnel
}

func (m *manage) tunnel(msg *proto.Msg, body *proto.TCPBody) {
	// if int64(len(m.data[msg.MsgId])) > body.MsgId {
	// 	m.data[msg.MsgId][body.MsgId].Receive(&body.Body)
	// 	return
	// }

	m.mutex.Lock()
	defer m.mutex.Unlock()
	ttcp := newTunnelTCP("transfer.tcp", msg, body)
	m.data[msg.MsgId] = append(m.data[msg.MsgId], ttcp)
	// m.data[msg.MsgId][body.MsgId] = ttcp
	ttcp.Start()
}

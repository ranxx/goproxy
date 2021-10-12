package service

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ranxx/goproxy/cconn"
	"github.com/ranxx/goproxy/pack"
	"github.com/ranxx/goproxy/proto"
)

// var
var (
	// 接收的消息
	ClientReadingMsgChannel = make(chan *proto.Msg, 1024)

	// 发送的消息
	ClientWritingMsgChannel = make(chan *proto.Msg, 1024)
)

// CheckHTTP ...
func CheckHTTP(msg *proto.Msg) bool {
	return proto.NetworkType_value[strings.ToUpper(msg.Network)] == int32(proto.NetworkType_HTTP)
}

// Client 客户端
// 具有重连机制，重连之后所有连接会全部断开
type Client struct {
	IP   string
	Port int
	Conn *cconn.Conn
	once *sync.Once
}

// NewClient ...
func NewClient(ip string, port int) *Client {
	return &Client{IP: ip, Port: port, once: new(sync.Once)}
}

func (c *Client) reset() {
	c.Conn = nil
	c.once = &sync.Once{}
	ClientReadingMsgChannel = make(chan *proto.Msg, 1024)
	ClientWritingMsgChannel = make(chan *proto.Msg, 1024)
}

// Close ...
func (c *Client) Close() {
	c.once.Do(func() {
		if c.Conn != nil {
			c.Conn.Close()
			c.Conn = nil
		}
		close(ClientReadingMsgChannel)
		close(ClientWritingMsgChannel)
		ClientReadingMsgChannel = nil
		ClientWritingMsgChannel = nil
	})
}

func (c *Client) dail() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.IP, c.Port), time.Second*10)
	return conn, err
}

// CheckConn 检查连接
func (c *Client) CheckConn() bool {
	if c.Conn == nil {
		return true
	}
	if c.Conn.Closed() {
		return true
	}
	return false
}

// CleanUpOnce ...
func (c *Client) CleanUpOnce() {
	c.Close()
}

// ReConnection 重连
func (c *Client) ReConnection() error {
	// 重连
	conn, err := c.dail()
	if err != nil {
		log.Println("client", err)
		return err
	}

	c.reset()

	log.Println("client", fmt.Sprintf("连上server %s -> %s", conn.LocalAddr(), conn.RemoteAddr()))

	c.Conn = cconn.NewConn(
		"client",
		conn,
		cconn.WithReadFunc(c._DefaultReadFunc()),
		cconn.WithWriteFunc(c._DefaultWriteFunc()))
	return nil
}

// ReStart ...
func (c *Client) ReStart() {
	c.StartConn()
}

// Start ...
func (c *Client) Start() {
	cconn.Checking(c, -1)
}

// StartConn ...
func (c *Client) StartConn() {
	// 开启读写
	c.Conn.Start()
	return
}

// _DefaultReadFunc ...
func (c *Client) _DefaultReadFunc() func(*cconn.Conn, <-chan struct{}) error {
	logPrefix := "client"
	return func(conn *cconn.Conn, closeC <-chan struct{}) error {
		scanner := pack.NewScanner(conn, pack.SplitFunc)
		for scanner.Scan() {
			scannedPack := new(pack.Package)
			err := scannedPack.UnpackBytes(scanner.Bytes())
			if err != nil {
				log.Println(logPrefix, "msg解包失败", err, string(scanner.Bytes()))
				return err
			}

			msg := new(proto.Msg)
			err = msg.XXX_Unmarshal(scannedPack.Msg)
			if err != nil {
				log.Println(logPrefix, "msg解码失败", err, string(scannedPack.Msg))
				return err
			}

			// 不打印心跳
			if msg.Type != proto.MsgType_Heartbeat {
				log.Println(logPrefix, "读取消息", len(msg.Body))
			}

			// 发现
			go c.msgHandler(msg)
		}
		return nil
	}
}

// _DefaultWriteFunc ...
func (c *Client) _DefaultWriteFunc() func(*cconn.Conn, <-chan struct{}) error {
	logPrefix := "client"
	return func(c *cconn.Conn, closeC <-chan struct{}) error {
		for {
			msg := new(proto.Msg)
			select {
			case msg = <-ClientWritingMsgChannel:
				if msg == nil {
					return nil
				}
			case <-closeC:
				return nil
			}

			body, err := msg.XXX_Marshal(nil, false)
			if err != nil {
				log.Println(logPrefix, "msg编码失败", err)
				return err
			}

			rbody, err := pack.NewPackage(body).PackBytes()
			if err != nil {
				log.Println(logPrefix, "mag打包失败", err)
				return err
			}

			// 不打印心跳
			if msg.Type != proto.MsgType_Heartbeat {
				log.Println(logPrefix, "写回消息", len(msg.Body))
			}

			wn, err := c.Write(rbody)
			if err != nil {
				log.Println(logPrefix, "写入消息失败", "len(msg):", len(rbody), err)
				return err
			}

			if wn != len(rbody) {
				log.Println(logPrefix, "写入消息未成功", err)
				return err
			}
		}
	}
}

func (c *Client) msgHandler(msg *proto.Msg) {
	switch msg.Type {
	case proto.MsgType_Heartbeat:
		// 发送心跳
		go c.heartbeat(msg)
	default:
		ClientReadingMsgChannel <- msg
	}
}

func (c *Client) heartbeat(msg *proto.Msg) {
	hb := new(proto.HeartBeat)
	hb.XXX_Unmarshal(msg.Body)
	log.Println("client", "收到心跳", hb.Now)
	hb = &proto.HeartBeat{Now: time.Now().Unix()}
	body, _ := hb.XXX_Marshal(nil, false)
	ClientWritingMsgChannel <- &proto.Msg{
		Body: body,
		Type: proto.MsgType_Heartbeat,
	}
}

package service

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ranxx/goproxy/cconn"
	"github.com/ranxx/goproxy/proto"
)

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
	ReadingMsgChannel = make(chan *proto.Msg, 1024)
	WritingMsgChannel = make(chan *proto.Msg, 1024)
}

// Close ...
func (c *Client) Close() {
	c.once.Do(func() {
		if c.Conn != nil {
			c.Conn.Close()
			c.Conn = nil
		}
		close(ReadingMsgChannel)
		close(WritingMsgChannel)
		ReadingMsgChannel = nil
		WritingMsgChannel = nil
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
		cconn.WithReadFunc(_DefaultReadFunc("client")),
		cconn.WithWriteFunc(_DefaultWriteFunc("client")))
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

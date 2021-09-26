package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/ranxx/goproxy/pack"
	"github.com/ranxx/goproxy/proto"
)

// Client ...
type Client struct {
	IP   string
	Port int
}

// NewClient ...
func NewClient(ip string, port int) *Client {
	return &Client{IP: ip, Port: port}
}

// Start ...
func (c *Client) Start() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.IP, c.Port))
	if err != nil {
		panic(err)
	}
	log.Printf("client 已连上server %s:%d\n", c.IP, c.Port)
	go c.Customer()
	c.StartConn(conn)
}

// StartConn ...
func (c *Client) StartConn(conn net.Conn) {
	// ReadFunc ...
	_ReadFunc := func(c *Conn) error {
		scanner := pack.NewScanner(c.Conn)
		for scanner.Scan() {
			scannedPack := new(pack.Package)
			err := scannedPack.UnpackBytes(scanner.Bytes())
			if err != nil {
				log.Println("解包pack失败", err, scanner.Bytes())
				return err
			}

			msg := new(proto.Msg)
			err = msg.XXX_Unmarshal(scannedPack.Msg)
			if err != nil {
				log.Println("解包msg失败", err, scanner.Bytes())
				return err
			}

			log.Println("client", "开始读取", len(msg.Body), string(msg.Body))
			ReadingMsgChannel <- msg
		}
		return nil
	}

	// WriteFunc ...
	_WriteFunc := func(c *Conn) error {
		for msg := range WritingMsgChannel {
			body, err := msg.XXX_Marshal(nil, false)
			if err != nil {
				log.Println("msg打包失败", err)
				return err
			}
			body, err = pack.NewPackage(body).PackBytes()
			if err != nil {
				log.Println("pack打包失败", err)
				return err
			}

			log.Println("client", "开始回写", len(msg.Body), string(msg.Body))
			c.Write(body)
		}
		return nil
	}

	// 开启读写
	cconn := NewConn(conn)
	cconn.WriteFunc = _WriteFunc
	cconn.ReadFunc = _ReadFunc
	cconn.Start()
}

// Customer ...
func (c *Client) Customer() {
	for v := range ReadingMsgChannel {
		// log.Println("client", "读取到service的消息", string(v.Body))
		// http,转发请求
		if CheckHTTP(v) {
			go c.HTTPHandler(v)
			continue
		}
		// tcp，主动连接
		go c.TCPHandler(v)
	}
}

// HTTPHandler ...
func (c *Client) HTTPHandler(msg *proto.Msg) {
	// 先解析
	body := new(proto.HTTPBody)
	if err := body.XXX_Unmarshal(msg.Body); err != nil {
		log.Println("client", "处理http 解析body失败", err)
		panic(err)
	}

	// log.Println("client", "读取到service的消息", string(body.Body))

	// 转发
	request, err := http.NewRequest(body.Method, fmt.Sprintf("%s:%d%s", "http://localhost", body.Laddr.Port, body.Url), bytes.NewReader(body.Body))
	if err != nil {
		log.Println("client", "处理http 转发http失败", err)
		panic(err)
	}
	defer request.Body.Close()

	for _, header := range body.Header {
		for _, v := range header.Value {
			request.Header.Add(header.Key, v)
		}
	}

	response, err := new(http.Client).Do(request)
	if err != nil {
		log.Println("client", "处理http 请求http失败", err)
		panic(err)
	}
	defer response.Body.Close()

	bbody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("client", "处理http 读取http返回值失败", err)
		panic(err)
	}

	body.Body = bbody
	body.Header = make([]*proto.Header, 0, len(response.Proto))
	for key, values := range response.Header {
		body.Header = append(body.Header, &proto.Header{Key: key, Value: values})
	}

	wbody, err := body.XXX_Marshal(nil, false)
	if err != nil {
		log.Println("client", "处理http 编码返回body失败", err)
		panic(err)
	}

	WritingMsgChannel <- &proto.Msg{
		Network: proto.NetworkType_HTTP.String(),
		MsgId:   msg.MsgId,
		Body:    wbody,
	}

	return
}

// TCPHandler ...
// 主动连接
func (c *Client) TCPHandler(msg *proto.Msg) {
	// 新
}

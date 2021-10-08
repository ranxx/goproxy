package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/ranxx/goproxy/pack"
	"github.com/ranxx/goproxy/proto"
)

// TunnelTCP ...
type TunnelTCP struct {
	logPrefix             string
	msg                   *proto.Msg
	body                  *proto.TCPBody
	once                  *sync.Once
	localConn, remoteConn net.Conn
}

// newTunnelTCP 新建tcp
func newTunnelTCP(logPrefix string, msg *proto.Msg, body *proto.TCPBody) *TunnelTCP {
	return &TunnelTCP{
		logPrefix: logPrefix,
		msg:       msg,
		body:      body,
		once:      new(sync.Once),
	}
}

// Start 开始
func (t *TunnelTCP) Start() {
	// 连接本地
	localConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.body.Laddr.Ip, t.body.Laddr.Port))
	if err != nil {
		// 连接失败直接报错
		panic(err)
	}
	t.localConn = localConn

	// 连接远程
	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.body.Raddr.Ip, t.body.Raddr.Port))
	if err != nil {
		// 连接失败直接报错
		panic(err)
	}
	t.remoteConn = remoteConn

	log.Println(t.logPrefix, fmt.Sprintf("%s -> %s", remoteConn.RemoteAddr(), localConn.RemoteAddr()))

	// 开始发送 bind 请求
	bind := proto.Bind{MsgId: t.body.MsgId}

	bindBody, err := bind.XXX_Marshal(nil, false)
	if err != nil {
		panic(err)
	}

	bindBodyBytes, err := pack.NewPackage(bindBody).PackBytes()
	if err != nil {
		panic(err)
	}

	_, err = remoteConn.Write(bindBodyBytes)
	if err != nil {
		panic(err)
	}

	// 开启读写
	go func() {
		defer t.Close()
		io.Copy(localConn, remoteConn)
	}()

	go func() {
		defer t.Close()
		io.Copy(remoteConn, localConn)
	}()
}

// Close 关闭
func (t *TunnelTCP) Close() {
	t.once.Do(func() {
		t.localConn.Close()
		t.remoteConn.Close()
	})
}

// Receive 接受数据
func (t *TunnelTCP) Receive(msg *[]byte) {
	return
}

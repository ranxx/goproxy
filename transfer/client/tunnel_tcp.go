package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/ranxx/goproxy/pack"
	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
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

func (t *TunnelTCP) newErrorProto(e string) *proto.ErrorBody {
	return &proto.ErrorBody{
		PMsgId: t.msg.MsgId,
		MsgId:  t.body.MsgId,
		Err:    e,
	}
}

func (t *TunnelTCP) sendError(format string, a ...interface{}) {
	ebody, _ := t.newErrorProto(fmt.Sprintf(format, a...)).XXX_Marshal(nil, false)
	service.ClientWritingMsgChannel <- &proto.Msg{
		MsgId: t.msg.MsgId,
		Type:  proto.MsgType_Error,
		Body:  ebody,
	}
}

// Start 开始
func (t *TunnelTCP) Start() {
	// 连接本地
	localConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.body.Laddr.Ip, t.body.Laddr.Port))
	if err != nil {
		// 连接失败直接报错
		log.Println(t.logPrefix, "连接local失败", err)
		t.sendError("连接local失败 %s", err.Error())
		return
	}
	t.localConn = localConn

	// 连接远程
	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.body.Raddr.Ip, t.body.Raddr.Port))
	if err != nil {
		// 连接失败直接报错
		log.Println(t.logPrefix, "连接remote失败", err)
		t.sendError("连接remote失败 %s", err.Error())
		return
	}
	t.remoteConn = remoteConn

	log.Println(t.logPrefix, fmt.Sprintf("%s -> %s", remoteConn.RemoteAddr(), localConn.RemoteAddr()))

	// 开始发送 bind 请求
	bind := proto.Bind{MsgId: t.body.MsgId}

	bindBody, _ := bind.XXX_Marshal(nil, false)
	bindBodyBytes, _ := pack.NewPackage(bindBody).PackBytes()

	_, err = remoteConn.Write(bindBodyBytes)
	if err != nil {
		log.Println(t.logPrefix, "写入remote失败", err)
		t.sendError("写入remote失败 %s", err.Error())
		return
	}

	// 开启读写
	go func() {
		defer t.Close()
		if rn, err := io.Copy(localConn, remoteConn); err != nil {
			log.Println(t.logPrefix, rn, err)
		}
	}()

	go func() {
		defer t.Close()
		if rn, err := io.Copy(remoteConn, localConn); err != nil {
			log.Println(t.logPrefix, rn, err)
		}
	}()
}

// Close 关闭
func (t *TunnelTCP) Close() {
	t.once.Do(func() {
		if t.localConn != nil {
			t.localConn.Close()
		}
		if t.remoteConn != nil {
			t.remoteConn.Close()
		}
	})
}

// Receive 接受数据
func (t *TunnelTCP) Receive(msg *[]byte) {
	return
}

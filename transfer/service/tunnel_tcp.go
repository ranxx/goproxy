package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ranxx/goproxy/pack"
	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	"github.com/ranxx/goproxy/utils"
)

type rw struct {
	net.Conn
	r io.Reader
	w io.Writer
}

// TunnelTCP ...
type TunnelTCP struct {
	msgID        int64
	indexManage  *utils.IndexI64
	indexs       []*rw
	laddr, raddr proto.Addr
	once         *sync.Once
}

// newTunnelTCP ...
func newTunnelTCP(msgID int64, localAddr, remoteAddr proto.Addr) Transfer {
	return &TunnelTCP{
		msgID:       msgID,
		indexManage: utils.NewIndexI64(),
		laddr:       localAddr,
		raddr:       remoteAddr,
		indexs:      make([]*rw, 512),
		once:        new(sync.Once),
	}
}

// Receive ...
func (t *TunnelTCP) Receive(body []byte) {
	// TODO: 处理报错的问题

	tcpBody := new(proto.Bind)
	if err := tcpBody.XXX_Unmarshal(body); err != nil {
		log.Println("transfer.tcp", "解码返回的消息失败", err)
		panic(err)
	}
	fmt.Println("收到消息----", tcpBody)
	return
}

// Start 开启服务
func (t *TunnelTCP) Start() {
	log.Println("transfer.tcp", fmt.Sprintf("%s:%d -> %s:%d runing", t.laddr.Ip, t.laddr.Port, t.raddr.Ip, t.raddr.Port))
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.laddr.Ip, t.laddr.Port))
	if err != nil {
		log.Println("transfer.tcp", "监听 %s:%d 失败", t.laddr.Ip, t.laddr.Port, err)
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		index := t.indexManage.NewIndex()
		rw := &rw{Conn: conn, r: conn, w: conn}

		t.indexs = append(t.indexs, rw)
		t.indexs[index] = rw

		go t.connection(index, conn)
	}
}

func (t *TunnelTCP) customerConnection(index int64, conn net.Conn, reader *bufio.Reader) {
	pTCPBody := &proto.TCPBody{
		MsgId: index,
		Laddr: &t.raddr,
		Raddr: &proto.Addr{Ip: "", Port: t.laddr.Port},
		Body:  nil,
	}
	tcpBody, err := pTCPBody.XXX_Marshal(nil, false)
	if err != nil {
		log.Println("transfer.tcp", "编码body失败", err)
		panic(err)
	}
	t.indexs[index].r = reader
	service.WritingMsgChannel <- &proto.Msg{
		Network: proto.NetworkType_TCP.String(),
		MsgId:   t.msgID,
		Body:    tcpBody,
	}
	return
}

func (t *TunnelTCP) connection(index int64, conn net.Conn) {
	// 首选设置超时时间，因为这个时候并不知道conn是client请求还是用户请求，客户端请求必定有数据
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))

	// 因为数据都是通过 pack打包，所以最先读取 pack 的前 n个字节
	prelen := int(pack.Empty().PreLength())

	// 这里为什么需要用 bufio.NewReader 是是因为有Peek这个方法，不会影响 后续的 read
	reader := bufio.NewReader(conn)
	pre, err := reader.Peek(prelen)

	// 重置超时时间，防止后续读写问题
	conn.SetDeadline(time.Time{})

	// 如果这里报错，或者读取的长度不够，可以判断不是client的请求，直接走 customerConnection
	if err != nil || len(pre) != prelen {
		t.customerConnection(index, conn, reader)
		return
	}

	// 使用pack解包前 n个字节，如果失败，直接走 customerConnection
	packer := new(pack.Package)
	if err := packer.ReadPre(bytes.NewReader(pre)); err != nil {
		t.customerConnection(index, conn, reader)
		return
	}

	// 判断该请求是否为client请求，否则认为该请求为 customerConnection
	if !packer.IsPackage(time.Second * 3) {
		t.customerConnection(index, conn, reader)
		return
	}

	// 开始读取pack后续字节 ，需要将 prelen 字节先读取
	reader.Read(make([]byte, prelen))
	packer.ReadLast(reader)

	// 开始解析 bind 数据，拿到 customer 对应的 msgid
	bind := new(proto.Bind)
	if err := bind.XXX_Unmarshal(packer.Msg); err != nil {
		// 这里解析失败，证明数据有问题，需要 panic
		log.Println("msg:", string(packer.Msg))
		panic(err)
	}

	// 拿到 customerConn
	if len(t.indexs) <= int(bind.MsgId) {
		// 这里属于无效conn
		panic("")
	}

	inRW := t.indexs[bind.MsgId]

	// 开启读写
	once := new(sync.Once)
	go func() {
		log.Println(io.Copy(conn, inRW.r))
		once.Do(func() {
			inRW.Conn.Close()
			conn.Close()
		})
	}()
	go func() {
		log.Println(io.Copy(inRW.w, reader))
		once.Do(func() {
			inRW.Conn.Close()
			conn.Close()
		})
	}()
}

// Close 关闭
func (t *TunnelTCP) Close() {
	t.once.Do(func() {
		for _, v := range t.indexs {
			if v == nil || v.Conn == nil {
				continue
			}
			v.Conn.Close()
		}
	})
}

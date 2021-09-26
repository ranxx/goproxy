package transfer

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	"github.com/ranxx/goproxy/utils"
)

// TCP ...
type TCP struct {
	MsgID                 int64
	IndexManage           *utils.IndexI64
	Indexs                []net.Conn
	LocalAddr, RemoteAddr proto.Addr
}

type tcpConn struct {
	net.Conn
	proto.TCPBody
}

// newTCP ...
func newTCP(msgID int64, localAddr, remoteAddr proto.Addr) Transfer {
	return &TCP{
		MsgID:       msgID,
		IndexManage: utils.NewIndexI64(),
		Indexs:      make([]net.Conn, 512),
		LocalAddr:   localAddr,
		RemoteAddr:  remoteAddr,
	}
}

// Receive ...
func (t *TCP) Receive(body []byte) {
	tcpBody := new(proto.TCPBody)
	if err := tcpBody.XXX_Unmarshal(body); err != nil {
		log.Println("transfer.tcp", "解码返回的消息失败", err)
		panic(err)
	}

	log.Println("transfer.tcp", fmt.Sprintf("msgID:%d 收到 %d 返回的数据", t.MsgID, tcpBody.MsgId))
	t.Indexs[tcpBody.MsgId].Write(tcpBody.Body)
}

// Start ...
func (t *TCP) Start() {
	log.Println("transfer.tcp", fmt.Sprintf("%s:%d -> %s:%d runing", t.LocalAddr.Ip, t.LocalAddr.Port, t.RemoteAddr.Ip, t.RemoteAddr.Port))
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.LocalAddr.Ip, t.LocalAddr.Port))
	if err != nil {
		log.Println("transfer.tcp", "监听 %s:%d 失败", t.LocalAddr.Ip, t.LocalAddr.Port, err)
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		index := t.IndexManage.NewIndex()

		t.Indexs = append(t.Indexs, conn)
		t.Indexs[index] = conn

		go t.connection(index, conn)
	}
}

func (t *TCP) connection(index int64, conn net.Conn) {
	// 开启读数据
	// 开启绑定
	//
	scanner := bufio.NewScanner(conn)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		log.Println("transfer.tcp", atEOF, len(data))
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if !atEOF {
			return len(data), data[:], nil
		}
		return 0, nil, nil
	})
	for scanner.Scan() {
		// 开启读
		rbody := scanner.Bytes()
		pTCPBody := &proto.TCPBody{
			MsgId: index,
			Laddr: &t.RemoteAddr,
			Raddr: &t.LocalAddr,
			Body:  rbody,
		}
		tcpBody, err := pTCPBody.XXX_Marshal(nil, false)
		if err != nil {
			log.Println("transfer.tcp", "编码body失败", err)
			panic(err)
		}
		service.WritingMsgChannel <- &proto.Msg{
			Network: proto.NetworkType_TCP.String(),
			MsgId:   t.MsgID,
			Body:    tcpBody,
		}
	}
}

package service

import (
	"log"
	"net"
	"time"

	"github.com/ranxx/goproxy/cconn"
	"github.com/ranxx/goproxy/pack"
	"github.com/ranxx/goproxy/proto"
)

type client struct {
	Conn              net.Conn
	cconn             *cconn.Conn
	ReadingMsgChannel chan *proto.Msg
	WritingMsgChannel chan *proto.Msg
	HeartBeatChannel  chan *proto.HeartBeat
	remove            func()
}

func newClient(conn net.Conn, remove func()) *client {
	cli := &client{
		Conn:              conn,
		ReadingMsgChannel: GlobalReadingMsgChannel,
		WritingMsgChannel: make(chan *proto.Msg, 512),
		HeartBeatChannel:  make(chan *proto.HeartBeat, 2),
		remove:            remove,
	}
	cli.cconn = cconn.NewConn(
		"service",
		cli.Conn,
		cconn.WithReadFunc(cli._DefaultReadFunc()),
		cconn.WithWriteFunc(cli._DefaultWriteFunc()))
	return cli
}

func (cli *client) Start() {
	go cli.cconn.Start()
	// 开启心跳
	go cli.HeartBeat()
}

func (cli *client) msgHandler(msg *proto.Msg) {
	switch msg.Type {
	case proto.MsgType_Heartbeat:
		go func() {
			hb := new(proto.HeartBeat)
			hb.XXX_Unmarshal(msg.Body)
			cli.HeartBeatChannel <- hb
		}()
	default:
		cli.ReadingMsgChannel <- msg
	}
}

// _DefaultReadFunc ...
func (cli *client) _DefaultReadFunc() func(*cconn.Conn, <-chan struct{}) error {
	return func(c *cconn.Conn, closeC <-chan struct{}) error {
		scanner := pack.NewScanner(c, pack.SplitFunc)
		for scanner.Scan() {
			scannedPack := new(pack.Package)
			err := scannedPack.UnpackBytes(scanner.Bytes())
			if err != nil {
				log.Println("service", "msg解包失败", err, string(scanner.Bytes()))
				return err
			}

			msg := new(proto.Msg)
			err = msg.XXX_Unmarshal(scannedPack.Msg)
			if err != nil {
				log.Println("service", "msg解码失败", err, string(scannedPack.Msg))
				return err
			}

			// 心跳消息不打印
			if msg.Type != proto.MsgType_Heartbeat {
				log.Println("service", "读取消息", len(msg.Body))
			}

			go cli.msgHandler(msg)
		}
		return nil
	}
}

// _DefaultWriteFunc ...
func (cli *client) _DefaultWriteFunc() func(*cconn.Conn, <-chan struct{}) error {
	return func(c *cconn.Conn, closeC <-chan struct{}) error {
		for {
			msg := new(proto.Msg)
			select {
			case msg = <-cli.WritingMsgChannel:
				if msg == nil {
					return nil
				}
			case <-closeC:
				return nil
			}

			body, err := msg.XXX_Marshal(nil, false)
			if err != nil {
				log.Println("service", "msg编码失败", err)
				return err
			}

			rbody, err := pack.NewPackage(body).PackBytes()
			if err != nil {
				log.Println("service", "mag打包失败", err)
				return err
			}

			// 心跳消息不打印
			if msg.Type != proto.MsgType_Heartbeat {
				log.Println("service", "写回消息", len(msg.Body))
			}

			wn, err := c.Write(rbody)
			if err != nil {
				log.Println("service", "写入消息失败", "len(msg):", len(rbody), err)
				return err
			}

			if wn != len(rbody) {
				log.Println("service", "写入消息未成功", err)
				return err
			}
		}
	}
}

func (cli *client) Closed() bool {
	return cli.cconn.Closed()
}

// HeartBeat 心跳检测，如果超时或者未响应，断开这个client，并移除
func (cli *client) HeartBeat() {
	for {
		hb := proto.HeartBeat{Now: time.Now().Unix()}
		body, _ := hb.XXX_Marshal(nil, false)

		cli.WritingMsgChannel <- &proto.Msg{
			Body: body,
			Type: proto.MsgType_Heartbeat,
		}

		// 等待返回值，超时则断开请求，并清理数据
		rhb := &proto.HeartBeat{}
		select {
		case <-time.After(time.Second * 3):
		case rhb = <-cli.HeartBeatChannel:
		}
		// 检测时间 是否在5秒钟内，这里不保证两端时间相同，所以取±3
		if rhb.Now-hb.Now >= -3 && rhb.Now-hb.Now <= 3 {
			// 通过
			time.Sleep(time.Second * 3)
		} else {
			log.Println("service", cli.Conn.RemoteAddr(), "心跳检测失败")
			// 未通过
			cli.remove()
			// 断开连接
			cli.Conn.Close()
			return
		}
	}
}

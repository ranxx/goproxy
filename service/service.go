package service

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ranxx/goproxy/proto"
)

// var
var (
	// 接收的消息
	GlobalReadingMsgChannel = make(chan *proto.Msg, 1024)

	// 发送的消息
	GlobalWritingMsgChannel = make(chan *proto.Msg, 1024)
)

// SendMsg ...
func SendMsg(machine string, msg *proto.Msg) bool {
	// 如果为空，随机
	if machine == "" {
		GlobalWritingMsgChannel <- msg
		return true
	}
	v, ok := Svc.Clients.Load(machine)
	if !ok {
		return false
	}
	client, ok := v.(*client)
	if !ok {
		return false
	}
	client.WritingMsgChannel <- msg
	return true
}

// Svc ...
var Svc *Service

// Service 隧道服务
type Service struct {
	IP       string
	Port     int
	listener net.Listener
	Clients  sync.Map // 每个client都应该有自己的读取和写入channel
}

// NewService ...
func NewService(ip string, port int) *Service {
	Svc = &Service{IP: ip, Port: port}
	return Svc
}

// Start ...
func (s *Service) Start() {
	// 开启监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Println("开启service失败")
		panic(err)
	}
	log.Println("start service", fmt.Sprintf("%s:%d", s.IP, s.Port))
	s.listener = listener

	// 开启全局消息检测
	go s.checkGlobalWritingMsgChannel()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("service监听失败")
			panic(err)
		}
		log.Println("service 新client", conn.RemoteAddr())
		go s.StartConn(conn)
	}
}

// StartConn ...
func (s *Service) StartConn(conn net.Conn) {
	client := newClient(conn, func() {
		s.Clients.Delete(conn.RemoteAddr().String())
	})
	s.Clients.LoadOrStore(conn.RemoteAddr().String(), client)
	client.Start()
}

// Close 关闭服务
func (s *Service) Close() {
	if s.listener != nil {
		// 关闭服务
		s.listener.Close()
	}
	// 关闭所有连接
	s.Clients.Range(func(key, value interface{}) bool {
		if value == nil {
			return true
		}
		cli, ok := value.(*client)
		if !ok {
			return true
		}
		cli.Conn.Close()
		return true
	})
}

// 每隔几秒检测一遍全局Channel
func (s *Service) checkGlobalWritingMsgChannel() {
	go func() {
		for v := range GlobalWritingMsgChannel {
			// 拿到一个client
			cli := s.checkClient()
			cli.WritingMsgChannel <- v
		}
	}()
}

func (s *Service) checkClient() *client {
	var cli *client
	for cli == nil {
		s.Clients.Range(func(key, value interface{}) bool {
			if value == nil {
				return true
			}
			_cli, ok := value.(*client)
			if !ok {
				return true
			}
			cli = _cli
			return false
		})
		if cli != nil {
			return cli
		}
		// 如果没查到，需要停留3秒
		time.Sleep(time.Second * 3)
	}
	return cli
}

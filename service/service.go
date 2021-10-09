package service

import (
	"fmt"
	"log"
	"net"

	"github.com/ranxx/goproxy/cconn"
)

// DefaultPort ...
var DefaultPort int = 12341

// Service 隧道服务
type Service struct {
	IP       string
	Port     int
	listener net.Listener
}

// NewService ...
func NewService(ip string, port int) *Service {
	return &Service{IP: ip, Port: port}
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("service监听失败")
			panic(err)
		}
		log.Println("service 新连接", conn.RemoteAddr())
		go s.StartConn(conn)
	}
}

// StartConn ...
func (s *Service) StartConn(conn net.Conn) {
	// 读取数据
	cconn := cconn.NewConn(
		"service",
		conn,
		cconn.WithReadFunc(_DefaultReadFunc("service")),
		cconn.WithWriteFunc(_DefaultWriteFunc("service")))
	cconn.Start()
}

// Close 关闭服务
func (s *Service) Close() {
	// 关闭服务
	s.listener.Close()
	// 关闭所有连接
	close(ReadingMsgChannel)
	close(WritingMsgChannel)
}

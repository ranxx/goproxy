package service

import (
	"fmt"
	"log"
	"net"
)

// Service 隧道服务
type Service struct {
	IP   string
	Port int
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("service监听失败")
			panic(err)
		}
		log.Println("service新连接", conn.RemoteAddr())
		go s.StartConn(conn)
	}
}

// StartConn ...
func (s *Service) StartConn(conn net.Conn) {
	// 读取数据
	cconn := NewConn("service", conn)
	cconn.WithReadFunc(_DefaultReadFunc("service", cconn)).
		WithWriteFunc(_DefaultWriteFunc("service", cconn)).Start()
}

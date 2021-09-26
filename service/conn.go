package service

import (
	"net"
	"time"
)

// Conn ...
type Conn struct {
	net.Conn

	ReadFunc func(*Conn) error

	WriteFunc func(*Conn) error
}

// NewConn ...
func NewConn(conn net.Conn) *Conn {
	return &Conn{Conn: conn}
}

// Start ...
func (conn *Conn) Start() {
	// 怎么规定 数据读写
	// 数据读
	go func() {
		defer func() {
			conn.Conn.Close()
		}()
		for {
			if conn.ReadFunc == nil {
				time.Sleep(time.Second)
				continue
			}
			conn.ReadFunc(conn)
			break
		}
		conn.Close()
	}()
	// 数据写
	go func() {
		defer func() {
			conn.Conn.Close()
		}()
		for {
			if conn.WriteFunc == nil {
				time.Sleep(time.Second)
				continue
			}
			conn.WriteFunc(conn)
		}
	}()
}

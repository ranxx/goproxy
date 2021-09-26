package service

import (
	"log"
	"net"
	"sync"
	"time"
)

// Conn ...
type Conn struct {
	net.Conn

	logPrefix string

	closeC chan struct{}

	once sync.Once

	readFunc func(*Conn, chan struct{}) error

	writeFunc func(*Conn, chan struct{}) error
}

// NewConn ...
func NewConn(logPrefix string, conn net.Conn) *Conn {
	return &Conn{Conn: conn, logPrefix: logPrefix, closeC: make(chan struct{})}
}

// Start ...
func (conn *Conn) Start() *Conn {
	// 怎么规定 数据读写
	// 数据读

	// 开启关闭
	go conn.closing()

	// 开启读
	go conn.reading()

	// 开启写
	go conn.writing()

	return conn
}

// WithReadFunc ...
func (conn *Conn) WithReadFunc(f func(*Conn, chan struct{}) error) *Conn {
	conn.readFunc = f
	return conn
}

// WithWriteFunc ...
func (conn *Conn) WithWriteFunc(f func(*Conn, chan struct{}) error) *Conn {
	conn.writeFunc = f
	return conn
}

func (conn *Conn) close() {
	conn.once.Do(
		func() {
			close(conn.closeC)
		},
	)
}

func (conn *Conn) closing() {
	<-conn.closeC
	conn.Conn.Close()
	log.Println(conn.logPrefix, "断开连接", conn.RemoteAddr())
}

// reading 读
func (conn *Conn) reading() {
	defer func() {
		conn.close()
	}()
	for {
		if conn.readFunc == nil {
			time.Sleep(time.Second)
			continue
		}
		conn.readFunc(conn, conn.closeC)
		break
	}
}

// writing ...
func (conn *Conn) writing() {
	defer func() {
		conn.close()
	}()
	for {
		if conn.writeFunc == nil {
			time.Sleep(time.Second)
			continue
		}
		conn.writeFunc(conn, conn.closeC)
	}
}

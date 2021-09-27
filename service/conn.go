package service

import (
	"fmt"
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

	closeB bool

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

func (conn *Conn) closed() bool {
	return conn.closeB
}

func (conn *Conn) closing() {
	<-conn.closeC
	conn.closeB = true
	conn.Conn.Close()
	log.Println(conn.logPrefix, fmt.Sprintf("断开连接 %s -> %s", conn.LocalAddr(), conn.RemoteAddr()))
	return
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
		break
	}
}

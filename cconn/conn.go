package cconn

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

	once *sync.Once

	readFunc func(*Conn, <-chan struct{}) error

	writeFunc func(*Conn, <-chan struct{}) error
}

// NewConn ...
func NewConn(logPrefix string, conn net.Conn, opts ...Options) *Conn {
	cconn := &Conn{Conn: conn, logPrefix: logPrefix, closeC: make(chan struct{}), once: new(sync.Once)}
	for _, v := range opts {
		v(cconn)
	}
	return cconn
}

// Start ...
func (conn *Conn) Start() *Conn {
	// 开启关闭监听
	go conn.closing()

	// 开启读
	go conn.reading()

	// 开启写
	go conn.writing()

	return conn
}

// Close 关闭连接
func (conn *Conn) Close() {
	conn.close()
}

// Closed 是否关闭
func (conn *Conn) Closed() bool {
	return conn.closed()
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
		if err := conn.readFunc(conn, conn.closeC); err != nil {
			log.Println(conn.logPrefix, "read:", err)
		}
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
		if err := conn.writeFunc(conn, conn.closeC); err != nil {
			log.Println(conn.logPrefix, "write:", err)
		}
		break
	}
}

package cconn

import "sync"

// Options 参数
type Options func(*Conn)

// WithCloseChannel 可配置关闭chanel
func WithCloseChannel(cn chan struct{}) Options {
	return func(c *Conn) {
		c.closeC = cn
	}
}

// WithSyncOnce 可配置once struct
func WithSyncOnce(once *sync.Once) Options {
	return func(c *Conn) {
		c.once = once
	}
}

// WithLogPrefix 可配置日志前缀
func WithLogPrefix(logPrefix string) Options {
	return func(c *Conn) {
		c.logPrefix = logPrefix
	}
}

// WithReadFunc 可配置读函数
func WithReadFunc(fn func(*Conn, <-chan struct{}) error) Options {
	return func(c *Conn) {
		c.readFunc = fn
	}
}

// WithWriteFunc 可配置写函数
func WithWriteFunc(fn func(*Conn, <-chan struct{}) error) Options {
	return func(c *Conn) {
		c.writeFunc = fn
	}
}

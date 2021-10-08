package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// const
const (
	reset = "reset"
)

// Reset 重置
func Reset() string {
	return reset
}

// SignalerWrap 包装器
type SignalerWrap func() string

// String SignalerWrap 包装器
func (s SignalerWrap) String() string { return s() }

// Signal SignalerWrap 包装器
func (s SignalerWrap) Signal() {}

// Options SignalerWrap 包装器
func (s SignalerWrap) Options() string { return s() }

// IgnoreSignal 忽略信号
func IgnoreSignal(fc func(), sig ...os.Signal) {
	defaultS := []os.Signal{syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL}
	sigMap := map[string]bool{}

	resetB := false
	for _, v := range sig {
		s, ok := v.(SignalerWrap)
		if !ok {
			sigMap[v.String()] = true
			continue
		}
		if s.Options() == reset {
			resetB = true
		}
	}

	if !resetB {
		for _, v := range defaultS {
			sigMap[v.String()] = true
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	for {
		si := <-ch
		switch si {
		default:
			if sigMap[si.String()] {
				fc()
				return
			}
		}
	}
}

package cconn

import "time"

// Checking 检查
func Checking(rc Reconn, recount int) (err error) {
	if recount < 0 {
		for {
			forloop(rc)
		}
	}
	for i := 0; i < recount; i++ {
		err = forloop(rc)
	}
	return
}

func forloop(rc Reconn) error {
	for {
		if !rc.CheckConn() {
			time.Sleep(time.Second * 2)
			continue
		}

		rc.CleanUpOnce()

		if err := rc.ReConnection(); err != nil {
			time.Sleep(time.Second * 3)
			return err
		}

		rc.ReStart()
	}
}

// Reconn 重连机制: 检查连接，清理数据，重连，重新开启
type Reconn interface {
	// 重连
	CheckConn() bool

	// 清理数据
	CleanUpOnce()

	// 重连
	ReConnection() error

	// 开启读写
	ReStart()
}

package service

// CManage ...
type CManage struct {
	data []ICustomer
}

// ICustomer ...
type ICustomer interface {
	// Start 开始处理
	Start()

	// 关闭
	Close()

	// Receive 接受消息
	Receive(*[]byte)
}

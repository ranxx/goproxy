package main

import (
	"fmt"

	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	"github.com/ranxx/goproxy/transfer"
)

func main() {
	fmt.Println("Hello World")

	// 开启转发
	go transfer.NewTransferWithIPPort("", 3334, "", 3333, proto.NetworkType_HTTP).Start()

	go transfer.NewTransferWithIPPort("", 3335, "", 3333, proto.NetworkType_HTTP).Start()

	go transfer.NewTransferWithIPPort("", 3336, "", 3333, proto.NetworkType_HTTP).Start()

	srv := service.NewService("", 12341)
	srv.Start()
}

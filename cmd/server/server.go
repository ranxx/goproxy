package main

import (
	"fmt"
	"log"

	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	transfer "github.com/ranxx/goproxy/transfer/service"
	"github.com/ranxx/goproxy/utils"
)

func main() {
	fmt.Println("Hello World")

	// 开启转发
	go transfer.NewTransferWithIPPort("", 3334, "", 3333, proto.NetworkType_HTTP).Start()

	go transfer.NewTransferWithIPPort("", 3335, "", 3333, proto.NetworkType_HTTP).Start()

	go transfer.NewTransferWithIPPort("", 3336, "", 3333, proto.NetworkType_HTTP).Start()

	go transfer.NewTransferWithIPPort("", 2022, "", 22, proto.NetworkType_TCP).Start()

	go transfer.NewTransferWithIPPort("", 3022, "", 22, proto.NetworkType_TCP).Start()

	go transfer.NewTransferWithIPPort("", 4022, "", 22, proto.NetworkType_TCP).Start()

	go transfer.NewTransferWithIPPort("", 5555, "", 4444, proto.NetworkType_TCP).Start()

	srv := service.NewService("", 12341)
	go srv.Start()

	utils.IgnoreSignal(func() {
		srv.Close()
		transfer.Manage.Close()
		log.Println("service", "退出")
	})
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/gin-gonic/gin"
	_ "github.com/ranxx/cuten"
	"github.com/ranxx/goproxy/api"
	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	transfer "github.com/ranxx/goproxy/transfer/service"
	"github.com/ranxx/goproxy/utils"
)

func main() {
	fmt.Println("Hello World")

	// test 开启转发
	go transfer.NewTransferWithIPPort("", 3334, "", 3333, proto.NetworkType_HTTP).Start()

	go transfer.NewTransferWithIPPort("", 3335, "", 3333, proto.NetworkType_HTTP).Start()

	go transfer.NewTransferWithIPPort("", 3336, "", 3333, proto.NetworkType_HTTP).Start()

	go transfer.NewTransferWithIPPort("", 2022, "", 22, proto.NetworkType_TCP).Start()

	go transfer.NewTransferWithIPPort("", 3022, "", 22, proto.NetworkType_TCP).Start()

	go transfer.NewTransferWithIPPort("", 4022, "", 22, proto.NetworkType_TCP).Start()

	go transfer.NewTransferWithIPPort("", 5555, "", 4444, proto.NetworkType_TCP).Start()

	srv := service.NewService("", 12341)
	go srv.Start()

	go func() {
		time.Sleep(time.Second * 3)
		e := api.Init()
		log.Println("starting api", "port :12351")
		log.Println(http.ListenAndServe(":12351", e))
	}()

	utils.IgnoreSignal(func() {
		srv.Close()
		transfer.Manage.Close()
		log.Println("service", "退出")
	})
}

package main

import (
	"flag"
	"log"

	"github.com/ranxx/goproxy/service"
	transfer "github.com/ranxx/goproxy/transfer/client"
	"github.com/ranxx/goproxy/utils"
)

// ip port
var ip string
var port int

func init() {
	flag.StringVar(&ip, "ip", "", "remote ip")
	flag.IntVar(&port, "port", service.DefaultPort, "remote port")
	flag.Parse()
}

func main() {
	// 49.233.211.140
	srv := service.NewClient(ip, port)
	go srv.Start()

	utils.IgnoreSignal(func() {
		srv.Close()
		transfer.Manage.Close()
		log.Println("client", "退出")
	})
}

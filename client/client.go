package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ranxx/goproxy/service"
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
	srv := service.NewClient(ip, port)
	go srv.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
			srv.Close()
			log.Println("client退出")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

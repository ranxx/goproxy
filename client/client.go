package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ranxx/goproxy/service"
)

func main() {
	fmt.Println("Hello World")

	srv := service.NewClient("", 12341)
	go srv.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
			srv.Close()
			fmt.Println("client退出")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

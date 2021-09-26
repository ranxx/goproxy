package main

import (
	"fmt"

	"github.com/ranxx/goproxy/service"
)

func main() {
	fmt.Println("Hello World")

	srv := service.NewClient("", 12341)
	srv.Start()
	select {}
}

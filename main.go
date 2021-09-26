package main

import (
	"fmt"
	"io"
	"net"

	proto3 "github.com/golang/protobuf/proto"
	"github.com/ranxx/goproxy/proto"
)

func test() {
	msg := proto.Msg{
		Network: "sds",
		Body:    []byte("Hello World"),
	}
	body, err := proto3.Marshal(&msg)
	fmt.Println(err, len(body), string(body))
	body, err = msg.XXX_Marshal(nil, false)
	fmt.Println(err, len(body), string(body))
	body, err = msg.XXX_Marshal(nil, true)
	fmt.Println(err, len(body), string(body))

}

func main() {

	test()
	return
	listenr, err := net.Listen("tcp", ":3022")
	if err != nil {
		panic(err)
	}
	defer listenr.Close()
	for {
		conn, err := listenr.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			defer func() {
				fmt.Println(recover())
			}()
			connection(conn)
		}()
	}
}

func connection(inConn net.Conn) {
	// 连接 远端 ip:port
	outConn, err := net.Dial("tcp", "49.233.211.140:22")
	if err != nil {
		panic(err)
	}

	// go ioCopy(inConn, outConn)

	// ioCopy(outConn, inConn)
	// inConn.Close()
	// outConn.Close()
	go readConn(inConn, outConn)
	go readConn(outConn, inConn)
}

func readConn(dst, src net.Conn) {
	ioCopy(dst, src)
	src.Close()
}

func ioCopy(dst io.WriteCloser, src io.ReadCloser) {
	for {
		wn, we := io.Copy(dst, src)
		if wn == 0 && we == nil {
			fmt.Println("退出")
			break
		}
		if we == io.EOF {
			fmt.Println("退出")
			// src.Close()
			break
		}
		fmt.Println(wn, we)
		break
	}
}

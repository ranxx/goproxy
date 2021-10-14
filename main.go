package main

import (
	"fmt"
	"io"
	"net"

	proto3 "github.com/golang/protobuf/proto"
	"github.com/ranxx/goproxy/config"
	"github.com/ranxx/goproxy/proto"
)

func testConfig() {
	cfg := config.ParseYamlFile("./config/config.yaml")
	fmt.Println(cfg)
}

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

func testpack() {
	body := &proto.TCPBody{
		MsgId: 2,
		Laddr: &proto.Addr{Ip: "12", Port: 12},
		Raddr: &proto.Addr{Ip: "23", Port: 23},
		Body:  []byte("H"),
	}

	nbody := new(proto.TCPBody)

	fmt.Println(body)
	fmt.Println(nbody)

}

func main() {
	testConfig()
	return
	testpack()
	return
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

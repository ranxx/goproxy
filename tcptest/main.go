package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

func main() {
	dial()
	return
	listenr, err := net.Listen("tcp", ":4444")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listenr.Accept()
		if err != nil {
			panic(err)
		}
		go echo(conn)
	}
}

func dial() {
	conn, err := net.Dial("tcp", ":4444")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	go func() {
		return
		go func() {
			index := 0
			scann := bufio.NewScanner(conn)
			scann.Buffer(nil, 1024*1024*1024)
			scann.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
				index++
				log.Println("dial", index, atEOF, len(data), cap(data))
				if atEOF && len(data) == 0 {
					return 0, nil, nil
				}
				if !atEOF {
					if len(data) < 6799801 {
						return 0, nil, nil
					}
					// 读取完成
					log.Println("读取完成")
					return 6799801, data[:6799801], nil
				}
				return 0, nil, nil
			})
			for scann.Scan() {
				log.Println("dial", "读取消息")
				bytess := scann.Bytes()
				// bytess, err := ioutil.ReadAll(conn)
				rn := len(bytess)
				//rn, err := conn.Read(bytess)
				// if err != nil {
				// 	panic(err)
				// }
				log.Println("读取", rn)
				conn.Write(bytess[:rn])
				time.Sleep(time.Second)
			}
			// 退出
			log.Println("echo", "退出", scann.Err())
		}()
	}()
	for {
		//       1323148
		body := make([]byte, 6799801)
		// _, err := conn.Write([]byte("Hello Wrold " + time.Now().String()))
		wn, err := conn.Write(body)
		if err != nil {
			panic(err)
		}
		log.Println("写入消息", wn)
		// bytess := make([]byte, 6799801)
		// rn, err := conn.Read(bytess)
		// if err != nil {
		// 	panic(err)
		// }
		// log.Println(string(bytess[:rn]))
		time.Sleep(time.Second)
	}
}

// 512读取

func echo(conn net.Conn) {
	// conn.SetNoDelay(true)
	go func() {
		index := 0
		scann := bufio.NewScanner(conn)
		scann.Buffer(nil, 1024*1024*1024)
		scann.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			index++
			log.Println("echo", index, atEOF, len(data), cap(data))
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if !atEOF {
				if len(data) < 6799801 {
					return 0, nil, nil
				}
				// 读取完成
				log.Println("读取完成")
				return 6799801, data[:6799801], nil
			}
			return 0, nil, nil
		})
		for scann.Scan() {
			// log.Println("echo", "读取消息")
			bytess := scann.Bytes()
			// bytess, err := ioutil.ReadAll(conn)
			rn := len(bytess)
			//rn, err := conn.Read(bytess)
			// if err != nil {
			// 	panic(err)
			// }
			log.Println("读取", rn)
			// conn.Write(bytess[:rn])
			time.Sleep(time.Second)
		}
		// 退出
		log.Println("echo", "退出", scann.Err())
	}()
}

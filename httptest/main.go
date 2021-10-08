package main

import (
	"fmt"
	"net/http"
	"syscall"
	"time"
)

func main() {
	// test()
	// return
	// server := http.Server{
	// 	Handler:      new(httpServer),
	// 	ReadTimeout:  20 * time.Second,
	// 	WriteTimeout: 20 * time.Second,
	// }
	// server.SetKeepAlivesEnabled(false)
	// listen, err := net.Listen("tcp4", ":3333")
	// if err != nil {
	// 	log.Printf("Failed to listen,err:%s\n", err.Error())
	// 	panic(err)
	// }
	// fmt.Println(server.Serve(listen))

	http.ListenAndServe(":3333", http.FileServer(http.Dir("/Users/axing")))
	// http.ListenAndServe(":3333", new(httpServer))
}

type httpServer struct {
}

//  ServeHTTP ...
func (h *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Close = true
	msg := fmt.Sprintf("Hello World: %s\nWelcone to %s\nnow time: %s\n", r.RemoteAddr, r.URL.String(), time.Now().Format("2006-01-02 15:04:05"))
	w.Write([]byte(msg))
}

func test() {

	resp := syscall.Rlimit{}
	fmt.Println(syscall.Getrlimit(syscall.RLIMIT_AS, &resp))
	fmt.Println(resp)
}

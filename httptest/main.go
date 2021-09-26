package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.ListenAndServe(":3333", http.FileServer(http.Dir("/Users/axing")))
}

type httpServer struct {
}

//  ServeHTTP ...
func (h *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Hello World: %s\nWelcone to %s\nnow time: %s\n", r.RemoteAddr, r.URL.String(), time.Now().Format("2006-01-02 15:04:05"))
	w.Write([]byte(msg))
}

package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
)

// HTTP ...
type HTTP struct {
	logPrefix string
	msg       *proto.Msg
	body      *proto.HTTPBody
	once      *sync.Once
}

// newHTTP 新建http
func newHTTP(logPrefix string, msg *proto.Msg, body *proto.HTTPBody) *HTTP {
	return &HTTP{
		logPrefix: logPrefix,
		msg:       msg,
		body:      body,
		once:      new(sync.Once),
	}
}

// Start ...
func (h *HTTP) Start() {
	h.Receive(&h.body.Body)
}

// Close ...
func (h *HTTP) Close() {}

// Receive ...
func (h *HTTP) Receive(body *[]byte) {
	request, err := http.NewRequest(h.body.Method, fmt.Sprintf("%s:%d/%s", "http://localhost", h.body.Laddr.Port, h.body.Url), bytes.NewReader(*body))
	if err != nil {
		log.Println(h.logPrefix, "初始化request失败", err)
		return
	}
	defer request.Body.Close()

	for _, header := range h.body.Header {
		for _, value := range header.Value {
			request.Header.Set(header.Key, value)
		}
	}
	request.Close = true

	cli := &http.Client{Transport: &http.Transport{DisableKeepAlives: true, MaxIdleConns: 0, MaxIdleConnsPerHost: 0}}
	response, err := cli.Do(request)
	if err != nil {
		log.Println(h.logPrefix, "请求http失败", err)
		return
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(h.logPrefix, "读取http返回值失败", err)
		return
	}

	h.body.Body = resBody
	h.body.Header = make([]*proto.Header, 0, len(response.Proto))
	for key, values := range response.Header {
		h.body.Header = append(h.body.Header, &proto.Header{Key: key, Value: values})
	}

	rbody, err := h.body.XXX_Marshal(nil, false)
	if err != nil {
		log.Println(h.logPrefix, "编码失败", err)
		return
	}

	service.WritingMsgChannel <- &proto.Msg{
		Network: h.msg.Network,
		MsgId:   h.msg.MsgId,
		Body:    rbody,
	}
}

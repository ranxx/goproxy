package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	"github.com/ranxx/goproxy/utils"
)

// HTTP ...
type HTTP struct {
	network      proto.NetworkType
	logPrefix    string
	msgID        int64
	indexManage  *utils.IndexI64
	indexs       []*httpBody
	laddr, raddr proto.Addr
	mutex        sync.Mutex
}

type httpBody struct {
	*proto.HTTPBody
	receiveBody chan *proto.HTTPBody
}

// newHTTP ...
func newHTTP(logPrefix string, network proto.NetworkType, msgID int64, localAddr, remoteAddr proto.Addr) Transfer {
	return &HTTP{
		network:     network,
		logPrefix:   fmt.Sprintf("%s %s", logPrefix, utils.TunnelAddrInfo(&localAddr, &remoteAddr)),
		msgID:       msgID,
		indexs:      make([]*httpBody, 512),
		laddr:       localAddr,
		raddr:       remoteAddr,
		indexManage: utils.NewIndexI64(),
	}
}

// Receive ...
func (h *HTTP) Receive(body *[]byte) {
	go func() {
		if body == nil {
			return
		}
		httpBody := new(proto.HTTPBody)
		if err := httpBody.XXX_Unmarshal(*body); err != nil {
			log.Println(h.logPrefix, "解码返回的消息失败", err)
			panic(err)
		}
		log.Println(h.logPrefix, fmt.Sprintf("msgId:%d 收到 %d 返回的数据", h.msgID, httpBody.MsgId))
		h.indexs[httpBody.MsgId].receiveBody <- httpBody
	}()
}

// Start ...
func (h *HTTP) Start() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", h.laddr.Ip, h.laddr.Port))
	if err != nil {
		return err
	}

	log.Println(h.logPrefix, "runing")
	go http.Serve(listen, h)
	// err := http.ListenAndServe(fmt.Sprintf("%s:%d", h.laddr.Ip, h.laddr.Port), h)
	// if err != nil {
	// 	log.Println(h.logPrefix, "监听 %s 失败", utils.AddrString(&h.laddr), err)
	// }
	return nil
}

// Close 关闭
func (h *HTTP) Close() {}

func (h *HTTP) send(body *proto.HTTPBody) (*httpBody, error) {
	sendBody, err := body.XXX_Marshal(nil, false)
	if err != nil {
		log.Println(h.logPrefix, "编码httpbody失败", err)
		return nil, err
	}

	hbody := &httpBody{HTTPBody: body, receiveBody: make(chan *proto.HTTPBody, 1)}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.indexs = append(h.indexs, nil)
	h.indexs[body.MsgId] = hbody

	service.GlobalWritingMsgChannel <- &proto.Msg{
		Network: proto.NetworkType_HTTP.String(),
		MsgId:   h.msgID,
		Body:    sendBody,
	}

	return hbody, nil
}

func (h *HTTP) newProtoHTTPBody(index int64, url, method string, header http.Header, body []byte) *proto.HTTPBody {
	httpBody := proto.HTTPBody{
		MsgId: index,
		Laddr: &h.raddr,
		// TODO: url 参数也要携带
		Url:    url,
		Method: method,
		Header: make([]*proto.Header, len(header)),
		Body:   body,
	}

	for key, values := range header {
		httpBody.Header = append(httpBody.Header, &proto.Header{Key: key, Value: values})
	}
	httpBody.Header = append(httpBody.Header, &proto.Header{Key: "Connection", Value: []string{"close"}})

	return &httpBody
}

func (h *HTTP) oneHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request, index int64) {
	// 读取 body
	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(h.logPrefix, "读取body失败", err)
		return
	}

	pHTTPBody := h.newProtoHTTPBody(index, r.URL.Path, r.Method, r.Header, bytesBody)

	httpBody, err := h.send(pHTTPBody)
	if err != nil {
		log.Println(h.logPrefix, "发送消息失败", err)
		return
	}
	respBody := <-httpBody.receiveBody

	// 设置返回header
	for _, v := range respBody.Header {
		for _, value := range v.Value {
			w.Header().Add(v.Key, value)
		}
	}

	// 写回消息
	wn, err := w.Write(respBody.Body)
	if err != nil {
		log.Println(h.logPrefix, "写入body失败", err)
		return
	}

	if wn != len(respBody.Body) {
		log.Println(h.logPrefix, "消息未写入完整", err)
		return
	}
}

// NetWork ...
func (h *HTTP) NetWork() proto.NetworkType {
	return h.network
}

// ServeHTTP ...
func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// index
	h.oneHTTP(context.TODO(), w, r, h.indexManage.NewIndex())
}

// Info ...
func (h *HTTP) Info() Info {
	return Info{
		Index:   h.msgID,
		Laddr:   h.laddr,
		Raddr:   h.raddr,
		Machine: "",
	}
}

package transfer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/service"
	"github.com/ranxx/goproxy/utils"
)

// HTTP ...
type HTTP struct {
	MsgID                 int64
	IndexManage           *utils.IndexI64
	Indexs                []*httpBody
	ReceiveBody           chan *proto.HTTPBody
	LocalAddr, RemoteAddr proto.Addr
	mutex                 sync.Mutex
}

type httpBody struct {
	*proto.HTTPBody
	receiveBody chan *proto.HTTPBody
}

// newHTTP ...
func newHTTP(msgID int64, localAddr, remoteAddr proto.Addr) Transfer {
	return &HTTP{
		MsgID:       msgID,
		Indexs:      make([]*httpBody, 512),
		ReceiveBody: make(chan *proto.HTTPBody, 1024*8),
		LocalAddr:   localAddr,
		RemoteAddr:  remoteAddr,
		IndexManage: utils.NewIndexI64(),
	}
}

// Receive ...
func (h *HTTP) Receive(body []byte) {
	go func() {
		httpBody := new(proto.HTTPBody)
		if err := httpBody.XXX_Unmarshal(body); err != nil {
			log.Println("transfer.http", "解码返回的消息失败", err)
			panic(err)
		}
		log.Println("transfer.http", fmt.Sprintf("msgID:%d 收到 %d 返回的数据", h.MsgID, httpBody.MsgId))
		h.Indexs[httpBody.MsgId].receiveBody <- httpBody
	}()
}

// Start ...
func (h *HTTP) Start() {
	log.Println("transfer.http", fmt.Sprintf("%s:%d -> %s:%d runing", h.LocalAddr.Ip, h.LocalAddr.Port, h.RemoteAddr.Ip, h.RemoteAddr.Port))
	http.ListenAndServe(fmt.Sprintf("%s:%d", h.LocalAddr.Ip, h.LocalAddr.Port), h)
}

func (h *HTTP) send(body *proto.HTTPBody) *httpBody {
	sendBody, err := body.XXX_Marshal(nil, false)
	if err != nil {
		log.Println("transfer.http", "编码httpbody失败", err)
		panic(err)
	}

	hbody := &httpBody{HTTPBody: body, receiveBody: make(chan *proto.HTTPBody, 1)}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.Indexs = append(h.Indexs, nil)
	h.Indexs[body.MsgId] = hbody

	service.WritingMsgChannel <- &proto.Msg{
		Network: proto.NetworkType_HTTP.String(),
		MsgId:   h.MsgID,
		Body:    sendBody,
	}
	return hbody
}

func (h *HTTP) newProtoHTTPBody(index int64, url, method string, header http.Header, body []byte) *proto.HTTPBody {
	httpBody := proto.HTTPBody{
		MsgId: index,
		Laddr: &h.RemoteAddr,
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
		log.Println("transfer.http", "读取body失败", err)
		panic(err)
	}

	pHTTPBody := h.newProtoHTTPBody(index, r.URL.Path, r.Method, r.Header, bytesBody)

	httpBody := h.send(pHTTPBody)
	respBody := <-httpBody.receiveBody

	// 设置返回header
	for _, v := range respBody.Header {
		for _, value := range v.Value {
			w.Header().Add(v.Key, value)
		}
	}

	// 写入消息
	wn, err := w.Write(respBody.Body)
	if err != nil {
		log.Println("transfer.http", "写入body失败", err)
		panic(err)
	}

	if wn != len(respBody.Body) {
		log.Println("transfer.http", "写入body失败,消息未写入完整", err)
		panic(err)
	}
}

// ServeHTTP ...
func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// index
	h.oneHTTP(context.TODO(), w, r, h.IndexManage.NewIndex())
}

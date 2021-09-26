package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/ranxx/goproxy/pack"
	"github.com/ranxx/goproxy/proto"
)

// var
var (
	// 接收的消息
	ReadingMsgChannel = make(chan *proto.Msg, 1024)

	// 发送的消息
	WritingMsgChannel = make(chan *proto.Msg, 1024)
)

// // Msg ...
// type Msg struct {
// 	Network string      `json:"network"`
// 	Body    interface{} `json:"body"`
// }

// // HTTPBody ...
// type HTTPBody struct {
// 	Laddr  Addr
// 	URL    string
// 	Hander http.Header
// 	Body   []byte
// }

// // TCPBody ...
// type TCPBody struct {
// 	Raddr, Laddr Addr
// 	Body         []byte
// }

// // IsHTTP ...
// func (msg *Msg) IsHTTP() bool {
// 	return strings.ToLower(msg.Network) == "http"
// }

// // PackBytes ...
// func (msg *Msg) PackBytes() ([]byte, error) {
// 	buffer := bytes.NewBuffer(make([]byte, 0, 1024*4))
// 	return buffer.Bytes(), msg.Pack(buffer)
// }

// // Pack ...
// func (msg *Msg) Pack(writer io.Writer) error {
// 	var err error
// 	err = binary.Write(writer, binary.BigEndian, &msg.Network)
// 	err = binary.Write(writer, binary.BigEndian, &msg.Body)
// 	return err
// }

// // UnpackBytes ...
// func (msg *Msg) UnpackBytes(body []byte) error {
// 	buffer := bytes.NewBuffer(body)
// 	return msg.Unpack(buffer)
// }

// // Unpack ...
// func (msg *Msg) Unpack(reader io.Reader) error {
// 	var err error
// 	err = binary.Read(reader, binary.BigEndian, &msg.Network)
// 	err = binary.Read(reader, binary.BigEndian, &msg.Body)
// 	return err
// }

// CheckHTTP ...
func CheckHTTP(msg *proto.Msg) bool {
	return proto.NetworkType_value[strings.ToUpper(msg.Network)] == int32(proto.NetworkType_HTTP)
}

// CheckTCP ...
func CheckTCP(msg *proto.Msg) bool {
	return proto.NetworkType_value[strings.ToUpper(msg.Network)] == int32(proto.NetworkType_TCP)
}

// WriteFunc ...
func WriteFunc(c *Conn) error {
	for msg := range WritingMsgChannel {
		body, err := msg.XXX_Marshal(nil, false)
		if err != nil {
			log.Println("msg打包失败", err)
			return err
		}
		// log.Println("service", "开始回写", len(body), string(body))

		body, err = pack.NewPackage(body).PackBytes()
		if err != nil {
			log.Println("pack打包失败", err)
			return err
		}

		log.Println("service", "开始回写", len(msg.Body), string(msg.Body))

		fmt.Println(c.Write(body))
	}
	return nil
}

// ReadFunc ...
func ReadFunc(c *Conn) error {
	scanner := pack.NewScanner(c)
	for scanner.Scan() {
		scannedPack := new(pack.Package)
		err := scannedPack.UnpackBytes(scanner.Bytes())
		if err != nil {
			log.Println("解包pack失败", err, scanner.Bytes())
			return err
		}

		msg := new(proto.Msg)
		err = msg.XXX_Unmarshal(scannedPack.Msg)
		if err != nil {
			log.Println("解包msg失败", err, scanner.Bytes())
			return err
		}

		log.Println("service", "开始读取", len(msg.Body), string(msg.Body))

		ReadingMsgChannel <- msg
	}
	return nil
}

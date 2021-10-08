package service

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ranxx/goproxy/cconn"
)

// Client 客户端
// 具有重连机制，重连之后所有连接会全部断开
type Client struct {
	IP    string
	Port  int
	Conn  *cconn.Conn
	mutex sync.Mutex
	once  *sync.Once
	// tcpConns [][]net.Conn
	// 初始化多个channel
	customers [65535][]ICustomer
	// customersV2 [65535][]ICustomer
	tcpConns [65535][]net.Conn
}

// NewClient ...
func NewClient(ip string, port int) *Client {
	return &Client{IP: ip, Port: port, once: new(sync.Once)}
}

func (c *Client) reset() {
	c.Conn = nil
	c.once = &sync.Once{}
	// c.customers = make([65535][]ICustomer)
	// ReadingMsgChannel = make(chan *proto.Msg, 1024)
	// WritingMsgChannel = make(chan *proto.Msg, 1024)
}

// Close ...
func (c *Client) Close() {
	c.once.Do(func() {
		if c.Conn != nil {
			c.Conn.Close()
			c.Conn = nil
		}
		for _, conns := range c.customers {
			for _, conn := range conns {
				conn.Close()
			}
		}
		for _, conns := range c.tcpConns {
			for _, conn := range conns {
				fmt.Println("关闭连接")
				conn.Close()
			}
		}
		// close(ReadingMsgChannel)
		// close(WritingMsgChannel)
	})
}

func (c *Client) dail() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.IP, c.Port), time.Second*10)
	return conn, err
}

// CheckConn 检查连接
func (c *Client) CheckConn() bool {
	if c.Conn == nil {
		return true
	}
	if c.Conn.Closed() {
		return true
	}
	return false
}

// CleanUpOnce ...
func (c *Client) CleanUpOnce() {
	c.Close()
}

// ReConnection 重连
func (c *Client) ReConnection() error {
	// 重连
	conn, err := c.dail()
	if err != nil {
		log.Println("client", err)
		return err
	}

	c.reset()

	log.Println("client", fmt.Sprintf("连上server %s -> %s", conn.LocalAddr(), conn.RemoteAddr()))

	c.Conn = cconn.NewConn(
		"client",
		conn,
		cconn.WithReadFunc(_DefaultReadFunc("client")),
		cconn.WithWriteFunc(_DefaultWriteFunc("client")))
	c.StartConn()
	return nil
}

// ReStart ...
func (c *Client) ReStart() {
	go c.Customer()
}

// Start ...
func (c *Client) Start() {
	cconn.Checking(c, -1)
}

// StartConn ...
func (c *Client) StartConn() {
	// 开启读写
	c.Conn.Start()
	return
}

// Customer ...
func (c *Client) Customer() {
	// log.Println("client", "开始消费")
	// for v := range ReadingMsgChannel {
	// 	if v == nil {
	// 		break
	// 	}
	// 	if CheckHTTP(v) {
	// 		go c.HTTPHandler(v)
	// 		continue
	// 	}
	// 	// tcp，主动连接
	// 	go c.tunnelTCP(v)
	// }
	// log.Println("client", "退出消费")
}

// // HTTPHandler ...
// func (c *Client) HTTPHandler(msg *proto.Msg) {
// 	// 先解析
// 	body := new(proto.HTTPBody)
// 	if err := body.XXX_Unmarshal(msg.Body); err != nil {
// 		log.Println("client", "处理http 解析body失败", err)
// 		panic(err)
// 	}

// 	// log.Println("client", "读取到service的消息", string(body.Body))

// 	// 转发
// 	request, err := http.NewRequest(body.Method, fmt.Sprintf("%s:%d%s", "http://localhost", body.Laddr.Port, body.Url), bytes.NewReader(body.Body))
// 	if err != nil {
// 		log.Println("client", "处理http 转发http失败", err)
// 		panic(err)
// 	}
// 	request.Close = true
// 	defer request.Body.Close()

// 	for _, header := range body.Header {
// 		for _, v := range header.Value {
// 			request.Header.Add(header.Key, v)
// 		}
// 	}

// 	cli := &http.Client{Transport: &http.Transport{DisableKeepAlives: true, MaxIdleConns: 0, MaxIdleConnsPerHost: 0}}
// 	response, err := cli.Do(request)
// 	if err != nil {
// 		log.Println("client", "处理http 请求http失败", err)
// 		panic(err)
// 	}

// 	bbody, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		log.Println("client", "处理http 读取http返回值失败", err)
// 		panic(err)
// 	}
// 	response.Body.Close()

// 	body.Body = bbody
// 	body.Header = make([]*proto.Header, 0, len(response.Proto))
// 	for key, values := range response.Header {
// 		body.Header = append(body.Header, &proto.Header{Key: key, Value: values})
// 	}

// 	wbody, err := body.XXX_Marshal(nil, false)
// 	if err != nil {
// 		log.Println("client", "处理http 编码返回body失败", err)
// 		panic(err)
// 	}

// 	WritingMsgChannel <- &proto.Msg{
// 		Network: proto.NetworkType_HTTP.String(),
// 		MsgId:   msg.MsgId,
// 		Body:    wbody,
// 	}

// 	return
// }

// // TCPHandler ...
// // 主动连接
// // func (c *Client) TCPHandler(msg *proto.Msg) {
// // 	// 先解析
// // 	body := new(proto.TCPBody)
// // 	if err := body.XXX_Unmarshal(msg.Body); err != nil {
// // 		log.Println("client.tcp", "解析body失败", err)
// // 		panic(err)
// // 	}

// // 	c.mutex.Lock()
// // 	defer c.mutex.Unlock()

// // 	if body.MsgId >= int64(len(c.tcpConns[msg.MsgId])) {
// // 		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", body.Laddr.Ip, body.Laddr.Port))
// // 		if err != nil {
// // 			log.Println("client.tcp", fmt.Sprintf("连接 %s:%d 失败", body.Laddr.Ip, body.Laddr.Port), err)
// // 			panic(err)
// // 		}
// // 		log.Println("client", fmt.Sprintf("新建tcp %s -> %s", conn.LocalAddr(), conn.RemoteAddr()))

// // 		c.tcpConns[msg.MsgId] = append(c.tcpConns[msg.MsgId], conn)
// // 		c.tcpConns[msg.MsgId][body.MsgId] = conn

// // 		scanner := bufio.NewScanner(conn)
// // 		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
// // 			log.Println("servcie.client.tcp", atEOF, len(data))
// // 			if atEOF && len(data) == 0 {
// // 				return 0, nil, nil
// // 			}
// // 			if !atEOF {
// // 				return len(data), data[:], nil
// // 			}
// // 			return 0, nil, nil
// // 		})
// // 		go func() {
// // 			for scanner.Scan() {
// // 				// 开启读
// // 				rbody := scanner.Bytes()

// // 				body.Body = rbody
// // 				body.Laddr, body.Raddr = body.Raddr, body.Laddr

// // 				tcpBody, err := body.XXX_Marshal(nil, false)
// // 				if err != nil {
// // 					log.Println("client.tcp", "编码body失败", err)
// // 					panic(err)
// // 				}
// // 				WritingMsgChannel <- &proto.Msg{
// // 					Network: proto.NetworkType_TCP.String(),
// // 					MsgId:   msg.MsgId,
// // 					Body:    tcpBody,
// // 				}
// // 			}
// // 			log.Println("client", fmt.Sprintf("关闭tcp %s -> %s", conn.LocalAddr(), conn.RemoteAddr()))
// // 		}()
// // 	}

// // 	// 开启连接写
// // 	c.tcpConns[msg.MsgId][body.MsgId].Write(body.Body)
// // }

// // TCPHandlerV2 ...
// // 主动连接
// func (c *Client) TCPHandlerV2(msg *proto.Msg) {
// 	// 主动连接

// 	// 先解析
// 	body := new(proto.TCPBody)
// 	if err := body.XXX_Unmarshal(msg.Body); err != nil {
// 		log.Println("client.tcp", "解析body失败", err)
// 		panic(err)
// 	}

// 	// 先判断有
// 	if body.MsgId < int64(len(c.customers[msg.MsgId])) {
// 		c.customers[msg.MsgId][body.MsgId].Receive(&body.Body)
// 		return
// 	}

// 	c.mutex.Lock()
// 	defer c.mutex.Unlock()
// 	if body.MsgId >= int64(len(c.customers[msg.MsgId])) {
// 		log.Println("client.tcp", msg.MsgId, body.MsgId)
// 		tcp := customer.NewTCP("client.tcp", WritingMsgChannel, msg, body)
// 		c.customers[msg.MsgId] = append(c.customers[msg.MsgId], tcp)
// 		c.customers[msg.MsgId][body.MsgId] = tcp
// 		c.customers[msg.MsgId][body.MsgId].Start()
// 	}

// 	c.customers[msg.MsgId][body.MsgId].Receive(&body.Body)
// }

// func (c *Client) tunnelTCP(msg *proto.Msg) {
// 	// 先解析TCP
// 	body := new(proto.TCPBody)
// 	if err := body.XXX_Unmarshal(msg.Body); err != nil {
// 		log.Println("client.tcp", "解析body失败", err)
// 		panic(err)
// 	}

// 	// 连接本地
// 	inConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", body.Laddr.Ip, body.Laddr.Port))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 连接远程
// 	outConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", body.Raddr.Ip, body.Raddr.Port))
// 	if err != nil {
// 		panic(err)
// 	}
// 	// 报告自己的 tunnel

// 	// outConn.LocalAddr()

// 	// nbody := proto.TCPBody{
// 	// 	MsgId: body.MsgId,
// 	// 	// Raddr:,
// 	// }

// 	// nbodyBytes, err := nbody.XXX_Marshal(nil, false)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// reflect.Copy(reflect.ValueOf(nbody), reflect.ValueOf(body))

// 	// rnbodyBytes, err := pack.NewPackage(nbodyBytes).PackBytes()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// outConn.Write(rnbodyBytes)

// 	// ipss := outConn.LocalAddr().String()
// 	// index := strings.LastIndex(ipss, ":")

// 	// fmt.Println(ipss[:index], ipss[index+1:])

// 	fmt.Println(outConn.LocalAddr().String(), outConn.RemoteAddr().String())
// 	fmt.Println(inConn.LocalAddr().String(), inConn.RemoteAddr().String())

// 	// lport, err := strconv.Atoi(ipss[index+1:])
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	nbody := proto.Bind{
// 		MsgId: body.MsgId,
// 		// Laddr: &proto.Addr{Ip: ipss[:index], Port: },
// 	}

// 	nbodyBytes, err := nbody.XXX_Marshal(nil, false)
// 	if err != nil {
// 		panic(err)
// 	}

// 	nnBytes, err := pack.NewPackage(nbodyBytes).PackBytes()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// WritingMsgChannel <- &proto.Msg{
// 	// 	Network: proto.NetworkType_BIND.String(),
// 	// 	MsgId:   msg.MsgId,
// 	// 	Body:    nbodyBytes,
// 	// }

// 	fmt.Println(outConn.Write(nnBytes))
// 	// fmt.Println(inConn.Write(body.Body))

// 	go io.Copy(inConn, outConn)
// 	go io.Copy(outConn, inConn)
// 	c.tcpConns[msg.MsgId] = append(c.tcpConns[msg.MsgId], outConn)
// 	c.tcpConns[msg.MsgId] = append(c.tcpConns[msg.MsgId], inConn)
// }

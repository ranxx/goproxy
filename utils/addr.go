package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ranxx/goproxy/proto"
)

// TunnelAddrInfo ...
func TunnelAddrInfo(laddr, raddr *proto.Addr) string {
	return fmt.Sprintf("%s:%d -> %s:%d", laddr.Ip, laddr.Port, raddr.Ip, raddr.Port)
}

// AddrString ...
func AddrString(addr *proto.Addr) string {
	return fmt.Sprintf("%s:%d", addr.Ip, addr.Port)
}

// ParseAddrString ...
func ParseAddrString(addr string) (ip string, port int, err error) {
	index := strings.LastIndex(addr, ":")
	if index < 0 || index == len(addr)-1 {
		return "", 0, fmt.Errorf("不支持的格式")
	}
	port, err = strconv.Atoi(addr[index+1:])
	if err != nil {
		return
	}

	ip = addr[:index]
	return
}

// RetRealHTTP 返回带有http或者https
func RetRealHTTP(ip string) string {
	if strings.HasPrefix(ip, "https://") {
		return ip
	}
	if strings.HasPrefix(ip, "http://") {
		return ip
	}
	if len(ip) <= 0 {
		return "http://localhost"
	}
	return fmt.Sprintf("http://%s", ip)
}

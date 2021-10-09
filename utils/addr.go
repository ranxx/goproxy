package utils

import (
	"fmt"

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

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ranxx/goproxy/api"
	"github.com/ranxx/goproxy/proto"
	"github.com/ranxx/goproxy/utils"
	"github.com/ranxx/grequests"
	"github.com/spf13/cobra"
)

type addCommand struct {
	laddrs  *[]string
	raddr   *string
	machine *string
}

func (a *addCommand) Cmd() *cobra.Command {
	add := &cobra.Command{
		Use:   "add",
		Short: "添加穿透的ip:port 列表",
		Long: `内网穿透服务端，需要指定被穿透的机器。
eg:
1) goproxy server nat add --laddr ip:port --raddr ip:port`,
		Run: a.Run,
	}

	a.laddrs = add.Flags().StringSliceP("laddr", "l", nil, "localhost addr")
	a.raddr = add.Flags().StringP("raddr", "r", "", "remote addr")
	add.MarkFlagRequired("laddr")
	add.MarkFlagRequired("raddr")
	return add
}

func (a *addCommand) Run(cmd *cobra.Command, args []string) {
	l := make([]proto.Addr, 0, len(*a.laddrs))
	for _, laddr := range *a.laddrs {
		laddrIP, laddrPort, err := utils.ParseAddrString(laddr)
		if err != nil {
			panic(err)
		}
		l = append(l, proto.Addr{Ip: laddrIP, Port: int32(laddrPort)})
	}
	raddrIP, raddrPort, err := utils.ParseAddrString(*a.raddr)
	if err != nil {
		panic(err)
	}
	r := proto.Addr{Ip: raddrIP, Port: int32(raddrPort)}

	url := fmt.Sprintf("%s:%d/transfer/tcp", "http://localhost", 12351)

	req := struct {
		Laddr []proto.Addr `json:"laddr"`
		Raddr proto.Addr   `json:"raddr"`
	}{
		Laddr: l, Raddr: r,
	}

	resp := api.Message{}
	err = grequests.Post(context.TODO(), url, req, &resp)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Code, resp.Msg, resp.Data)
}

type removeCommand struct{}

func (r *removeCommand) Cmd() *cobra.Command {
	remove := &cobra.Command{
		Use:   "rm",
		Short: "移除监听的端口",
		Long: `内网穿透服务端，按端口移除transfer
eg:
1) goproxy server rm port [port...]`,
		Run: r.Run,
	}
	return remove
}

func (r *removeCommand) Run(cmd *cobra.Command, args []string) {
	// 端口
	// 处理args
	ports := make([]int, 0, len(args))
	for _, v := range args {
		port, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		ports = append(ports, port)
	}

	url := fmt.Sprintf("%s:%d/transfer/port", "http://localhost", 12351)

	req := struct {
		Ports []int `json:"ports"`
	}{
		Ports: ports,
	}

	resp := api.Message{}
	err := grequests.Delete(context.TODO(), url, req, &resp)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Code, resp.Msg, resp.Data)
}

type listCommand struct{}

func (l *listCommand) Cmd() *cobra.Command {
	remove := &cobra.Command{
		Use:   "list",
		Short: "列表transfer",
		Long: `内网穿透服务端，列表transfer
eg:
1) goproxy server list`,
		Run: l.Run,
	}
	return remove
}

func (l *listCommand) Run(cmd *cobra.Command, args []string) {

	url := fmt.Sprintf("%s:%d/transfer", "http://localhost", 12351)

	resp := api.Message{}
	err := grequests.Get(context.TODO(), url, nil, &resp)
	if err != nil {
		panic(err)
	}
	if resp.Data == nil {
		fmt.Println(resp.Code, resp.Msg, resp.Data)
		return
	}

	data, err := json.MarshalIndent(resp.Data, "", "\t")
	fmt.Println(string(data))
}

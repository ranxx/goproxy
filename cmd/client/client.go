package client

import (
	"log"

	"github.com/ranxx/goproxy/service"
	transfer "github.com/ranxx/goproxy/transfer/client"
	"github.com/ranxx/goproxy/utils"
	"github.com/spf13/cobra"
)

// ip port
var ip *string
var port *int

// Command ...
func Command() *cobra.Command {
	client := &cobra.Command{
		Use:   "client",
		Short: "client. Network Address Translation，NAT",
		Long: `内网穿透客户端，在被穿透的机器上启动
eg:
1) goproxy client
2) goproxy client --ip 127.0.0.1
2) goproxy client --port 12341
2) goproxy client --ip 127.0.0.1 --port 12341`,
		Run: func(cmd *cobra.Command, args []string) {
			srv := service.NewClient(*ip, *port)
			go srv.Start()
			utils.IgnoreSignal(func() {
				srv.Close()
				transfer.Manage.Close()
				log.Println("client", "退出")
			})
		},
	}
	ip = client.Flags().StringP("ip", "i", "", `ip (default "")`)
	port = client.Flags().IntP("port", "p", 12341, "port")
	return client
}

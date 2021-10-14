package client

import (
	"log"

	"github.com/ranxx/goproxy/config"
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

注意：默认ip和port为配置文件的server.ip，server.port配置

1) goproxy client
2) goproxy client -i 127.0.0.1
2) goproxy client -p 12341
2) goproxy client -i 127.0.0.1 -p 12341`,
		Run: func(cmd *cobra.Command, args []string) {
			// 如果命令行没有设置，则使用config中的配置
			if !cmd.Flag("ip").Changed {
				*ip = config.Cfg.Server.IP
			}

			if !cmd.Flag("port").Changed {
				*port = config.Cfg.Server.Port
			}

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

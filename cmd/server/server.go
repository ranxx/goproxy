package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ranxx/goproxy/api"
	"github.com/ranxx/goproxy/service"
	transfer "github.com/ranxx/goproxy/transfer/service"
	"github.com/ranxx/goproxy/utils"
	"github.com/spf13/cobra"
)

var port *int
var apiPort *int

// Command ...
func Command() *cobra.Command {
	server := &cobra.Command{
		Use:   "server",
		Short: "server. Network Address Translation，NAT",
		Long: `内网穿透服务端，在公网服务机器上启动
eg:
1) goproxy server
2) goproxy server --port 12341
3) goproxy server --api-port 12351
4) goproxy server --port 12341 --api-port 12351`,
		Run: func(cmd *cobra.Command, args []string) {
			srv := service.NewService("", *port)
			go srv.Start()

			// 自动开启 api
			go func() {
				time.Sleep(time.Second * 2)
				e := api.Init()
				log.Println("starting api", "port", fmt.Sprintf(":%d", *apiPort))
				log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *apiPort), e))
			}()

			utils.IgnoreSignal(func() {
				srv.Close()
				transfer.Manage.Close()
				log.Println("service", "退出")
			})
		},
	}
	port = server.Flags().IntP("port", "p", 12341, "port")
	apiPort = server.Flags().IntP("api-port", "a", 12351, "api port")
	server.AddCommand(new(addCommand).Cmd(), new(removeCommand).Cmd(), new(listCommand).Cmd())
	return server
}

package main

import (
	"fmt"

	"github.com/ranxx/goproxy/cmd/client"
	"github.com/ranxx/goproxy/cmd/server"
	"github.com/ranxx/goproxy/config"
	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{
		Use:   "goproxy",
		Short: "proxy in golang",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd, args)
		},
	}
	root.PersistentFlags().String("config", "./config.yaml", "config file")
	root.AddCommand(client.Command(), server.Command())
	root.Execute()
}

func initConfig(cmd *cobra.Command, args []string) error {
	if cmd.Flag("config") == nil {
		return fmt.Errorf("找不到配置文件")
	}
	file := cmd.Flag("config").Value.String()
	config.ParseFile(file)
	fmt.Println(config.Cfg)
	return nil
}

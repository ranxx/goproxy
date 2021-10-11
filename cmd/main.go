package main

import (
	"github.com/ranxx/goproxy/cmd/client"
	"github.com/ranxx/goproxy/cmd/server"
	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{
		Use:   "goproxy",
		Short: "proxy in golang",
	}
	root.AddCommand(client.Command(), server.Command())
	root.Execute()
}

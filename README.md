# goproxy

本项目使用 golang 编写的内网穿透工具，分为 server 和 client，目前需要在先启动 [server](server/server.go)，然后在需要被穿透的设备启动 [client](client/client.go)。
# goproxy

## 启动

### Linux

如果安装了docker，server可直接敲 `make` 命令启动，否则请移步`教程`

## 教程

本项目使用 golang 编写的内网穿透工具，分为 server 和 client，

启动server `go run cmd/main.go server`
```bash
➜  proxy-test git:(main) ✗ go run cmd/main.go server
2021/10/11 15:23:58 start service :12341
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /transfer                 --> github.com/ranxx/goproxy/api.(*Transfer).List-fm (1 handlers)
[GIN-debug] POST   /transfer/tcp             --> github.com/ranxx/goproxy/api.(*Transfer).NewTCP-fm (1 handlers)
[GIN-debug] DELETE /transfer/tcp/:id         --> github.com/ranxx/goproxy/api.(*Transfer).CloseTCP-fm (1 handlers)
[GIN-debug] DELETE /transfer/port            --> github.com/ranxx/goproxy/api.(*Transfer).RemovePorts-fm (1 handlers)
2021/10/11 15:24:00 starting api port :12351
```

在被穿透的机器上启动 `go run cmd/main.go client --ip ServerIp`
```bash
➜  proxy-test git:(main) ✗ go run cmd/main.go client                    
2021/10/11 15:24:58 client 连上server 127.0.0.1:61397 -> 127.0.0.1:12341
```
此时server端会打印
```bash
2021/10/11 15:24:58 service 新连接 127.0.0.1:61397
```

此时穿透服务已启动，开始配置穿透的端口

启动穿透端口，执行 `go run cmd/main.go server add -l :访问的port -r :被穿透的port` (目前需要在server机器上运行此命令)
```bash
➜  proxy-test git:(main) ✗ go run cmd/main.go server add -l :2222 -r :22
200 ok <nil>
➜  proxy-test git:(main) ✗ 
```
此时server端打印，表明可以通过2222端口访问22端口
```
2021/10/11 15:26:24 transfer.tcp :2222 -> :22 runing
```

执行 `ssh axing@localhost -p 2222` 命令，按提示输入密码，就会成功登录
```bash
➜  proxy-test git:(main) ✗ ssh axing@localhost -p 2222
Password:
Last login: Mon Oct 11 15:17:29 2021 from 127.0.0.1
➜  ~
```
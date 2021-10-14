.PHONY: server client .IGNORE cleanServer cleanClient build runs runc

appName=goproxy
dockerContainerNameServer=goproxys
dockerContainerNameClient=goproxyc

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0


tag=`git branch | grep \* | cut -d ' ' -f2`

server: cleanServer .IGNORE build runs

build:
	# 先编译
	go build -o ${appName} cmd/main.go
	# 编译镜像
	docker build -t ${appName}:${tag} -f ./Dockerfile .

cleanServer:
	@-docker rm -f ${dockerContainerNameServer}

runs:
	docker run --rm -it -d --name ${dockerContainerNameServer} --network host ${appName}:${tag} /cmd/${appName} server

.IGNORE:
	@docker rmi -f ${appName}:${tag}

client: cleanClient .IGNORE build runc

cleanClient:
	@-docker rm -f ${dockerContainerNameClient}

runc:
	docker run --rm -it -d --name ${dockerContainerNameClient} --network host ${appName}:${tag} /cmd/${appName} client
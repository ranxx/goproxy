.PHONY: client

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

tag=`git branch | grep \* | cut -d ' ' -f2`

docker client:
	echo $0
	# 先编译
	go build -o client client/client.go
	# 编译镜像
	docker build -t client:${tag} -f ./Dockerfile_client .
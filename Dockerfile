FROM ubuntu:latest

WORKDIR /cmd

ADD ./goproxy /cmd

CMD ["/cmd/goproxy"]
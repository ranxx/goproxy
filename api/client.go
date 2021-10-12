package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranxx/goproxy/service"
)

// Client 客户端
type Client struct{}

func initClientRouter(e *gin.Engine) {
	cli := new(Client)
	g := e.Group("/client")
	g.GET("", cli.List)
}

type clientListResp struct {
	Items []string `json:"items"`
}

// List 列表
func (cli *Client) List(ctx *gin.Context) {
	resp := &clientListResp{Items: make([]string, 0, 512)}
	service.Svc.Clients.Range(func(key, value interface{}) bool {
		name, ok := key.(string)
		if !ok {
			return ok
		}
		resp.Items = append(resp.Items, name)
		return true
	})
	ctx.JSON(http.StatusOK, Message{
		Code: http.StatusOK,
		Msg:  "ok",
		Data: resp,
	})
}

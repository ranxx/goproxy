package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranxx/goproxy/proto"
	transfer "github.com/ranxx/goproxy/transfer/service"
)

// Transfer ...
type Transfer struct{}

func initTransferRouter(e *gin.Engine) {
	t := new(Transfer)
	g := e.Group("/transfer")
	g.GET("", t.List)

	// tcp
	g.POST("/tcp", t.NewTCP)
	g.DELETE("/tcp/:id", t.CloseTCP)

	// http
	g.POST("/http", t.NewHTTP)
	g.DELETE("/http/:id", t.CloseHTTP)
}

type listResp struct {
	ID      int64
	Network proto.NetworkType
	Laddr   proto.Addr `json:"laddr"`
	Raddr   proto.Addr `json:"raddr"`
}

// List 列表
func (t *Transfer) List(ctx *gin.Context) {
	res := make([]*listResp, 0, 512)
	transfer.Manage.Range(func(v transfer.Transfer) {
		id, laddr, raddr := v.Info()
		tmp := &listResp{
			ID:      id,
			Network: v.NetWork(),
			Laddr:   laddr,
			Raddr:   raddr,
		}
		res = append(res, tmp)
	})
	ctx.JSON(http.StatusOK, res)
}

type newTCPReq struct {
	Laddr proto.Addr `json:"laddr"`
	Raddr proto.Addr `json:"raddr"`
}

// NewTCP 新建tcp
func (t *Transfer) NewTCP(ctx *gin.Context) {
	req := new(newTCPReq)
	if err := ctx.BindJSON(req); err != nil {
		ctx.Writer.Write([]byte(fmt.Sprintf("%s", err)))
		ctx.Abort()
		return
	}

	err := transfer.NewTransferWithIPPort(req.Laddr.Ip, int(req.Laddr.Port), req.Raddr.Ip, int(req.Raddr.Port), proto.NetworkType_TCP).Start()
	if err != nil {
		ctx.Writer.Write([]byte(fmt.Sprintf("%s", err)))
		ctx.Abort()
		return
	}

	ctx.String(http.StatusOK, "ok")
}

// CloseTCP 关闭tcp
func (t *Transfer) CloseTCP(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		ctx.Writer.Write([]byte(fmt.Sprintf("%s", "没有必要参数id")))
		ctx.Abort()
		return
	}

	_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.Writer.Write([]byte(fmt.Sprintf("%s", "id不是一个整数类型")))
		ctx.Abort()
		return
	}

	transfer.Manage.Remove(_id)
	ctx.String(http.StatusOK, "ok")
}

// NewHTTP ...
func (t *Transfer) NewHTTP(ctx *gin.Context) {
	req := new(newTCPReq)
	if err := ctx.BindJSON(req); err != nil {
		ctx.Writer.Write([]byte(fmt.Sprintf("%s", err)))
		ctx.Abort()
		return
	}

	err := transfer.NewTransferWithIPPort(req.Laddr.Ip, int(req.Laddr.Port), req.Raddr.Ip, int(req.Raddr.Port), proto.NetworkType_HTTP).Start()
	if err != nil {
		ctx.Writer.Write([]byte(fmt.Sprintf("%s", err)))
		ctx.Abort()
		return
	}
	ctx.String(http.StatusOK, "ok")
}

// CloseHTTP ...
func (t *Transfer) CloseHTTP(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		ctx.Writer.Write([]byte(fmt.Sprintf("%s", "没有必要参数id")))
		ctx.Abort()
		return
	}

	_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.Writer.Write([]byte(fmt.Sprintf("%s", "id不是一个整数类型")))
		ctx.Abort()
		return
	}

	transfer.Manage.Remove(_id)
	ctx.String(http.StatusOK, "ok")
}

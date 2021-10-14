package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranxx/goproxy/config"
	"github.com/ranxx/goproxy/proto"
	transfer "github.com/ranxx/goproxy/transfer/service"
)

// Message ...
type Message struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func initTransferRouter(e *gin.Engine) {
	t := new(Transfer)
	g := e.Group("/transfer")
	g.GET("", t.List)

	// tcp
	g.POST("/tcp", t.NewTCP)
	g.DELETE("/tcp/:id", t.CloseTCP)
	g.DELETE("/port", t.RemovePorts)
}

type listResp struct {
	ID      int64             `json:"id"`
	Network proto.NetworkType `json:"network"`
	Laddr   proto.Addr        `json:"laddr"`
	Raddr   proto.Addr        `json:"raddr"`
	Machine string            `json:"machine"`
}

// Transfer ...
type Transfer struct{}

// List 列表
func (t *Transfer) List(ctx *gin.Context) {
	res := make([]*listResp, 0, 512)
	transfer.Manage.Range(func(v transfer.Transfer) {
		info := v.Info()
		tmp := &listResp{
			ID:      info.Index,
			Network: v.NetWork(),
			Laddr:   info.Laddr,
			Raddr:   info.Raddr,
			Machine: info.Machine,
		}
		res = append(res, tmp)
	})

	ctx.JSON(http.StatusOK, Message{
		Code: http.StatusOK,
		Msg:  "ok",
		Data: res,
	})
}

type newTCPReq struct {
	Laddr   []proto.Addr `json:"laddr"`
	Raddr   proto.Addr   `json:"raddr"`
	Machine string       `json:"machine"`
}

// NewTCP 新建tcp
func (t *Transfer) NewTCP(ctx *gin.Context) {
	req := new(newTCPReq)
	if err := ctx.BindJSON(req); err != nil {
		ctx.JSON(http.StatusOK, Message{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}

	for _, laddr := range req.Laddr {
		if laddr.Ip == "" {
			laddr.Ip = config.Cfg.Server.IP
		}
		err := transfer.NewTransferWithIPPort(req.Machine, laddr.Ip, int(laddr.Port), req.Raddr.Ip, int(req.Raddr.Port), proto.NetworkType_TCP).Start()
		if err != nil {
			ctx.JSON(http.StatusOK, Message{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			})
			ctx.Abort()
			return
		}
	}

	ctx.JSON(http.StatusOK, Message{
		Code: http.StatusOK,
		Msg:  "ok",
	})
}

// CloseTCP 关闭tcp
func (t *Transfer) CloseTCP(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		ctx.JSON(http.StatusOK, Message{
			Code: http.StatusBadRequest,
			Msg:  fmt.Sprintf("%s", "没有必要参数id"),
		})
		ctx.Abort()
		return
	}

	_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Message{
			Code: http.StatusBadRequest,
			Msg:  fmt.Sprintf("%s", "id不是一个整数类型"),
		})
		ctx.Abort()
		return
	}

	transfer.Manage.Remove(_id)
	ctx.JSON(http.StatusOK, Message{
		Code: http.StatusOK,
		Msg:  "ok",
	})
}

type removePorts struct {
	Ports []int `json:"ports"`
}

// RemovePorts 按端口移除
func (t *Transfer) RemovePorts(ctx *gin.Context) {
	req := new(removePorts)
	if err := ctx.BindJSON(req); err != nil {
		ctx.JSON(http.StatusOK, Message{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}

	transfer.Manage.RemoveByPort(req.Ports...)
	ctx.JSON(http.StatusOK, Message{
		Code: http.StatusOK,
		Msg:  "ok",
	})
}

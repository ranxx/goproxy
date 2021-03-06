package api

import "github.com/gin-gonic/gin"

// Init ...
func Init() *gin.Engine {
	e := gin.New()

	initTransferRouter(e)
	initClientRouter(e)

	return e
}

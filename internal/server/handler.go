package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
	return
}

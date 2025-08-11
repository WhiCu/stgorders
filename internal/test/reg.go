package test

import (
	"github.com/WhiCu/stgorders/internal/test/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	handle := NewHandler()

	r.GET("/ping", handle.Ping)

}

func NewHandler() *handler.Handler {
	return handler.NewHandler(nil)
}

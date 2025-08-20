package handler

import "github.com/gin-gonic/gin"

func (h *Handler) Ping(c *gin.Context) {
	c.String(200, "pong")
}

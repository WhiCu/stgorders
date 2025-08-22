package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

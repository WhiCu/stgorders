package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Inter(c *gin.Context) {
	h.log.With(slog.String("IP", c.ClientIP())).Debug("GET")
	c.File("internal/web-interface/templates/index.html")
}

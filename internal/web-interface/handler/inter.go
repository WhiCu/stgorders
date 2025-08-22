package handler

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templates embed.FS

func (h *Handler) Inter(c *gin.Context) {
	h.log.With(slog.String("IP", c.ClientIP())).Debug("GET")
	f, err := templates.Open("templates/index.html")
	if err != nil {
		h.log.Error("cannot load index.html", slog.String("ERR", err.Error()))
		c.String(http.StatusInternalServerError, "cannot load index.html")
		return
	}

	stat, err := f.Stat()
	if err != nil {
		h.log.Error("cannot get stat of index.html", slog.String("ERR", err.Error()))
		c.String(http.StatusInternalServerError, "cannot get stat of index.html")
		return
	}

	c.DataFromReader(http.StatusOK, stat.Size(), "text/html; charset=utf-8", f, nil)
}

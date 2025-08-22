package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID string `uri:"orderUID" binding:"required"`
}

func (h *Handler) Order(c *gin.Context) {
	o := new(Order)
	if err := c.ShouldBindUri(o); err != nil {
		h.log.Debug("invalid order id", slog.String("ERR", err.Error()))
		c.String(http.StatusBadRequest, "invalid order id")
		return
	}
	log := h.log.With(slog.String("order_id", o.ID), slog.String("IP", c.ClientIP()))

	order, err := h.service.GetJsonOrderByUID(c, o.ID)
	if err != nil {
		log.Debug("could not get order", slog.String("ERR", err.Error()))
		c.String(http.StatusNotFound, "order not found")
		return
	}
	log.Info("order found")
	c.JSON(http.StatusOK, order)
}

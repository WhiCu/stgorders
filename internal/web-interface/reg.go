package wi

import (
	"log/slog"

	"github.com/WhiCu/stgorders/db/cache"
	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/db/storage"
	"github.com/WhiCu/stgorders/internal/web-interface/client"
	"github.com/WhiCu/stgorders/internal/web-interface/handler"
	"github.com/WhiCu/stgorders/internal/web-interface/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, log *slog.Logger, storage *storage.Storage, cache cache.Cache[string, model.JsonOrder]) {
	stg := client.NewStorage(storage, cache, log.WithGroup("storageAdapter"))
	srv := service.NewService(stg, log.WithGroup("service"))
	h := handler.NewHandler(srv, log.WithGroup("handler"))
	r.GET(":orderUID", h.Order)
	r.GET("", h.Inter)
}

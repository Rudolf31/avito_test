package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewGin(lc fx.Lifecycle) *gin.Engine {
	router := gin.Default()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				servErr := server.ListenAndServe()
				if servErr != nil {
					panic(servErr.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		},
	})

	return router
}

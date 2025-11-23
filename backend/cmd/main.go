package main

import (
	"avito_test/backend/internal/db"
	routes "avito_test/backend/internal/routers"
	"avito_test/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func main() {

	fx.New(
		fx.Provide(
			db.NewPool,
			routes.NewGin,
		),
		services.Module,

		fx.Invoke(func(pool *pgxpool.Pool, router *gin.Engine, teamService services.TeamService, userService services.UserService, prService services.PullRequestService) {
			routes.RegisterRoutes(router, teamService, userService, prService)
		}),
	).Run()
}

package router

import (
	"github.com/Rahmans11/final-phase-3/internal/controller"
	"github.com/Rahmans11/final-phase-3/internal/middleware"
	"github.com/Rahmans11/final-phase-3/internal/repository"
	"github.com/Rahmans11/final-phase-3/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func RegisterFollowsRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	followsRouter := app.Group("/follows")

	followsRepository := repository.NewFollowsRepository()
	followsService := service.NewFollowsService(followsRepository, db, rdb)
	followsController := controller.NewFollowsController(followsService)

	followsRouter.Use(middleware.VerifyJWT, middleware.RecognizedOnly(rdb, "USER"))
	{
		followsRouter.POST("/add-followed", followsController.AddFollowed)
	}
}

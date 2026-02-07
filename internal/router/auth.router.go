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

func RegisterAuthRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := app.Group("/auth")

	authRepository := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepository, db, rdb)
	authController := controller.NewAuthController(authService)

	authRouter.POST("/", authController.Login)
	authRouter.POST("/register", authController.Register)

	authRouter.Use(middleware.VerifyJWT, middleware.RecognizedOnly(rdb, "USER"))
	{
		authRouter.DELETE("/logout", authController.Logout)
	}
}

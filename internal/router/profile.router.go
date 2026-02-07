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

func RegisterProfilesRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {

	profilesRouter := app.Group("/profile")

	profileRepository := repository.NewProfileRepository()
	profileService := service.NewProfileService(profileRepository, db, rdb)
	profileController := controller.NewProfileController(profileService)

	profilesRouter.Use(middleware.VerifyJWT, middleware.RecognizedOnly(rdb, "USER"))
	{
		profilesRouter.GET("/", profileController.GetProfile)
		profilesRouter.GET("/user/:id", profileController.GetOtherProfile)
		profilesRouter.PATCH("/", profileController.EditProfile)
		// profilesRouter.PUT("/change-password", profileController.ChangePassword)
	}

}

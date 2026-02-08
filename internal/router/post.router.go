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

func RegisterPostsRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {

	postsRouter := app.Group("/post")

	postRepository := repository.NewPostsRepository()
	postService := service.NewPostsService(postRepository, db, rdb)
	postController := controller.NewPostController(postService)

	postsRouter.Use(middleware.VerifyJWT, middleware.RecognizedOnly(rdb, "USER"))
	{
		postsRouter.POST("/", postController.CreatePost)
	}

}

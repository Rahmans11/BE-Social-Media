package router

import (
	"net/http"

	_ "github.com/Rahmans11/final-phase-3/docs"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Static("/static/img", "public")
	app.Static("/static/pages", "public/html")

	RegisterAuthRouter(app, db, rdb)
	RegisterProfilesRouter(app, db, rdb)

	app.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/static/pages/not-found.html")
	})
}

package main

import (
	"log"

	"github.com/Rahmans11/final-phase-3/internal/config"
	"github.com/Rahmans11/final-phase-3/internal/middleware"
	"github.com/Rahmans11/final-phase-3/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           Final phase
// @version         1.0
// @description     Social Media service for Final phase
// @host      		localhost:8080
// @BasePath  		/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env")
		return
	}

	db, err := config.InitDB()
	if err != nil {
		log.Println("Failed to connect database")
		return
	}

	defer db.Close()

	rdb := config.InitRedis()

	defer rdb.Close()

	app := gin.Default()

	app.Use(middleware.CORSMiddleware)

	router.Init(app, db, rdb)

	app.Run("localhost:8080")
}

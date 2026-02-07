package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	user := os.Getenv("RDS_USER")
	pass := os.Getenv("RDS_PASS")
	host := os.Getenv("RDS_HOST")
	port := os.Getenv("RDS_PORT")
	dtbs := os.Getenv("RDS_DTBS")

	db, _ := strconv.Atoi(dtbs)

	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: user,
		Password: pass,
		DB:       db,
	})
}

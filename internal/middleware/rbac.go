package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/Rahmans11/final-phase-3/internal/dto"
	"github.com/Rahmans11/final-phase-3/pkg"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RecognizedOnly(rdb *redis.Client, role ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, isExist := c.Get("token")
		if !isExist {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Response{
				Msg:     "Forbidden Access",
				Success: false,
				Data:    []any{},
				Error:   "Access Denied",
			})
			return
		}

		accessToken, ok := token.(pkg.JWTClaims)

		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Msg:     "Internal Server Error",
				Success: false,
				Data:    []any{},
				Error:   "internal server error",
			})
			return
		}

		isAllowed := slices.Contains(role, accessToken.Role)

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Response{
				Msg:     "Forbidden Access",
				Success: false,
				Data:    []any{},
				Error:   "Access Denied",
			})
			return
		}

		//var rdb redis.Client
		rkey := fmt.Sprintf("rahman:social-media:whitelist-token:%d", accessToken.Id)
		_, err := rdb.Get(c, rkey).Result()

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Msg:     "Unauthorized Access",
				Success: false,
				Data:    []any{},
				Error:   "Access Denied",
			})
			return
		}

		// isInWhiteList :=

		c.Next()
	}
}

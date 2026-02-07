package middleware

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(c *gin.Context) {
	defer c.Next()

	// Allowed Origin
	whiteListOrigin := []string{"http://localhost:5500", "http://localhost:5173", "http://localhost:8080"}
	origin := c.GetHeader("Origin")
	if slices.Contains(whiteListOrigin, origin) {
		c.Header("Access-Control-Allow-Origin", origin)
	} else {
		log.Printf("Origin is not in the Whitelist: %s", origin)
	}
	// Allowed Header
	allowedHeaders := []string{"x-who-am-i", "Content-Type", "Authorization", "Origin"}
	c.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
	// Allowed Methods
	allowedMethod := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodHead}
	c.Header("Access-Control-Allow-Methods", strings.Join(allowedMethod, ", "))
	// Handling preflight
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
}

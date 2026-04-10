package middleware

import (
	"strconv"
	"strings"

	"kbmanage/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func CORS(cfg repository.Config) gin.HandlerFunc {
	allowOriginSet := make(map[string]struct{}, len(cfg.CORSAllowOrigins))
	for _, o := range cfg.CORSAllowOrigins {
		allowOriginSet[o] = struct{}{}
	}

	allowMethods := strings.Join(cfg.CORSAllowMethods, ", ")
	allowHeaders := strings.Join(cfg.CORSAllowHeaders, ", ")
	exposeHeaders := strings.Join(cfg.CORSExposeHeaders, ", ")
	maxAge := "600"
	if cfg.CORSMaxAgeSeconds > 0 {
		maxAge = strconv.Itoa(cfg.CORSMaxAgeSeconds)
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowOriginSet[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
				c.Header("Access-Control-Allow-Methods", allowMethods)
				c.Header("Access-Control-Allow-Headers", allowHeaders)
				if exposeHeaders != "" {
					c.Header("Access-Control-Expose-Headers", exposeHeaders)
				}
				if cfg.CORSAllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
				c.Header("Access-Control-Max-Age", maxAge)
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

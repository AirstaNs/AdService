package httpgin

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now().UTC()

		c.Next()
		latencyTime := time.Since(startTime)

		logger.Printf("[%s] | %s  | %s  | %s | %d | %s  | %s\n",
			time.Now().UTC().Format(time.DateTime),
			c.ClientIP(),
			c.Request.Method,
			c.Request.Proto,
			c.Writer.Status(),
			latencyTime.String(),
			c.Request.URL.Path,
		)
	}
}

func RecoveryMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("PANIC ERROR: %v\n", r)

				c.JSON(http.StatusInternalServerError, ErrorResponse(errors.New("internal server error")))
			}
		}()
		c.Next()
	}
}

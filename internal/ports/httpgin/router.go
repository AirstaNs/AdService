package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/app"
	"log"
	"net/http"
	"net/http/pprof"
)

func AppRouter(r *gin.RouterGroup, a app.App, logger *log.Logger) {
	r.Use(LoggerMiddleware(logger))
	r.Use(RecoveryMiddleware(logger))

	r.GET("/ads/:ad_id", getAdByID(a))
	r.GET("/ads", getAdsByFilter(a))
	r.POST("/ads", createAd(a))
	r.PUT("/ads/:ad_id/status", changeAdStatus(a))
	r.PUT("/ads/:ad_id", updateAd(a))
	r.DELETE("/ads/:ad_id", deleteAd(a))

	r.GET("/users/:user_id", getUserByID(a))
	r.POST("/users", createUser(a))
	r.PUT("/users/:user_id", updateUser(a))
	r.DELETE("/users/:user_id", deleteUser(a))
	// регистрируем маршруты для обработки запросов pprof
	r.GET("/debug/pprof/", gin.WrapH(http.HandlerFunc(pprof.Index)))
	r.GET("/debug/pprof/cmdline", gin.WrapH(http.HandlerFunc(pprof.Cmdline)))
	r.GET("/debug/pprof/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
	r.GET("/debug/pprof/symbol", gin.WrapH(http.HandlerFunc(pprof.Symbol)))
	r.GET("/debug/pprof/trace", gin.WrapH(http.HandlerFunc(pprof.Trace)))
}

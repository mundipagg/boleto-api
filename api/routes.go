package api

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/middleware"
)

func Base(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "./boleto/favicon.ico")
	router.GET("/boleto", getBoleto)
	router.GET("/boleto/memory-check/:unit", memory)
	router.GET("/boleto/memory-check/", memory)
	router.GET("/boleto/confirmation", confirmation)
	router.POST("/boleto/confirmation", confirmation)
}

//V1 configura as rotas da v1
func V1(router *gin.Engine) {
	v1 := router.Group("v1")
	v1.Use(timingMetrics())
	v1.Use(returnHeaders())
	v1.POST("/boleto/register", middleware.Authentication, middleware.ParseBoleto, middleware.ValidateRegisterV1, middleware.Logger, registerBoleto)
	v1.GET("/boleto/:id", getBoletoByID)
}

//V2 configura as rotas da v2
func V2(router *gin.Engine) {
	v2 := router.Group("v2")
	v2.Use(timingMetrics())
	v2.Use(returnHeaders())
	v2.POST("/boleto/register", middleware.Authentication, middleware.ParseBoleto, middleware.ValidateRegisterV2, middleware.Logger, registerBoleto)
}

func returnHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func executionController() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.IsRunning() {
			c.AbortWithError(500, errors.New("a aplicação está sendo finalizada"))
			return
		}
	}
}

func timingMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		total := end.Sub(start)
		s := float64(total.Seconds())
		metrics.PushTimingMetric("request-time", s)
	}
}

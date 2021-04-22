package api

import "github.com/gin-gonic/gin"

//V1 instala a api versao 1
func V1(router *gin.Engine) {
	v1 := router.Group("v1")
	v1.Use(timingMetrics())
	v1.Use(returnHeaders())
	v1.POST("/boleto/register", Authentication, ParseBoleto, Logger, registerBoleto)
	v1.GET("/boleto/:id", getBoletoByID)
}

func Base(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "./boleto/favicon.ico")
	router.GET("/boleto/memory-check/:unit", memory)
	router.GET("/boleto/memory-check/", memory)
	router.GET("/boleto", getBoleto)
	router.GET("/boleto/confirmation", confirmation)
	router.POST("/boleto/confirmation", confirmation)
}


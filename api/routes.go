package api

import (
	"github.com/gin-gonic/gin"
)

func Base(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "./boleto/favicon.ico")
	router.GET("/boleto", getBoleto)
	router.GET("/boleto/memory-check/:unit", memory)
	router.GET("/boleto/memory-check/", memory)
	router.GET("/boleto/confirmation", confirmation)
	router.POST("/boleto/confirmation", confirmation)
}

//V1 instala a api versao 1
func V1(router *gin.Engine) {
	v1 := router.Group("v1")
	v1.Use(timingMetrics())
	v1.Use(returnHeaders())
	v1.POST("/boleto/register", Authentication, ParseBoleto, ValidateRegisterV1, Logger, registerBoleto)
	v1.GET("/boleto/:id", getBoletoByID)
}

//V2 intala a api versao 2
func V2(router *gin.Engine) {
	v2 := router.Group("v2")
	v2.Use(timingMetrics())
	v2.Use(returnHeaders())
	v2.POST("/boleto/register", Authentication, ParseBoleto, ValidateRegisterV2, Logger, registerBoleto)
}

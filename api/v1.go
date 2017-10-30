package api

import "github.com/gin-gonic/gin"

//InstallV1 instala a api versao 1
func InstallV1(router *gin.Engine) {
	v1 := router.Group("v1")
	v1.Use(timingMetrics())
	v1.Use(ReturnHeaders())
	v1.Use(ParseBoleto())
	v1.POST("/boleto/register", registerBoleto)
	v1.GET("/boleto/:id", getBoletoByID)
}

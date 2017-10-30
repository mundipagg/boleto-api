package mock

import "github.com/gin-gonic/gin"

//Run sobe uma aplicação web para mockar a integração com os Bancos
func Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/oauth/token", authBB)
	router.POST("/registrarBoleto", registerBoletoBB)
	router.POST("/caixa/registrarBoleto", registerBoletoCaixa)
	router.POST("/citi/registrarBoleto", registerBoletoCiti)
	router.POST("/santander/get-ticket", getTicket)
	router.POST("/santander/register", registerBoletoSantander)
	router.POST("/bradesco/registrarBoleto", registerBoletoBradesco)

	router.Run(":9091")
}

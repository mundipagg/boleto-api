package mock

import "github.com/gin-gonic/gin"
import "github.com/mundipagg/boleto-api/env"

//Run sobe uma aplicação web para mockar a integração com os Bancos
func Run(port string) {
	env.ConfigMock(port)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/oauth/token", authBB)
	router.POST("/registrarBoleto", registerBoletoBB)
	router.POST("/caixa/registrarBoleto", registerBoletoCaixa)
	router.POST("/citi/registrarBoleto", registerBoletoCiti)
	router.POST("/santander/get-ticket", getTicket)
	router.POST("/santander/register", registerBoletoSantander)
	router.POST("/bradescoshopfacil/registrarBoleto", registerBoletoBradescoShopFacil)
	router.POST("/itau/gerarToken", getTokenItau)
	router.POST("/itau/registrarBoleto", registerItau)
	router.POST("/bradesconetempresa/registrarBoleto", registerBoletoBradescoNetEmpresa)
	router.POST("/pfisa/gerarToken", getTokenPfisa)
	router.POST("/pfisa/registrarBoleto", registerPfisa)
	router.Run(":" + port)
}

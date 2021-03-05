package mock

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/env"
)

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
	router.POST("/pefisa/gerarToken", getTokenPefisa)
	router.POST("/pefisa/registrarBoleto", registerPefisa)
	router.Run(":" + port)
}

//StartMockService Inicializa servidor de mock para testes unitários
func StartMockService(port string) {
	env.Config(true, true, true)
	go Run(port)
	time.Sleep(2 * time.Second)
}

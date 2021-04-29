package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
)

//InstallRestAPI "instala" e sobe o servico de rest
func InstallRestAPI() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(executionController())
	if config.Get().DevMode && !config.Get().MockMode {
		router.Use(gin.Logger())
	}

	Base(router)
	V1(router)
	V2(router)

	router.Run(config.Get().APIPort)
}

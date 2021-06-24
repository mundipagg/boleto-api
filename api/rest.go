package api

import (
	"net/http/httputil"

	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/mundipagg/boleto-api/metrics"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
)

//InstallRestAPI "instala" e sobe o servico de rest
func InstallRestAPI() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(executionController())
	useNewRelic(router)

	if config.Get().DevMode && !config.Get().MockMode {
		router.Use(gin.Logger())
	}

	Base(router)
	V1(router)
	V2(router)

	router.Run(config.Get().APIPort)
}

func useNewRelic(router *gin.Engine) {
	if !config.Get().TelemetryEnabled {
		return
	}

	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName(config.Get().NewRelicAppName),
		newrelic.ConfigLicense(config.Get().NewRelicLicence),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	router.Use(nrgin.Middleware(app))
}

func memory(c *gin.Context) {
	unit := c.Param("unit")
	c.JSON(200, metrics.GetMemoryReport(unit))
}

func confirmation(c *gin.Context) {
	if dump, err := httputil.DumpRequest(c.Request, true); err == nil {
		l := log.CreateLog()
		l.BankName = "BradescoShopFacil"
		l.Operation = "BoletoConfirmation"
		l.Request(string(dump), c.Request.URL.String(), nil)
	}
	c.String(200, "OK")
}

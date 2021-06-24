package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

//Memory HealthCheck com relatório de memória
func memory(c *gin.Context) {
	unit := c.Param("unit")
	c.JSON(200, metrics.GetMemoryReport(unit))
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

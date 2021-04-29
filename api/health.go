package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/metrics"
)

//Memory HealthCheck com relatório de memória
func memory(c *gin.Context) {
	unit := c.Param("unit")
	c.JSON(200, metrics.GetMemoryReport(unit))
}

package metrics

import (
	. "github.com/PMoneda/telemetry"
	"github.com/PMoneda/telemetry/registry"
	"github.com/mundipagg/boleto-api/config"
)

var business *Telemetry

//InstallBusinessMetrics Instância a telemetria de negócio
func InstallBusinessMetrics(cnf registry.Config) {
	value := Database("boleto-api").RetentionPolicy("business").Measurement("boletos").Tag("host").Value("host0")
	business = BuildTelemetryContext(cnf, Context(value))
	go business.StartTelemetry(true)
}

//GetBusinessMetrics Obtém um objeto de telemetria do negócio
func GetBusinessMetrics() *Telemetry {
	return business
}

//PushBusinessMetric Envio dados de negócio para a telemetria
func PushBusinessMetric(tag string, value interface{}) {
	if config.Get().EnableMetrics {
		GetBusinessMetrics().Push(tag, value)
	}
}

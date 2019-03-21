package metrics

import (
	. "github.com/PMoneda/telemetry"
	"github.com/PMoneda/telemetry/registry"
	"github.com/mundipagg/boleto-api/config"
)

var timing *Telemetry

//InstallTimingMetrics Instância a telemetria de tempo de resposta
func InstallTimingMetrics(cnf registry.Config) {
	value := Database("boleto-api").RetentionPolicy("runtime").Measurement("timing").Tag("host").Value("host0")
	timing = BuildTelemetryContext(cnf, Context(value))
	go timing.StartTelemetry(true)
}

//GetTimingMetrics Obtém um objeto de telemetria do tempo de resposta
func GetTimingMetrics() *Telemetry {
	return timing
}

//PushTimingMetric Envio dados de tempo de resposta para a telemetria
func PushTimingMetric(tag string, value interface{}) {
	if config.Get().EnableMetrics {
		GetTimingMetrics().Push(tag, value)
	}
}

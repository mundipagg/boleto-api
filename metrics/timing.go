package metrics

import (
	. "github.com/PMoneda/telemetry"
	"github.com/PMoneda/telemetry/registry"
	"github.com/mundipagg/boleto-api/config"
)

var timing *Telemetry

func InstallTimingMetrics(cnf registry.Config) {
	value := Database("boleto-api").RetentionPolicy("runtime").Measurement("timing").Tag("host").Value("host0")
	timing = BuildTelemetryContext(cnf, Context(value))
	go timing.StartTelemetry(true)
}

func GetTimingMetrics() *Telemetry {
	return timing
}

func PushTimingMetric(tag string, value interface{}) {
	if config.Get().EnableMetrics {
		GetTimingMetrics().Push(tag, value)
	}
}

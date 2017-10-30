package metrics

import (
	. "github.com/PMoneda/telemetry"
	"github.com/PMoneda/telemetry/registry"
)

var timing *Telemetry

func InstallTimingMetrics(cnf registry.Config) {
	value := Database("boleto-api").RetentionPolicy("runtime").Measurement("timing").Tag("host").Value("host0")
	timing = BuildTelemetryContext(cnf, Context(value))
	go timing.StartTelemetry()
}

func GetTimingMetrics() *Telemetry {
	return timing
}

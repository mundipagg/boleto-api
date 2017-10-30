package metrics

import "github.com/PMoneda/telemetry/registry"
import . "github.com/PMoneda/telemetry"

var business *Telemetry

func InstallBusinessMetrics(cnf registry.Config) {
	value := Database("boleto-api").RetentionPolicy("business").Measurement("boletos").Tag("host").Value("host0")
	business = BuildTelemetryContext(cnf, Context(value))
	go business.StartTelemetry()
}

func GetBusinessMetrics() *Telemetry {
	return business
}

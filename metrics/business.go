package metrics

import "github.com/PMoneda/telemetry/registry"
import . "github.com/PMoneda/telemetry"

var business *Telemetry

func InstallBusinessMetrics(cnf registry.Config) {
	value := Database("boleto-api").RetentionPolicy("business").Measurement("boletos").Tag("host").Value("host0")
	business = BuildTelemetryContext(cnf, Context(value))
	business.StartTelemetry(false)
}

func GetBusinessMetrics() *Telemetry {
	return business
}

package metrics

import "github.com/PMoneda/telemetry/registry"
import . "github.com/PMoneda/telemetry"

func InstallRuntime(cnf registry.Config) {
	runtime := Database("boleto-api").RetentionPolicy("runtime").Measurement("metrics").Tag("host").Value("host0")
	telemetryRuntime := BuildTelemetryContext(cnf, Context(runtime))
	go telemetryRuntime.StartRuntimeTelemetry()
}

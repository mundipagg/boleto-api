package seq

import (
	"github.com/mralves/tracer"
)

const (
	Debug       = "Debug"
	Information = "Information"
	Warning     = "Warning"
	Error       = "Error"
	Fatal       = "Fatal"
	Verbose     = "Verbose"
)

func Level(level uint8) string {
	switch level {
	case tracer.Debug:
		return Debug
	case tracer.Informational:
		return Information
	case tracer.Notice, tracer.Warning:
		return Warning
	case tracer.Critical, tracer.Error:
		return Error
	case tracer.Alert, tracer.Fatal:
		return Fatal
	default:
		return Verbose
	}
}

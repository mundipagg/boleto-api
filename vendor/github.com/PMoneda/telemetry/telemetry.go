package telemetry

import (
	"runtime"
	"time"

	"github.com/PMoneda/telemetry/registry"
)

type Telemetry struct {
	Registry registry.Registry
	Root     string
}

type Context func(*Telemetry) *Telemetry

type RetentionPolicy Context

type Measurement Context

type Value Context

type Tag Context

func (r Tag) Tag(value string) Value {
	return func(t *Telemetry) *Telemetry {
		r(t)
		t.Root = t.Root + "." + value
		return t
	}
}

func (r Value) Value(value string) Context {
	return func(t *Telemetry) *Telemetry {
		r(t)
		t.Root = t.Root + "." + value
		return t
	}
}

func (r Measurement) Measurement(tag string) Tag {
	return func(t *Telemetry) *Telemetry {
		r(t)
		t.Root = t.Root + "." + tag
		return t
	}
}

func (r RetentionPolicy) RetentionPolicy(s string) Measurement {
	return func(t *Telemetry) *Telemetry {
		r(t)
		t.Root = t.Root + "." + s
		return t
	}
}

func (c Context) Child(s string) Context {
	return func(t *Telemetry) *Telemetry {
		c(t)
		t.Root = t.Root + "." + s
		return t
	}
}

func Database(s string) RetentionPolicy {
	return (func(t *Telemetry) *Telemetry {
		t.Root = s
		return t
	})
}

func NewTelemetryForInfluxDB(config registry.Config, context string) *Telemetry {
	telemetry := new(Telemetry)
	telemetry.Root = context
	telemetry.Registry = registry.NewInfluxClient(config)
	return telemetry
}

func BuildTelemetryContext(config registry.Config, ctx Context) *Telemetry {
	telemetry := new(Telemetry)
	ctx(telemetry)
	telemetry.Registry = registry.NewInfluxClient(config)
	return telemetry
}

func (t *Telemetry) PushAndFlush(tag string, value interface{}) (err error) {
	err = t.Registry.Registry(t.Root+"."+tag, value)
	if err != nil {
		return
	}
	err = t.Registry.Flush(tag)
	return
}

func (t *Telemetry) Push(tag string, value interface{}) (err error) {
	err = t.Registry.Registry(t.Root+"."+tag, value)
	return
}

func (t *Telemetry) Flush() (err error) {
	err = t.Registry.FlushAll()
	return
}

func (t *Telemetry) StartRuntimeTelemetry() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			alloc := float64(m.Alloc) / (1024)
			totalAlloc := float64(m.TotalAlloc) / (1024)
			gcRun := float64(m.NumGC)
			t.Push("memory-alloc", alloc)
			t.Push("total-memory-alloc", totalAlloc)
			t.Push("garbage-collector-total-run", gcRun)
			t.Push("num-goroutines", runtime.NumGoroutine())
			t.Flush()
		}
	}
}

func (t *Telemetry) StartTelemetry() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			t.Flush()

		}
	}
}

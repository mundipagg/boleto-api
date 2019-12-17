# SEQ Writer
This library is a writer(sink) [SEQ](https://datalust.co/) to use with [Tracer](https://github.com/mralves/tracer)

## How to install
Using go get (not recommended):
```bash
go get github.com/mundipagg/tracer-seq-writer
```

Using [dep](github.com/golang/dep) (recommended):
```bash
dep ensure --add github.com/mundipagg/tracer-seq-writer@<version>
```

## Configuration

|Field|Type|Mandatory?|Default|Description|
|---|---|---|:---:|---|
|Address|string|Y||Seq **full** endpoint (i.e. http://seq.io/api/events/eaw)|
|Key|string|N|""|Seq [API Key](https://docs.datalust.co/docs/api-keys)|
|Application|string|Y||Application name|
|Minimum Level|uint8|N|DEBUG|Minimum Level to log following the [syslog](https://en.wikipedia.org/wiki/Syslog#Severity_level) standard|
|Timeout|time.Duration|N|0 (infinite)|Timeout of the HTTP client|
|MessageEnvelop|string|N|"%v"|A envelop that *wraps* the original message|
|DefaultProperties|seq.Entry|N|{}|A generic object to append to *every* log entry, but can be overwritten by the original log entry|
|Buffer.Cap|int|N|100| Maximum capacity of the log buffer, when the buffer is full all logs are sent at once|
|Buffer.OnWait|int|N|100| Maximum size of the queue to send to Seq|
|Buffer.BackOff|time.Duration|N|60 seconds| Delay between retries to Seq|

## How to use

Below follows a simple example of how to use this lib:

```go
package main

import (
	"fmt"
	"time"

	"github.com/mralves/tracer"
	seq "github.com/mundipagg/tracer-seq-writer"

	bsq "github.com/mundipagg/tracer-seq-writer/buffer"
)

type Safe struct {
	tracer.Writer
}

type LogEntry = map[string]interface{}

func (s *Safe) Write(entry tracer.Entry) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%v", err)
		}
	}()
	s.Writer.Write(entry)
}

func configureTracer() {
	var writers []tracer.Writer
	tracer.DefaultContext.OverwriteChildren()
	writers = append(writers, &Safe{seq.New(seq.Config{
		Timeout:      3 * time.Second,
		MinimumLevel: tracer.Debug,
		DefaultProperties: LogEntry{
			"Application": "MyApplicationName",
			"Environment": "Production",
		},
		Key:         "yYka3WMkG5sfqVNlUUYj",
		Address:     "http://localhost:5341/api/events/raw",
		Buffer: bsq.Config{
			OnWait:     2,
			BackOff:    1 * time.Second,
			Expiration: 5 * time.Second,
		},
	})})

	for _, writer := range writers {
		tracer.RegisterWriter(writer)
	}
}

func inner() {
	logger := tracer.GetLogger("moduleA.inner")
	logger.Info("don't know which transaction is this")
	logger.Info("but this log in this transaction")
	logger = logger.Trace()
	go func() {
		logger.Info("this is also inside the same transaction")
		func() {
			logger := tracer.GetLogger("moduleA.inner.nested")
			logger.Info("but not this one...")

		}()
	}()
}

func main() {
	configureTracer()
	logger := tracer.GetLogger("moduleA")
	logger.Info("logging in transaction 'A'", "B")
	logger.Info("logging in transaction 'B'", "B")
	logger.Info("logging in transaction 'B'", "B")
	logger.Info("logging in transaction 'A'", "A")
	logger.Info("logging in transaction 'A'", "A")
	logger = logger.Trace("C") // now all logs on this logger will be on the transaction C
	logger.Info("logging in transaction 'C'")
	logger.Info("logging in transaction 'C'", "A")
	inner()
	fmt.Scanln()
}

```

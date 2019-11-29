# Splunk Writer
This library is a writer(sink) [Splunk](https://www.splunk.com/en_us) to use with [Tracer](https://github.com/mralves/tracer)

## How to install
Using go get (not recommended):
```bash
go get github.com/mundipagg/tracer-splunk-writer
```

Using [dep](github.com/golang/dep) (recommended):
```bash
dep ensure --add github.com/mundipagg/tracer-splunk-writer@<version>
```

## Configuration

|Field|Type|Mandatory?|Default|Description|
|---|---|---|:---:|---|
|Address|string|Y||Splunk **full** endpoint (i.e. http://localhost:8088/services/collector)|
|Key|string|N|""|Splunk [Token Key](https://docs.splunk.com/Documentation/Splunk/8.0.0/Data/UsetheHTTPEventCollector)|
|Application|string|Y||Application name|
|Minimum Level|uint8|N|DEBUG|Minimum Level to log following the [syslog](https://en.wikipedia.org/wiki/Syslog#Severity_level) standard|
|Timeout|time.Duration|N|0 (infinite)|Timeout of the HTTP client|
|MessageEnvelop|string|N|"%v"|A envelop that *wraps* the original message|
|DefaultProperties|splunk.Entry|N|{}|A generic object to append to *every* log entry, but can be overwritten by the original log entry|
|Buffer.Cap|int|N|100| Maximum capacity of the log buffer, when the buffer is full all logs are sent at once|
|Buffer.OnWait|int|N|100| Maximum size of the queue to send to Splunk|
|Buffer.BackOff|time.Duration|N|60 seconds| Delay between retries to Splunk|
|ConfigLineLog |  map[string]interface{} | S | {} | Properties needed to insert log in splunk (host, source, sourcetype and index) 
|DefaultPropertiesSplunk | map[string]interface{} | S | {} | Properties set by administrador on splunk
|DefaultPropertiesApp | map[string]interface{} | S | {} | Properties to information about your application

## How to use

Below follows a simple example of how to use this lib:

```go
package main

import (
	"fmt"
	splunk "github.com/mundipagg/tracer-splunk-writer"
	"time"

	"github.com/mralves/tracer"


	bsp "github.com/mundipagg/tracer-splunk-writer/buffer"
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
	writers = append(writers, &Safe{splunk.New(splunk.Config{
		Timeout:      3 * time.Second,
		MinimumLevel: tracer.Debug,
		ConfigLineLog: LogEntry{
			"host":       "MyMachineName",
			"source":     "MySourceLog",
			"sourcetype": "_json",
			"index":      "main",
		},
		DefaultPropertiesSplunk: LogEntry{
			"ProcessName":    "MyProcessSourceLog",
			"ProductCompany": "MyCompanyName",
			"ProductName":    "MyProductName",
			"ProductVersion": "1.0",
		},
		DefaultPropertiesApp: LogEntry{
			"Properties": "Add here your default properties.",

		},
		Application: "MyApplicationName",
		Key:         "dd4a1733-75c8-48b2-abba-c102af7b9523",
		Address:     "http://localhost:8088/services/collector",
		Buffer: bsp.Config{
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

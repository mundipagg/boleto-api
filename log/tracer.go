package log

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/mundipagg/boleto-api/config"
	seq "github.com/mundipagg/tracer-seq-writer"
	splunk "github.com/mundipagg/tracer-splunk-writer"

	"github.com/mralves/tracer"

	bsq "github.com/mundipagg/tracer-seq-writer/buffer"
	bsp "github.com/mundipagg/tracer-splunk-writer/buffer"
)

type Safe struct {
	tracer.Writer
}

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

	WaitTimeLog := toInt(config.Get().WaitSecondsRetentationLog, 1)

	if config.Get().SeqEnabled == true {
		writers = append(writers, &Safe{seq.New(seq.Config{
			Timeout:      3 * time.Second,
			MinimumLevel: tracer.Debug,
			DefaultProperties: LogEntry{
				"Application": config.Get().ApplicationName,
				"Environment": config.Get().Environment,
				"Domain":      config.Get().SEQDomain,
				"MachineName": config.Get().MachineName,
			},
			Application: config.Get().ApplicationName,
			Key:         config.Get().SEQAPIKey,
			Address:     config.Get().SEQUrl,
			Buffer: bsq.Config{
				OnWait:     2,
				BackOff:    time.Duration(WaitTimeLog) * time.Second,
				Expiration: 5 * time.Second,
			},
		})})
	}

	if config.Get().SplunkEnabled == true {
		writers = append(writers, &Safe{splunk.New(splunk.Config{
			Timeout:      3 * time.Second,
			MinimumLevel: tracer.Debug,
			ConfigLineLog: LogEntry{
				"host":       config.Get().MachineName,
				"source":     "BoletoOnline",
				"sourcetype": config.Get().SplunkSourceType,
				"index":      config.Get().SplunkIndex,
			},
			DefaultPropertiesSplunk: LogEntry{
				"ProcessName":    "BoletoApi",
				"ProductCompany": "Mundipagg",
				"ProductName":    "BoletoOnline",
				"ProductVersion": "1.0",
			},
			DefaultPropertiesApp: LogEntry{
				"Application": config.Get().ApplicationName,
				"Environment": config.Get().Environment,
				"Domain":      config.Get().SEQDomain,
				"MachineName": config.Get().MachineName,
			},
			Application: config.Get().ApplicationName,
			Key:         config.Get().SplunkKey,
			Address:     config.Get().SplunkAddress,
			Buffer: bsp.Config{
				OnWait:     2,
				BackOff:    time.Duration(WaitTimeLog) * time.Second,
				Expiration: 5 * time.Second,
			},
		})})
	}

	for _, writer := range writers {
		tracer.RegisterWriter(writer)
	}
}

func toInt(str string, defaultValue ...int) int {
	if isBlank(str) {
		return 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		panic(err)
	}
	return i
}

var emptyOrWhitespacePattern = regexp.MustCompile(`^\s*$`)

// Function to check if a string is empty or contain only whitespaces.
func isBlank(str string) bool {
	return emptyOrWhitespacePattern.MatchString(str)
}

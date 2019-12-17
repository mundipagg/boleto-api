package splunk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/mralves/tracer"
	"github.com/mundipagg/tracer-splunk-writer/buffer"
	"github.com/mundipagg/tracer-splunk-writer/json"
	s "github.com/mundipagg/tracer-splunk-writer/strings"
)

type Writer struct {
	sync.Locker
	address                 string
	key                     string
	configLineLog           map[string]interface{}
	defaultPropertiesSplunk map[string]interface{}
	defaultPropertiesApp    map[string]interface{}
	client                  *http.Client
	buffer                  buffer.Buffer
	minimumLevel            uint8
	marshaller              jsoniter.API
	messageEnvelop          string
}

var punctuation = regexp.MustCompile(`(.+?)[?;:\\.,!]?$`)

//Used when message contains properties to replace.
var r = strings.NewReplacer("{", "{{.", "}", "}}")

func (sw *Writer) Write(entry tracer.Entry) {
	go func(sw *Writer, entry tracer.Entry) {
		defer func() {
			if err := recover(); err != nil {
				stderr("COULD NOT SEND SPLUNK TO SPLUNK BECAUSE %v", err)
			}
		}()
		if entry.Level > sw.minimumLevel {
			return
		}

		properties := NewEntry(append(entry.Args, sw.defaultPropertiesApp))
		message := punctuation.FindStringSubmatch(s.Capitalize(entry.Message))[1]
		message = s.ProcessString(r.Replace(message), properties)

		l := NewEntry(sw.configLineLog)
		l.Add("time", time.Now().UTC().UnixNano())
		e := NewEntry(Entry{
			"AdditionalData": properties,
			"Message":        message,
			"Severity":       Level(entry.Level),
		}, sw.defaultPropertiesSplunk)
		l.Add("event", e)

		sw.buffer.Write(l)
	}(sw, entry)
}

func (sw *Writer) send(events []interface{}) error {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	body, err := sw.marshaller.Marshal(events)
	if err != nil {
		stderr("COULD NOT SEND LOG TO SPLUNK BECAUSE %v", err)
		return err
	}

	request, _ := http.NewRequest(http.MethodPost, sw.address, bytes.NewBuffer(body))
	if len(sw.key) > 0 {
		request.Header.Set("Authorization", "Splunk "+sw.key)
	}
	request.Header.Set("Content-Type", "application/json")

	var response *http.Response
	response, err = sw.client.Do(request)
	if err != nil {
		stderr("COULD NOT SEND LOG TO SPLUNK BECAUSE %v", err)
		return err
	}
	response.Body.Close()
	if response.StatusCode != 200 {
		stderr("COULD NOT SEND LOG TO SPLUNK BECAUSE %v", response.Status)
		return errors.New(fmt.Sprintf("request returned %v", response.StatusCode))
	}

	return nil
}

func stderr(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message+"\n", args...)
}

type Config struct {
	Address                 string
	Key                     string
	Application             string
	Buffer                  buffer.Config
	MinimumLevel            uint8
	Timeout                 time.Duration
	ConfigLineLog           Entry
	DefaultPropertiesSplunk Entry
	DefaultPropertiesApp    Entry
	MessageEnvelop          string
}

func New(config Config) *Writer {
	writer := Writer{
		Locker:  &sync.RWMutex{},
		address: config.Address,
		key:     config.Key,
		client: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				TLSHandshakeTimeout: config.Timeout,
				IdleConnTimeout:     config.Timeout,
			},
		},
		messageEnvelop:          config.MessageEnvelop,
		minimumLevel:            config.MinimumLevel,
		configLineLog:           config.ConfigLineLog,
		defaultPropertiesSplunk: config.DefaultPropertiesSplunk,
		defaultPropertiesApp:    config.DefaultPropertiesApp,
		marshaller:              json.NewWithCaseStrategy(s.UseAnnotation),
	}
	config.Buffer.OnOverflow = writer.send
	writer.buffer = buffer.New(config.Buffer)
	return &writer
}

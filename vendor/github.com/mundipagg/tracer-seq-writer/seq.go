package seq

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/mralves/tracer"
	"github.com/mundipagg/tracer-seq-writer/buffer"
	"github.com/mundipagg/tracer-seq-writer/json"
	"github.com/mundipagg/tracer-seq-writer/strings"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

type Writer struct {
	sync.Locker
	address           string
	key               string
	defaultProperties map[string]interface{}
	client            *http.Client
	buffer            buffer.Buffer
	minimumLevel      uint8
	marshaller        jsoniter.API
	messageEnvelop    string
}

type seqLog struct {
	Events []interface{} `json:"Events"`
}

type event struct {
	Timestamp       string      `json:"Timestamp"`
	Level           string      `json:"Level"`
	MessageTemplate string      `json:"MessageTemplate"`
	Properties      interface{} `json:"Properties"`
}

var punctuation = regexp.MustCompile(`(.+?)[?;:\\.,!]?$`)

func (sw *Writer) Write(entry tracer.Entry) {
	go func(sw *Writer, entry tracer.Entry) {
		defer func() {
			if err := recover(); err != nil {
				stderr("COULD NOT SEND LOG TO SEQ BECAUSE %v", err)
			}
		}()
		if entry.Level > sw.minimumLevel {
			return
		}
		extraProperties := map[string]interface{}{
			"RequestKey": entry.TransactionId,
			"Caller":     entry.StackTrace[0].String(),
		}
		if strings.IsBlank(entry.TransactionId) {
			delete(extraProperties, "RequestKey")
		}
		properties := NewEntry(append(entry.Args, extraProperties, sw.defaultProperties))
		message := strings.Capitalize(entry.Message)
		message = punctuation.FindStringSubmatch(message)[1]
		if len(sw.messageEnvelop) > 0 {
			message = fmt.Sprintf(sw.messageEnvelop, message)
		}
		e := event{
			Level:           Level(entry.Level),
			Timestamp:       entry.Time.UTC().Format(time.RFC3339Nano),
			MessageTemplate: message,
			Properties:      properties,
		}
		sw.buffer.Write(e)
	}(sw, entry)
}

func (sw *Writer) send(events []interface{}) error {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()
	log := seqLog{
		Events: events,
	}
	body, err := sw.marshaller.Marshal(log)
	if err != nil {
		stderr("COULD NOT SEND LOG TO SEQ BECAUSE %v", err)
		return err
	}

	request, _ := http.NewRequest(http.MethodPost, sw.address, bytes.NewBuffer(body))
	if len(sw.key) > 0 {
		request.Header.Set("X-Seq-ApiKey", sw.key)
	}
	request.Header.Set("Content-Type", "application/json")
	var response *http.Response
	response, err = sw.client.Do(request)
	if err != nil {
		stderr("COULD NOT SEND LOG TO SEQ BECAUSE %v", err)
		return err
	}
	response.Body.Close()
	if response.StatusCode != 201 {
		stderr("COULD NOT SEND LOG TO SEQ BECAUSE %v", response.Status)
		return errors.New(fmt.Sprintf("request returned %v", response.StatusCode))
	}
	return nil
}

func stderr(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message+"\n", args...)
}

type Config struct {
	Address           string
	Key               string
	Application       string
	Buffer            buffer.Config
	MinimumLevel      uint8
	Timeout           time.Duration
	DefaultProperties Entry
	MessageEnvelop    string
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
		messageEnvelop:    config.MessageEnvelop,
		minimumLevel:      uint8(config.MinimumLevel),
		defaultProperties: config.DefaultProperties,
		marshaller:        json.NewWithCaseStrategy(strings.ToPascalCase),
	}
	config.Buffer.OnOverflow = writer.send
	writer.buffer = buffer.New(config.Buffer)
	return &writer
}

package log

import (
	"fmt"
	"github.com/mralves/tracer"
	"github.com/mundipagg/boleto-api/util"
	"strings"

	"github.com/mundipagg/boleto-api/config"
)

type LogEntry = map[string]interface{}

var logger tracer.Logger

//Operation a operacao usada na API
var Operation string

//Recipient o nome do banco
var Recipient string

//Log struct com os elemtos do log
type Log struct {
	Operation   string
	Recipient   string
	RequestKey  string
	BankName    string
	IPAddress   string
	ServiceUser string
	NossoNumero uint
	ElapsedTime int64
	logger      tracer.Logger
}

//Install instala o "servico" de log do SEQ
func Install() {
	configureTracer()
	logger = tracer.GetLogger("boleto")
}

func formatter(message string) string {
	return "[{Application}: {Operation}] - {MessageType} " + message
}

//CreateLog cria uma nova instancia do Log
func CreateLog() *Log {
	return &Log{
		logger: logger,
	}
}

//Request loga o request para algum banco
func (l *Log) Request(content interface{}, url string, headers map[string]string) {
	if config.Get().DisableLog {
		return
	}

	go (func() {
		props := l.defaultProperties("Request", content)
		props["Headers"] = headers
		props["URL"] = url
		props["ExecutionDate"] = util.UtcNow()

		action := strings.Split(url, "/")
		msg := formatter(fmt.Sprintf("to {BankName} (%s) | {Recipient}", action[len(action)-1]))

		l.logger.Info(msg, props)
	})()
}

//Response loga o response para algum banco
func (l *Log) Response(content interface{}, url string) {
	if config.Get().DisableLog {
		return
	}

	go (func() {

		action := strings.Split(url, "/")
		msg := formatter(fmt.Sprintf("from {BankName} (%s) | {Recipient}", action[len(action)-1]))

		props := l.defaultProperties("Response", content)
		props["URL"] = url
		props["ExecutionDate"] = util.UtcNow()
		props["ElapsedTimeInMilliseconds"] = l.ElapsedTime

		l.logger.Info(msg, props)
	})()
}

//RequestApplication loga o request que chega na boleto api
func (l *Log) RequestApplication(content interface{}, url string, headers map[string]string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Request", content)
		props["Headers"] = headers
		props["URL"] = url
		props["ExecutionDate"] = util.UtcNow()

		msg := formatter("from {IPAddress} | {Recipient}")

		l.logger.Info(msg, props)
	})()
}

//ResponseApplication loga o response que sai da boleto api
func (l *Log) ResponseApplication(content interface{}, url string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Response", content)
		props["URL"] = url
		props["ExecutionDate"] = util.UtcNow()

		msg := formatter("{Operation} | {Recipient}")

		l.logger.Info(msg, props)
	})()
}

//Info loga mensagem do level INFO
func (l *Log) Info(msg string) {
	if config.Get().DisableLog {
		return
	}

	props := LogEntry{}
	props["ExecutionDate"] = util.UtcNow()

	go l.logger.Info(msg, props)
}

//Warn loga mensagem do leve Warning
func (l *Log) Warn(content interface{}, msg string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Warning", content)
		m := formatter(msg)
		props["ExecutionDate"] = util.UtcNow()

		l.logger.Warn(m, props)
	})()
}

func (l *Log) Error(content interface{}, msg string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Error", content)
		m := formatter(msg)
		props["ExecutionDate"] = util.UtcNow()

		l.logger.Error(m, props)
	})()
}

// Fatal loga erros da aplicação
func (l *Log) Fatal(content interface{}, msg string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Fatal", content)
		m := formatter(msg)
		props["ExecutionDate"] = util.UtcNow()

		l.logger.Fatal(m, props)
	})()
}

//InitRobot loga o inicio da execução do robô de recovery
func (l *Log) InitRobot() {
	go (func() {
		props := defaultRobotProperties("Execute", l.Operation, "")
		msg := formatter("- Starting execution")
		props["ExecutionDate"] = util.UtcNow()
		logger.Info(msg, props)
	})()
}

//ResumeRobot loga um resumo de Recovery do robô de recovery
func (l *Log) ResumeRobot(key string) {
	go (func() {
		props := defaultRobotProperties("RecoveryBoleto", l.Operation, key)
		msg := formatter(key)
		props["ExecutionDate"] = util.UtcNow()
		logger.Info(msg, props)
	})()
}

//EndRobot loga o fim da execução do robô de recovery
func (l *Log) EndRobot() {
	go (func() {
		props := defaultRobotProperties("Finish", l.Operation, "")
		msg := formatter("- Finishing execution")
		props["ExecutionDate"] = util.UtcNow()
		logger.Info(msg, props)
	})()
}

func (l *Log) defaultProperties(messageType string, content interface{}) LogEntry {
	props := LogEntry{
		"MessageType": messageType,
		"Content":     content,
		"Recipient":   l.Recipient,
		"Operation":   l.Operation,
		"NossoNumero": l.NossoNumero,
		"RequestKey":  l.RequestKey,
		"BankName":    l.BankName,
		"IPAddress":   l.IPAddress,
		"ServiceUser": l.ServiceUser,
	}
	return props
}

func defaultRobotProperties(msgType, op, key string) LogEntry {
	props := LogEntry{
		"MessageType": msgType,
		"Operation":   op,
	}

	if key != "" {
		props["BoletoKey"] = key
	}
	return props
}

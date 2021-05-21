package log

import (
	"fmt"
	"strings"

	"github.com/mralves/tracer"

	"github.com/mundipagg/boleto-api/config"
)

const (
	FailGetBoletoMessage    = "Falha ao recuperar Boleto"
	SuccessGetBoletoMessage = "Boleto recuperado com sucesso"
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

		msg := formatter("{Operation} | {Recipient}")

		l.logger.Info(msg, props)
	})()
}

//Info loga mensagem do level INFO
func (l *Log) Info(msg string) {
	if config.Get().DisableLog {
		return
	}
	go l.logger.Info(msg, nil)
}

//Info loga mensagem do level INFO
func Info(msg string) {
	if config.Get().DisableLog {
		return
	}
	go logger.Info(msg, nil)
}

//Warn loga mensagem do leve Warning
func (l *Log) Warn(content interface{}, msg string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Warning", content)
		m := formatter(msg)

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

		l.logger.Fatal(m, props)
	})()
}

//InitRobot loga o inicio da execução do robô de recovery
func (l *Log) InitRobot(totalRecords int) {
	msg := formatter("- Starting execution")
	go func() {
		props := defaultRobotProperties("Execute", l.Operation, "")
		props["TotalRecords"] = totalRecords
		logger.Info(msg, props)
	}()
}

//ResumeRobot loga um resumo de Recovery do robô de recovery
func (l *Log) ResumeRobot(key string) {
	msg := formatter(key)
	go func() {
		props := defaultRobotProperties("RecoveryBoleto", l.Operation, key)
		props["RequestKey"] = l.RequestKey
		logger.Info(msg, props)
	}()
}

//EndRobot loga o fim da execução do robô de recovery
func (l *Log) EndRobot() {
	msg := formatter("- Finishing execution")
	go logger.Info(msg, defaultRobotProperties("Finish", l.Operation, ""))
}

func (l *Log) defaultProperties(messageType string, content interface{}) LogEntry {
	props := LogEntry{
		"Content":     content,
		"Recipient":   l.Recipient,
		"NossoNumero": l.NossoNumero,
		"RequestKey":  l.RequestKey,
		"BankName":    l.BankName,
		"ServiceUser": l.ServiceUser,
	}

	for k, v := range l.basicProperties(messageType) {
		props[k] = v
	}

	return props
}

//GetBoleto Loga mensagem de recuperação de boleto
func (l *Log) GetBoleto(content interface{}, msgType string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.basicProperties(msgType)
		props["Content"] = content

		switch msgType {
		case "Warning":
			l.logger.Warn(formatter(FailGetBoletoMessage), props)
		case "Error":
			l.logger.Error(formatter(FailGetBoletoMessage), props)
		default:
			l.logger.Info(formatter(SuccessGetBoletoMessage), props)
		}
	})()
}

func (l *Log) basicProperties(messageType string) LogEntry {
	props := LogEntry{
		"MessageType": messageType,
		"Operation":   l.Operation,
		"IPAddress":   l.IPAddress,
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

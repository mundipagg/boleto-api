package log

import (
	"fmt"
	"strings"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/goseq"
)

var logger *goseq.Logger

// Operation a operacao usada na API
var Operation string

// Recipient o nome do banco
var Recipient string

// Log struct com os elemtos do log
type Log struct {
	Operation   string
	Recipient   string
	RequestKey  string
	BankName    string
	IPAddress   string
	NossoNumero uint
	logger      *goseq.Logger
}

//Install instala o "servico" de log do SEQ
func Install() error {

	_logger, err := goseq.GetLogger(config.Get().SEQUrl, config.Get().SEQAPIKey)
	if err != nil {
		return err
	}
	_logger.SetDefaultProperties(map[string]interface{}{
		"Application": config.Get().ApplicationName,
		"Environment": config.Get().Environment,
		"Domain":      config.Get().SEQDomain,
		"MachineName": config.Get().MachineName,
	})
	logger = _logger
	return nil
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

// Request loga o request para algum banco
func (l Log) Request(content interface{}, url string, headers map[string]string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Request", content)
		props.AddProperty("Headers", headers)
		props.AddProperty("URL", url)
		action := strings.Split(url, "/")
		msg := formatter(fmt.Sprintf("to {BankName} (%s) | {Recipient}", action[len(action)-1]))

		l.logger.Information(msg, props)
	})()
}

// Response loga o response para algum banco
func (l Log) Response(content interface{}, url string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Response", content)
		props.AddProperty("URL", url)
		action := strings.Split(url, "/")
		msg := formatter(fmt.Sprintf("from {BankName} (%s) | {Recipient}", action[len(action)-1]))

		l.logger.Information(msg, props)
	})()
}

// Request loga o request que chega na boleto api
func (l Log) RequestApplication(content interface{}, url string, headers map[string]string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Request", content)
		props.AddProperty("Headers", headers)
		props.AddProperty("URL", url)
		msg := formatter("from {IPAddress} | {Recipient}")

		l.logger.Information(msg, props)
	})()
}

// Response loga o response que sai da boleto api
func (l Log) ResponseApplication(content interface{}, url string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Response", content)
		props.AddProperty("URL", url)
		msg := formatter(" {Operation} | {Recipient}")

		l.logger.Information(msg, props)
	})()
}

//Info loga mensagem do level INFO
func (l Log) Info(msg string) {
	if config.Get().DisableLog {
		return
	}
	go logger.Information(msg, goseq.NewProperties())
}

func Info(msg string) {
	if config.Get().DisableLog {
		return
	}
	go logger.Information(msg, goseq.NewProperties())
}

//Warn loga mensagem do leve Warning
func (l Log) Warn(content interface{}, msg string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Warning", content)
		m := formatter(msg)

		l.logger.Warning(m, props)
	})()
}

// Fatal loga erros da aplicação
func (l Log) Fatal(content interface{}, msg string) {
	if config.Get().DisableLog {
		return
	}
	go (func() {
		props := l.defaultProperties("Error", content)
		m := formatter(msg)

		l.logger.Fatal(m, props)
	})()
}

//InitRobot loga o inicio da execução do robô de recovery
func (l Log) InitRobot() {
	msg := formatter("- Starting execution")
	go logger.Information(msg, defaultRobotProperties("Execute", l.Operation, ""))
}

//ResumeRobot loga um resumo de Recovery do robô de recovery
func (l Log) ResumeRobot(key string) {
	msg := formatter(key)
	go logger.Information(msg, defaultRobotProperties("RecoveryBoleto", l.Operation, key))
}

//EndRobot loga o fim da execução do robô de recovery
func (l Log) EndRobot() {
	msg := formatter("- Finishing execution")
	go logger.Information(msg, defaultRobotProperties("Finish", l.Operation, ""))
}

func (l Log) defaultProperties(messageType string, content interface{}) goseq.Properties {
	props := goseq.NewProperties()
	props.AddProperty("MessageType", messageType)
	props.AddProperty("Content", content)
	props.AddProperty("Recipient", l.Recipient)
	props.AddProperty("Operation", l.Operation)
	props.AddProperty("NossoNumero", l.NossoNumero)
	props.AddProperty("RequestKey", l.RequestKey)
	props.AddProperty("BankName", l.BankName)
	props.AddProperty("IPAddress", l.IPAddress)
	return props
}

func defaultRobotProperties(msgType, op, key string) goseq.Properties {
	props := goseq.NewProperties()
	props.AddProperty("MessageType", msgType)
	props.AddProperty("Operation", op)

	if key != "" {
		props.AddProperty("BoletoKey", key)
	}

	return props
}

//Close fecha a conexao com o SEQ
func Close() {
	if !config.Get().DisableLog && logger.Async {
		fmt.Println("Closing SEQ Connection")
		logger.Close()
	}
}

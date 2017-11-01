package config

import (
	"os"
	"sync"
	"sync/atomic"
)

//Config é a estrutura que tem todas as configurações da aplicação
type Config struct {
	InfluxDBHost               string
	InfluxDBPort               string
	APIPort                    string
	PdfAPIURL                  string
	Version                    string
	SEQUrl                     string
	SEQAPIKey                  string
	EnableRequestLog           bool
	EnablePrintRequest         bool
	Environment                string
	SEQDomain                  string
	ApplicationName            string
	URLBBRegisterBoleto        string
	URLCaixaRegisterBoleto     string
	URLBBToken                 string
	URLCitiBoleto              string
	URLCiti                    string
	MockMode                   bool
	DevMode                    bool
	HTTPOnly                   bool
	AppURL                     string
	ElasticURL                 string
	MongoURL                   string
	MongoUser                  string
	MongoPassword              string
	BoletoJSONFileStore        string
	DisableLog                 bool
	CertBoletoPathCrt          string
	CertBoletoPathKey          string
	CertBoletoPathCa           string
	URLTicketSantander         string
	URLRegisterBoletoSantander string
	URLBradesco                string
}

var cnf Config
var scnf sync.Once
var running uint64
var mutex sync.Mutex

//Get retorna o objeto de configurações da aplicação
func Get() Config {
	return cnf
}
func Install(mockMode, devMode, disableLog bool) {
	atomic.StoreUint64(&running, 0)
	cnf = Config{
		APIPort:                    ":" + os.Getenv("API_PORT"),
		PdfAPIURL:                  os.Getenv("PDF_API"),
		Version:                    os.Getenv("API_VERSION"),
		SEQUrl:                     os.Getenv("SEQ_URL"),                        //Pegar o SEQ de dev
		SEQAPIKey:                  os.Getenv("SEQ_API_KEY"),                    //Staging Key:
		EnableRequestLog:           os.Getenv("ENABLE_REQUEST_LOG") == "true",   // Log a cada request no SEQ
		EnablePrintRequest:         os.Getenv("ENABLE_PRINT_REQUEST") == "true", // Imprime algumas informacoes da request no console
		Environment:                os.Getenv("ENVIRONMENT"),
		SEQDomain:                  "One",
		ApplicationName:            "BoletoOnline",
		URLBBRegisterBoleto:        os.Getenv("URL_BB_REGISTER_BOLETO"),
		URLCaixaRegisterBoleto:     os.Getenv("URL_CAIXA"),
		URLBBToken:                 os.Getenv("URL_BB_TOKEN"),
		URLCitiBoleto:              os.Getenv("URL_CITI_BOLETO"),
		URLCiti:                    os.Getenv("URL_CITI"),
		MockMode:                   mockMode,
		AppURL:                     os.Getenv("APP_URL"),
		ElasticURL:                 os.Getenv("ELASTIC_URL"),
		DevMode:                    devMode,
		DisableLog:                 disableLog,
		MongoURL:                   os.Getenv("MONGODB_URL"),
		MongoUser:                  os.Getenv("MONGODB_USER"),
		MongoPassword:              os.Getenv("MONGODB_PASSWORD"),
		BoletoJSONFileStore:        os.Getenv("BOLETO_JSON_STORE"),
		CertBoletoPathCrt:          os.Getenv("CERT_BOLETO_CRT"),
		CertBoletoPathKey:          os.Getenv("CERT_BOLETO_KEY"),
		CertBoletoPathCa:           os.Getenv("CERT_BOLETO_CA"),
		URLTicketSantander:         os.Getenv("URL_SANTANDER_TICKET"),
		URLRegisterBoletoSantander: os.Getenv("URL_SANTANDER_REGISTER"),
		URLBradesco:                os.Getenv("URL_BRADESCO"),
		InfluxDBHost:               os.Getenv("INFLUXDB_HOST"),
		InfluxDBPort:               os.Getenv("INFLUXDB_PORT"),
	}
}

//IsRunning verifica se a aplicação tem que aceitar requisições
func IsRunning() bool {
	return atomic.LoadUint64(&running) > 0
}

//IsNotProduction returns true if application is running in DevMode or MockMode
func IsNotProduction() bool {
	return cnf.DevMode || cnf.MockMode
}

//Stop faz a aplicação parar de receber requisições
func Stop() {
	atomic.StoreUint64(&running, 1)
}

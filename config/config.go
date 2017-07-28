package config

import (
	"os"
	"sync"
	"sync/atomic"
)

//Config é a estrutura que tem todas as configurações da aplicação
type Config struct {
	APIPort                string
	Version                string
	SEQUrl                 string
	SEQAPIKey              string
	EnableRequestLog       bool
	EnablePrintRequest     bool
	Environment            string
	SEQDomain              string
	ApplicationName        string
	URLBBRegisterBoleto    string
	URLCaixaRegisterBoleto string
	URLBBToken             string
	URLCiti                string
	MockMode               bool
	DevMode                bool
	HTTPOnly               bool
	AppURL                 string
	ElasticURL             string
	MongoURL               string
	BoletoJSONFileStore    string
	DisableLog             bool
	TLSCertPath            string
	TLSKeyPath             string
}

var cnf Config
var scnf sync.Once
var running uint64
var mutex sync.Mutex

//Get retorna o objeto de configurações da aplicação
func Get() Config {
	return cnf
}
func Install(mockMode, devMode, disableLog, httpOnly bool) {
	atomic.StoreUint64(&running, 0)
	cnf = Config{
		APIPort:                ":" + os.Getenv("API_PORT"),
		Version:                os.Getenv("API_VERSION"),
		SEQUrl:                 os.Getenv("SEQ_URL"),                        //Pegar o SEQ de dev
		SEQAPIKey:              os.Getenv("SEQ_API_KEY"),                    //Staging Key:
		EnableRequestLog:       os.Getenv("ENABLE_REQUEST_LOG") == "true",   // Log a cada request no SEQ
		EnablePrintRequest:     os.Getenv("ENABLE_PRINT_REQUEST") == "true", // Imprime algumas informacoes da request no console
		Environment:            os.Getenv("ENVIRONMENT"),
		SEQDomain:              "One",
		ApplicationName:        "BoletoOnline",
		URLBBRegisterBoleto:    os.Getenv("URL_BB_REGISTER_BOLETO"),
		URLCaixaRegisterBoleto: os.Getenv("URL_CAIXA"),
		URLBBToken:             os.Getenv("URL_BB_TOKEN"),
		URLCiti:                os.Getenv("URL_CITI"),
		MockMode:               mockMode,
		AppURL:                 os.Getenv("APP_URL"),
		ElasticURL:             os.Getenv("ELASTIC_URL"),
		DevMode:                devMode,
		DisableLog:             disableLog,
		HTTPOnly:               httpOnly,
		MongoURL:               os.Getenv("MONGODB_URL"),
		BoletoJSONFileStore:    os.Getenv("BOLETO_JSON_STORE"),
		TLSCertPath:            os.Getenv("TLS_CERT_PATH"),
		TLSKeyPath:             os.Getenv("TLS_KEY_PATH"),
	}
}

//IsRunning verifica se a aplicação tem que aceitar requisições
func IsRunning() bool {
	return atomic.LoadUint64(&running) > 0
}

//Stop faz a aplicação parar de receber requisições
func Stop() {
	atomic.StoreUint64(&running, 1)
}

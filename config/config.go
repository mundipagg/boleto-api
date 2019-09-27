package config

import (
	"os"
	"sync"
	"sync/atomic"
)

//Config é a estrutura que tem todas as configurações da aplicação
type Config struct {
	InfluxDBHost                    string
	InfluxDBPort                    string
	APIPort                         string
	MachineName                     string
	PdfAPIURL                       string
	Version                         string
	SEQUrl                          string
	SEQAPIKey                       string
	EnableRequestLog                bool
	EnablePrintRequest              bool
	Environment                     string
	SEQDomain                       string
	ApplicationName                 string
	URLBBRegisterBoleto             string
	CaixaEnv                        string
	URLCaixaRegisterBoleto          string
	URLBBToken                      string
	URLCitiBoleto                   string
	URLCiti                         string
	MockMode                        bool
	DevMode                         bool
	HTTPOnly                        bool
	AppURL                          string
	ElasticURL                      string
	MongoURL                        string
	MongoUser                       string
	MongoPassword                   string
	RedisURL                        string
	RedisPassword                   string
	RedisDatabase                   string
	RedisExpirationTime             string
	RedisSSL                        bool
	BoletoJSONFileStore             string
	DisableLog                      bool
	CertBoletoPathCrt               string
	CertBoletoPathKey               string
	CertBoletoPathCa                string
	CertICP_PathPkey                string
	CertICP_PathChainCertificates   string
	URLTicketSantander              string
	URLRegisterBoletoSantander      string
	URLBradescoShopFacil            string
	URLBradescoNetEmpresa           string
	ItauEnv                         string
	SantanderEnv                    string
	URLTicketItau                   string
	URLRegisterBoletoItau           string
	RecoveryRobotExecutionEnabled   string
	RecoveryRobotExecutionInMinutes string
	TimeoutRegister                 string
	TimeoutToken                    string
	TimeoutDefault                  string
	URLPefisaToken                  string
	URLPefisaRegister               string
	EnableMetrics                   bool
	CertificatesPath                string
	AzureTenantId                   string
	AzureClientId                   string
	AzureClientSecret               string
	VaultName                       string
	CertificateICPName              string
	PswCertificateICP               string
	CertificateSSLName              string
	PswCertificateSSL               string
	EnableFileServerCertificate     bool
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
	hostName := getHostName()

	cnf = Config{
		APIPort:                         ":" + os.Getenv("API_PORT"),
		PdfAPIURL:                       os.Getenv("PDF_API"),
		Version:                         os.Getenv("API_VERSION"),
		MachineName:                     hostName,
		SEQUrl:                          os.Getenv("SEQ_URL"),                        //Pegar o SEQ de dev
		SEQAPIKey:                       os.Getenv("SEQ_API_KEY"),                    //Staging Key:
		EnableRequestLog:                os.Getenv("ENABLE_REQUEST_LOG") == "true",   // Log a cada request no SEQ
		EnablePrintRequest:              os.Getenv("ENABLE_PRINT_REQUEST") == "true", // Imprime algumas informacoes da request no console
		Environment:                     os.Getenv("ENVIRONMENT"),
		SEQDomain:                       "One",
		ApplicationName:                 "BoletoOnline",
		URLBBRegisterBoleto:             os.Getenv("URL_BB_REGISTER_BOLETO"),
		CaixaEnv:                        os.Getenv("CAIXA_ENV"),
		URLCaixaRegisterBoleto:          os.Getenv("URL_CAIXA"),
		URLBBToken:                      os.Getenv("URL_BB_TOKEN"),
		URLCitiBoleto:                   os.Getenv("URL_CITI_BOLETO"),
		URLCiti:                         os.Getenv("URL_CITI"),
		MockMode:                        mockMode,
		AppURL:                          os.Getenv("APP_URL"),
		ElasticURL:                      os.Getenv("ELASTIC_URL"),
		DevMode:                         devMode,
		DisableLog:                      disableLog,
		MongoURL:                        os.Getenv("MONGODB_URL"),
		MongoUser:                       os.Getenv("MONGODB_USER"),
		MongoPassword:                   os.Getenv("MONGODB_PASSWORD"),
		RedisURL:                        os.Getenv("REDIS_URL"),
		RedisPassword:                   os.Getenv("REDIS_PASSWORD"),
		RedisDatabase:                   os.Getenv("REDIS_DATABASE"),
		RedisExpirationTime:             os.Getenv("REDIS_EXPIRATION_TIME_IN_SECONDS"),
		RedisSSL:                        os.Getenv("REDIS_SSL") == "true",
		CertBoletoPathCrt:               os.Getenv("CERT_BOLETO_CRT"),
		CertBoletoPathKey:               os.Getenv("CERT_BOLETO_KEY"),
		CertBoletoPathCa:                os.Getenv("CERT_BOLETO_CA"),
		CertICP_PathPkey:                os.Getenv("CERT_ICP_BOLETO_KEY"),
		CertICP_PathChainCertificates:   os.Getenv("CERT_ICP_BOLETO_CHAIN_CA"),
		URLTicketSantander:              os.Getenv("URL_SANTANDER_TICKET"),
		URLRegisterBoletoSantander:      os.Getenv("URL_SANTANDER_REGISTER"),
		ItauEnv:                         os.Getenv("ITAU_ENV"),
		SantanderEnv:                    os.Getenv("SANTANDER_ENV"),
		URLTicketItau:                   os.Getenv("URL_ITAU_TICKET"),
		URLRegisterBoletoItau:           os.Getenv("URL_ITAU_REGISTER"),
		URLBradescoShopFacil:            os.Getenv("URL_BRADESCO_SHOPFACIL"),
		URLBradescoNetEmpresa:           os.Getenv("URL_BRADESCO_NET_EMPRESA"),
		InfluxDBHost:                    os.Getenv("INFLUXDB_HOST"),
		InfluxDBPort:                    os.Getenv("INFLUXDB_PORT"),
		RecoveryRobotExecutionEnabled:   os.Getenv("RECOVERYROBOT_EXECUTION_ENABLED"),
		RecoveryRobotExecutionInMinutes: os.Getenv("RECOVERYROBOT_EXECUTION_IN_MINUTES"),
		TimeoutRegister:                 os.Getenv("TIMEOUT_REGISTER"),
		TimeoutToken:                    os.Getenv("TIMEOUT_TOKEN"),
		TimeoutDefault:                  os.Getenv("TIMEOUT_DEFAULT"),
		URLPefisaToken:                  os.Getenv("URL_PEFISA_TOKEN"),
		URLPefisaRegister:               os.Getenv("URL_PEFISA_REGISTER"),
		EnableMetrics:                   os.Getenv("ENABLE_METRICS") == "true",
		CertificatesPath:                os.Getenv("PATH_CERTIFICATES"),
		AzureTenantId:                   os.Getenv("AZURE_TENANT_ID"),
		AzureClientId:                   os.Getenv("AZURE_CLIENT_ID"),
		AzureClientSecret:               os.Getenv("AZURE_CLIENT_SECRET"),
		VaultName:                       os.Getenv("VAULT_NAME"),
		CertificateICPName:              os.Getenv("CERTIFICATE_ICP_NAME"),
		PswCertificateICP:               os.Getenv("PSW_CERTIFICATE_ICP_NAME"),
		CertificateSSLName:              os.Getenv("CERTIFICATE_SSL_NAME"),
		PswCertificateSSL:               os.Getenv("PSW_CERTIFICATE_SSL_NAME"),
		EnableFileServerCertificate:     os.Getenv("ENABLE_FILESERVER_CERTIFICATE") == "true",
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

func getHostName() string {
	machineName, err := os.Hostname()
	if err != nil {
		return ""
	}
	return machineName
}

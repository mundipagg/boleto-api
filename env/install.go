package env

import (
	"os"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

//Config Realiza a configuração da aplicação
func Config(devMode, mockMode, disableLog bool) {
	configFlags(devMode, mockMode, disableLog)
	flow.RegisterConnector("logseq", util.SeqLogConector)
	flow.RegisterConnector("apierro", models.BoletoErrorConector)
	flow.RegisterConnector("tls", util.TlsConector)
	metrics.Install()
}

//ConfigMock Criar configurações de desenvolvimento
func ConfigMock(port string) {
	os.Setenv("URL_BB_REGISTER_BOLETO", "http://localhost:"+port+"/registrarBoleto")
	os.Setenv("URL_BB_TOKEN", "http://localhost:"+port+"/oauth/token")
	os.Setenv("URL_CAIXA", "http://localhost:"+port+"/caixa/registrarBoleto")
	os.Setenv("URL_CITI", "http://localhost:"+port+"/citi/registrarBoleto")
	os.Setenv("URL_SANTANDER_TICKET", "tls://localhost:"+port+"/santander/get-ticket")
	os.Setenv("URL_SANTANDER_REGISTER", "tls://localhost:"+port+"/santander/register")
	os.Setenv("URL_BRADESCO_SHOPFACIL", "http://localhost:"+port+"/bradescoshopfacil/registrarBoleto")
	os.Setenv("URL_ITAU_TICKET", "http://localhost:"+port+"/itau/gerarToken")
	os.Setenv("URL_ITAU_REGISTER", "http://localhost:"+port+"/itau/registrarBoleto")
	os.Setenv("URL_BRADESCO_NET_EMPRESA", "http://localhost:"+port+"/bradesconetempresa/registrarBoleto")
	os.Setenv("URL_PEFISA_TOKEN", "http://localhost:"+port+"/pefisa/gerarToken")
	os.Setenv("URL_PEFISA_REGISTER", "http://localhost:"+port+"/pefisa/registrarBoleto")
	os.Setenv("MONGODB_URL", "localhost:27017")
	os.Setenv("MONGODB_USER", "")
	os.Setenv("MONGODB_PASSWORD", "")
	os.Setenv("REDIS_URL", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "123456")
	os.Setenv("REDIS_DATABASE", "8")
	os.Setenv("REDIS_SSL", "false")
	os.Setenv("REDIS_EXPIRATION_TIME_IN_SECONDS", "2880")
	os.Setenv("RECOVERYROBOT_EXECUTION_ENABLED", "true")
	os.Setenv("RECOVERYROBOT_EXECUTION_IN_MINUTES", "2")
	os.Setenv("SEQ_URL", "http://localhost:5341")
	os.Setenv("SEQ_API_KEY", "4jZzTybZ9bUHtJiPdh6")
	os.Setenv("TIMEOUT_REGISTER", "30")
	os.Setenv("TIMEOUT_TOKEN", "20")
	os.Setenv("TIMEOUT_DEFAULT", "50")
	config.Install(true, true, true)
}

func configFlags(devMode, mockMode, disableLog bool) {
	if devMode {
		os.Setenv("INFLUXDB_HOST", "http://localhost")
		os.Setenv("INFLUXDB_PORT", "8086")
		os.Setenv("PDF_API", "http://localhost:7070/topdf")
		os.Setenv("API_PORT", "3000")
		os.Setenv("API_VERSION", "0.0.1")
		os.Setenv("ENVIRONMENT", "Development")
		os.Setenv("SEQ_URL", "http://localhost:5341")
		os.Setenv("SEQ_API_KEY", "4jZzTybZ9bUHtJiPdh6")
		os.Setenv("ENABLE_REQUEST_LOG", "false")
		os.Setenv("ENABLE_PRINT_REQUEST", "true")
		os.Setenv("URL_BB_REGISTER_BOLETO", "https://cobranca.homologa.bb.com.br:7101/registrarBoleto")
		os.Setenv("URL_BB_TOKEN", "https://oauth.hm.bb.com.br/oauth/token")
		os.Setenv("CAIXA_ENV", "SGCBS01D")
		os.Setenv("URL_CAIXA", "https://des.barramento.caixa.gov.br/sibar/ManutencaoCobrancaBancaria/Boleto/Externo")
		os.Setenv("URL_CITI", "https://citigroupsoauat.citigroup.com/comercioeletronico/registerboleto/RegisterBoletoSOAP")
		os.Setenv("URL_CITI_BOLETO", "https://ebillpayer.uat.brazil.citigroup.com/ebillpayer/jspInformaDadosConsulta.jsp")
		os.Setenv("APP_URL", "http://localhost:3000/boleto")
		os.Setenv("ELASTIC_URL", "http://localhost:9200")
		os.Setenv("MONGODB_URL", "localhost:27017")
		os.Setenv("MONGODB_USER", "")
		os.Setenv("MONGODB_PASSWORD", "")
		os.Setenv("REDIS_URL", "localhost:6379")
		os.Setenv("REDIS_PASSWORD", "123456")
		os.Setenv("REDIS_DATABASE", "8")
		os.Setenv("REDIS_SSL", "false")
		os.Setenv("REDIS_EXPIRATION_TIME_IN_SECONDS", "2880")
		os.Setenv("CERT_BOLETO_CRT", "C:\\cert_boleto_api\\certificate.crt")
		os.Setenv("CERT_BOLETO_KEY", "C:\\cert_boleto_api\\pkey.key")
		os.Setenv("CERT_BOLETO_CA", "C:\\cert_boleto_api\\ca-cert.ca")
		os.Setenv("CERT_ICP_BOLETO_KEY", "C:\\cert_boleto_api\\ICP_PKey.key")
		os.Setenv("CERT_ICP_BOLETO_CHAIN_CA", "C:\\cert_boleto_api\\ICP_cadeiaCerts.pem")
		os.Setenv("URL_SANTANDER_TICKET", "https://ymbdlb.santander.com.br/dl-ticket-services/TicketEndpointService")
		os.Setenv("URL_SANTANDER_REGISTER", "https://ymbcash.santander.com.br/ymbsrv/CobrancaEndpointService")
		os.Setenv("URL_BRADESCO_SHOPFACIL", "https://homolog.meiosdepagamentobradesco.com.br/apiboleto/transacao")
		os.Setenv("ITAU_ENV", "1")
		os.Setenv("SANTANDER_ENV", "T")
		os.Setenv("URL_ITAU_REGISTER", "https://gerador-boletos.itau.com.br/router-gateway-app/public/codigo_barras/registro")
		os.Setenv("URL_ITAU_TICKET", "https://oauth.itau.com.br/identity/connect/token")
		os.Setenv("URL_BRADESCO_NET_EMPRESA", "https://cobranca.bradesconetempresa.b.br/ibpjregistrotitulows/registrotitulohomologacao")
		os.Setenv("RECOVERYROBOT_EXECUTION_ENABLED", "true")
		os.Setenv("RECOVERYROBOT_EXECUTION_IN_MINUTES", "2")
		os.Setenv("TIMEOUT_REGISTER", "30")
		os.Setenv("TIMEOUT_TOKEN", "20")
		os.Setenv("TIMEOUT_DEFAULT", "50")
		os.Setenv("URL_PEFISA_TOKEN", "https://psdo-hom.pernambucanas.com.br:444/sdcobr/api/oauth/token")
		os.Setenv("URL_PEFISA_REGISTER", "https://psdo-hom.pernambucanas.com.br:444/sdcobr/api/v2/titulos")
		os.Setenv("ENABLE_METRICS", "true")
	}
	config.Install(mockMode, devMode, disableLog)
}

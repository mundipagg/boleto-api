package env

import (
	"os"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

func Config(devMode, mockMode, disableLog bool) {
	configFlags(devMode, mockMode, disableLog)
	flow.RegisterConnector("logseq", util.SeqLogConector)
	flow.RegisterConnector("apierro", models.BoletoErrorConector)
	flow.RegisterConnector("tls", util.TlsConector)
	metrics.Install()
}

func ConfigMock(port string) {
	os.Setenv("URL_BB_REGISTER_BOLETO", "http://localhost:"+port+"/registrarBoleto")
	os.Setenv("URL_BB_TOKEN", "http://localhost:"+port+"/oauth/token")
	os.Setenv("URL_CAIXA", "http://localhost:"+port+"/caixa/registrarBoleto")
	os.Setenv("URL_CITI", "http://localhost:"+port+"/citi/registrarBoleto")
	os.Setenv("URL_SANTANDER_TICKET", "tls://localhost:"+port+"/santander/get-ticket")
	os.Setenv("URL_SANTANDER_REGISTER", "tls://localhost:"+port+"/santander/register")
	os.Setenv("URL_BRADESCO", "http://localhost:"+port+"/bradesco/registrarBoleto")
	os.Setenv("URL_ITAU_TOKEN", "http://localhost:"+port+"/itau/gerarToken")
	os.Setenv("URL_ITAU_REGISTER", "http://localhost:"+port+"/itau/registrarBoleto")
	config.Install(true, true, true)
}

func configFlags(devMode, mockMode, disableLog bool) {
	if devMode {
		os.Setenv("INFLUXDB_HOST", "http://localhost")
		os.Setenv("INFLUXDB_PORT", "8086")
		os.Setenv("API_PORT", "3000")
		os.Setenv("API_VERSION", "0.0.1")
		os.Setenv("ENVIROMENT", "Development")
		os.Setenv("SEQ_URL", "http://localhost:5341")   // http://stglog.mundipagg.com/ 192.168.8.119:5341
		os.Setenv("SEQ_API_KEY", "4jZzTybZ9bUHtJiPdh6") //4jZzTybZ9bUHtJiPdh6
		os.Setenv("ENABLE_REQUEST_LOG", "false")
		os.Setenv("ENABLE_PRINT_REQUEST", "true")
		os.Setenv("URL_BB_REGISTER_BOLETO", "https://cobranca.homologa.bb.com.br:7101/registrarBoleto")
		os.Setenv("URL_BB_TOKEN", "https://oauth.hm.bb.com.br:43000/oauth/token")
		os.Setenv("URL_CAIXA", "https://des.barramento.caixa.gov.br/sibar/ManutencaoCobrancaBancaria/Boleto/Externo")
		os.Setenv("URL_CITI", "https://citigroupsoauat.citigroup.com/comercioeletronico/registerboleto/RegisterBoletoSOAP")
		os.Setenv("URL_CITI_BOLETO", "https://ebillpayer.uat.brazil.citigroup.com/ebillpayer/jspInformaDadosConsulta.jsp")
		os.Setenv("APP_URL", "http://localhost:3000/boleto")
		os.Setenv("ELASTIC_URL", "http://localhost:9200")
		os.Setenv("MONGODB_URL", "localhost:27017")
		os.Setenv("MONGODB_USER", "")
		os.Setenv("MONGODB_PASSWORD", "")
		os.Setenv("BOLETO_JSON_STORE", "/home/philippe/boletodb/upMongo")
		os.Setenv("CERT_BOLETO_CRT", "C:\\cert_boleto_api\\certificate.crt")
		os.Setenv("CERT_BOLETO_KEY", "C:\\cert_boleto_api\\mundi.key")
		os.Setenv("CERT_BOLETO_CA", "C:\\cert_boleto_api\\ca-cert.ca")
		os.Setenv("URL_SANTANDER_TICKET", "https://ymbdlb.santander.com.br/dl-ticket-services/TicketEndpointService")
		os.Setenv("URL_SANTANDER_REGISTER", "https://ymbcash.santander.com.br/ymbsrv/CobrancaEndpointService")
		os.Setenv("URL_BRADESCO", "https://homolog.meiosdepagamentobradesco.com.br/api/transacao")
		os.Setenv("URL_ITAU_REGISTER", "https://gerador-boletos.itau.com.br/router-gateway-app/public/codigo_barras/registro")
		os.Setenv("URL_ITAU_TOKEN", "https://oauth.itau.com.br/identity/connect/token")
	}
	config.Install(mockMode, devMode, disableLog)
}

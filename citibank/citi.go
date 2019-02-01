package citibank

import (
	"net/http"
	"strconv"

	. "github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

type bankCiti struct {
	validate  *models.Validator
	log       *log.Log
	transport *http.Transport
}

func New() bankCiti {
	b := bankCiti{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}
	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(citiValidateAgency)
	b.validate.Push(citiValidateAccount)
	b.validate.Push(citiValidateWallet)
	transp, err := util.BuildTLSTransport(config.Get().CertBoletoPathCrt, config.Get().CertBoletoPathKey, config.Get().CertBoletoPathCa)
	if err != nil {
		//TODO
	}
	b.transport = transp
	return b
}

//Log retorna a referencia do log
func (b bankCiti) Log() *log.Log {
	return b.log
}

func (b bankCiti) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	timing := metrics.GetTimingMetrics()
	boleto.Title.OurNumber = calculateOurNumber(boleto)
	serviceURL := config.Get().URLCiti
	
	bod := NewFlow().From("message://?source=inline", boleto, getRequestCiti(), tmpl.GetFuncMaps())
	b.log.RequestCustom(bod.GetBody().(string), bod.GetHeader(), map[string]string{"URL" : serviceURL})

	var response string
	var status int
	var err error

	duration := util.Duration(func() {
		response, status, err = b.sendRequest(bod.GetBody().(string))
	})

	timing.Push("citibank-register-boleto-online", duration.Seconds())
	b.log.ResponseCustom(response, map[string]string{"ContentStatusCode": strconv.Itoa(status), "URL" : serviceURL, "Duration" : util.ConvertDuration(duration)})

	if err != nil {
		return models.BoletoResponse{}, err
	}

	var result interface{}

	if status == 200 {
		result = tmpl.TransformFromXML(response, getResponseCiti(), getAPIResponseCiti(), new(models.BoletoResponse))
	}

	switch t := result.(type) {
	case *models.BoletoResponse:
		return *t, nil
	case error:
		return models.BoletoResponse{}, t
	default:
		return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
	}
}

func (b bankCiti) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	return b.RegisterBoleto(boleto)
}

func (b bankCiti) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

func (b bankCiti) sendRequest(body string) (string, int, error) {
	serviceURL := config.Get().URLCiti
	if config.Get().MockMode {
		return util.Post(serviceURL, body, config.Get().TimeoutDefault, map[string]string{"Soapaction": "RegisterBoleto"})
	} else {
		return util.PostTLS(serviceURL, body, config.Get().TimeoutDefault, map[string]string{"Soapaction": "RegisterBoleto"}, b.transport)
	}
}

//GetBankNumber retorna o codigo do banco
func (b bankCiti) GetBankNumber() models.BankNumber {
	return models.Citibank
}

func calculateOurNumber(boleto *models.BoletoRequest) uint {
	ourNumberWithDigit := strconv.Itoa(int(boleto.Title.OurNumber)) + util.OurNumberDv(strconv.Itoa(int(boleto.Title.OurNumber)), util.MOD11)
	value, _ := strconv.Atoi(ourNumberWithDigit)
	return uint(value)
}

func (b bankCiti) GetBankNameIntegration() string {
	return "Citibank"
}

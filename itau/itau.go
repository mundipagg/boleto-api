package itau

import (
	"strconv"
	"errors"
	"strings"

	. "github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

type bankItau struct {
	validate *models.Validator
	log      *log.Log
}

func New() bankItau {
	b := bankItau{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}
	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(itauValidateAccount)
	b.validate.Push(itauValidateAgency)
	return b
}

//Log retorna a referencia do log
func (b bankItau) Log() *log.Log {
	return b.log
}

func (b bankItau) GetTicket(boleto *models.BoletoRequest) (string, error) {

	timing := metrics.GetTimingMetrics()
	pipe := NewFlow()
	url := config.Get().URLTicketItau
	
	pipe.From("message://?source=inline", boleto, getRequestTicket(), tmpl.GetFuncMaps())
	b.log.RequestCustom(pipe.GetBody().(string), pipe.GetHeader(), map[string]string{"URL" : url, "Operation":"GenerateToken"})
	
	var response string
	var status int
	var err error

	duration := util.Duration(func() {
		response, status, err = util.Post(url, pipe.GetBody().(string), config.Get().TimeoutToken, pipe.GetHeader())
	})
	
	timing.Push("itau-get-ticket-boleto-time", duration.Seconds())
	b.log.ResponseCustom(response, map[string]string{"ContentStatusCode": strconv.Itoa(status), "URL" : url, "Duration" : util.ConvertDuration(duration), "Operation" : "GenerateToken"})

	if err != nil {
		return "", err
	}

	var result interface{}

	if status == 200 {
		result = tmpl.TransformFromJSON(response, getTicketResponse(), `{{.access_token}}`, nil)
	}else if status == 400{
		errResult := tmpl.TransformFromJSON(response, getTicketResponse(), `{{.errorMessage}}`, nil)
		result = errors.New(errResult.(string))
	}else if status == 403{
		result = errors.New("403 Forbidden")
	}else if status == 500{
		errResult := tmpl.TransformFromJSON(response, getTicketErrorResponse(), `{{.errorMessage}}`, nil)
		result = errors.New(errResult.(string))
	}

	switch t := result.(type) {
	case string:
		return t, nil
	case error:
		return "", t
	default:
		return "", errors.New("Integration error")
	}
}

func (b bankItau) RegisterBoleto(input *models.BoletoRequest) (models.BoletoResponse, error) {
	timing := metrics.GetTimingMetrics()
	itauURL := config.Get().URLRegisterBoletoItau

	exec := NewFlow().From("message://?source=inline", input, getRequestItau(), tmpl.GetFuncMaps())
	b.log.RequestCustom(exec.GetBody().(string), exec.GetHeader(), map[string]string{"URL" : itauURL})

	var response string
	var status int
	var err error

	duration := util.Duration(func() {
		response, status, err = util.Post(itauURL, exec.GetBody().(string), config.Get().TimeoutRegister, exec.GetHeader())
	})

	timing.Push("itau-register-boleto-time", duration.Seconds())
	b.log.ResponseCustom(response, map[string]string{"ContentStatusCode": strconv.Itoa(status), "URL" : itauURL, "Duration" : util.ConvertDuration(duration)})

	if err != nil {
		return models.BoletoResponse{}, err
	}

	var result interface{}
	
	if status == 200 {
		result = tmpl.TransformFromJSON(response, getResponseItau(), getAPIResponseItau(), new(models.BoletoResponse))
	} else if strings.Contains(response, "text/html"){
		result = models.NewHTTPNotFound("404", "Page not found")
	} else if status == 400 {
		errResult := tmpl.TransformFromJSON(response, getResponseErrorItau(), `{{.errorMessage}}`, nil)
		result = models.NewErrorResponse("400", errResult.(string));
	} else {
		result = tmpl.TransformFromJSON(response, getResponseErrorItau(), getAPIResponseItau(), new(models.BoletoResponse))
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

func (b bankItau) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	if ticket, err := b.GetTicket(boleto); err != nil {
		return models.BoletoResponse{Errors: errs}, err
	} else {
		boleto.Authentication.AuthorizationToken = ticket
	}
	return b.RegisterBoleto(boleto)
}

func (b bankItau) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

//GetBankNumber retorna o codigo do banco
func (b bankItau) GetBankNumber() models.BankNumber {
	return models.Itau
}

func (b bankItau) GetBankNameIntegration() string {
	return "Itau"
}

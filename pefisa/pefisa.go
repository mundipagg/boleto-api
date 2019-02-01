package pefisa

import (
	"errors"
	"strconv"
	s "strings"
	. "github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

type bankPefisa struct {
	validate *models.Validator
	log      *log.Log
}

func New() bankPefisa {
	b := bankPefisa{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)

	return b
}

func (b bankPefisa) Log() *log.Log {
	return b.log
}

func (b bankPefisa) GetToken(boleto *models.BoletoRequest) (string, error) {

	timing := metrics.GetTimingMetrics()
	pipe := NewFlow()
	url := config.Get().URLPefisaToken

	pipe.From("message://?source=inline", boleto, getRequestToken(), tmpl.GetFuncMaps())
	b.log.RequestCustom(pipe.GetBody().(string), pipe.GetHeader(), map[string]string{"URL" : url, "Operation":"GenerateToken"})
	
	var response string
	var status int
	var err error

	duration := util.Duration(func() {
		response, status, err = util.Post(url, pipe.GetBody().(string), config.Get().TimeoutToken, pipe.GetHeader())
	})

	timing.Push("pefisa-get-token-boleto-time", duration.Seconds())
	b.log.ResponseCustom(response, map[string]string{"ContentStatusCode": strconv.Itoa(status), "URL" : url, "Duration" : util.ConvertDuration(duration), "Operation" : "GenerateToken"})

	if err != nil {
		return "", err
	}

	var result interface{}

	if status == 200 {
		result = tmpl.TransformFromJSON(response, getTokenResponse(), `{{.access_token}}`, nil)
	}else if status == 401 {
		errResult := tmpl.TransformFromJSON(response, getTokenErrorResponse(), `{{.errorMessage}}`, nil)
		result = models.NewErrorResponse("401", errResult.(string));
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

func (b bankPefisa) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	timing := metrics.GetTimingMetrics()
	pefisaURL := config.Get().URLPefisaRegister
	
	exec := NewFlow().From("message://?source=inline", boleto, getRequestPefisa(), tmpl.GetFuncMaps())
	b.log.RequestCustom(exec.GetBody().(string), nil, map[string]string{"URL" : pefisaURL})

	var response string
	var status int
	var err error

	duration := util.Duration(func() {
		headers := map[string]string{"Authorization": "Bearer " + boleto.Authentication.AuthorizationToken, "Content-Type": "application/json"}
		response, status, err = util.Post(pefisaURL, exec.GetBody().(string), config.Get().TimeoutRegister, headers)
	})

	timing.Push("pefisa-register-boleto-time", duration.Seconds())
	b.log.ResponseCustom(response, map[string]string{"ContentStatusCode": strconv.Itoa(status), "URL" : pefisaURL, "Duration" : util.ConvertDuration(duration)})

	if err != nil {
		return models.BoletoResponse{}, err
	}

	var result interface{}

	if status == 200 {
		result = tmpl.TransformFromJSON(response, getResponsePefisa(), getAPIResponsePefisa(), new(models.BoletoResponse))
	}else if status == 401 {
		result = tmpl.TransformFromJSON(response, getResponseErrorPefisa(), getAPIResponsePefisa(), new(models.BoletoResponse))
	}else if status == 400 {
		dataError := util.ParseJSON(response, new(models.ArrayDataError)).(*models.ArrayDataError)
		newBody := s.Replace(util.Stringify(dataError.Error[0]), "\\\"", "", -1)

		result = tmpl.TransformFromJSON(newBody, getResponseErrorPefisaArray(), getAPIResponsePefisa(), new(models.BoletoResponse))
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

func (b bankPefisa) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)

	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	if token, err := b.GetToken(boleto); err != nil {
		return models.BoletoResponse{Errors: errs}, err
	} else {
		boleto.Authentication.AuthorizationToken = token
	}

	return b.RegisterBoleto(boleto)
}

func (b bankPefisa) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

//GetBankNumber retorna o codigo do banco
func (b bankPefisa) GetBankNumber() models.BankNumber {
	return models.Pefisa
}

func (b bankPefisa) GetBankNameIntegration() string {
	return "Pefisa"
}


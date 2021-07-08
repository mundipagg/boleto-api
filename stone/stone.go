package stone

import (
	"fmt"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

type bankStone struct {
	validate *models.Validator
	log      *log.Log
}

//New Create a new Stone Integration Instance
func New() bankStone {
	b := bankStone{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)

	b.validate.Push(stoneValidateAccessKeyNotEmpty)

	return b
}

func (b bankStone) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}

	if accToken, err := authenticate(boleto.Authentication.AccessKey, b.log); err != nil {
		return models.BoletoResponse{Errors: errs}, err
	} else {
		boleto.Authentication.AuthorizationToken = accToken
	}
	return b.RegisterBoleto(boleto)
}

func (b bankStone) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	var response string
	var status int
	var err error

	stoneURL := config.Get().URLStoneRegister
	boleto.Title.BoletoType, boleto.Title.BoletoTypeCode = getBoletoType(boleto)

	body := flow.NewFlow().From("message://?source=inline", boleto, templateRequest, tmpl.GetFuncMaps()).GetBody().(string)
	head := hearders(boleto.Authentication.AuthorizationToken)
	b.log.Request(body, stoneURL, head)

	duration := util.Duration(func() {
		response, status, err = util.Post(stoneURL, body, config.Get().TimeoutRegister, head)
	})
	metrics.PushTimingMetric("stone-register-boleto-time", duration.Seconds())

	b.log.Response(response, stoneURL)

	return mapStoneResponse(boleto, response, status, err), nil
}

func mapStoneResponse(request *models.BoletoRequest, response string, status int, httpErr error) models.BoletoResponse {
	f := flow.NewFlow().To("set://?prop=body", response)
	switch status {
	case 0, 504:
		var msg string
		if httpErr != nil {
			msg = fmt.Sprintf("%v", httpErr)
		} else {
			msg = "GatewayTimeout"
		}
		return models.GetBoletoResponseError("MPTimeout", msg)
	case 201:
		f.To("transform://?format=json", templateResponse, templateAPI, tmpl.GetFuncMaps())
		f.To("unmarshall://?format=json", new(models.BoletoResponse))
	default:
		f.To("transform://?format=json", templateError, templateAPI, tmpl.GetFuncMaps())
		f.To("unmarshall://?format=json", new(models.BoletoResponse))
	}

	switch t := f.GetBody().(type) {
	case *models.BoletoResponse:
		if hasOurNumberFail(t) {
			return models.GetBoletoResponseError("MPOurNumberFail", "our number was not returned by the bank")
		} else {
			return *t
		}
	case error:
		return models.GetBoletoResponseError("MP500", t.Error())
	case string:
		return models.GetBoletoResponseError("MP500", t)
	}

	return models.GetBoletoResponseError("MP500", "Internal Error")
}

func (b bankStone) ValidateBoleto(request *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(request))
}

func (b bankStone) GetBankNumber() models.BankNumber {
	return models.Stone
}

func (b bankStone) GetBankNameIntegration() string {
	return "Stone"
}

func (b bankStone) Log() *log.Log {
	return b.log
}

func getBoletoType(boleto *models.BoletoRequest) (bt string, btc string) {
	return "DM", "bill_of_exchange"
}

func hearders(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token, "Content-Type": "application/json"}
}

func hasOurNumberFail(response *models.BoletoResponse) bool {
	return !response.HasErrors() && response.OurNumber == ""
}

package itau

import (
	"errors"

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
	return b
}

//Log retorna a referencia do log
func (b bankItau) Log() *log.Log {
	return b.log
}

func (b bankItau) GetTicket(boleto *models.BoletoRequest) (string, error) {
	timing := metrics.GetTimingMetrics()
	pipe := NewFlow()
	url := config.Get().URLItauToken
	pipe.From("message://?source=inline", boleto, getRequestTicket(), tmpl.GetFuncMaps())
	pipe.To("logseq://?type=request&url="+url, b.log)
	duration := util.Duration(func() {
		pipe.To(url, map[string]string{"method": "POST", "insecureSkipVerify": "true"})
	})
	timing.Push("itau-get-ticket-boleto-time", duration.Seconds())
	pipe.To("logseq://?type=response&url="+url, b.log)
	ch := pipe.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=json", getTicketResponse(), `{{.access_token}}`, tmpl.GetFuncMaps())
	ch.When(Header("status").IsEqualTo("400"))
	ch.To("transform://?format=json", getTicketResponse(), `{{.errorMessage}}`, tmpl.GetFuncMaps())
	ch.To("set://?prop=body", errors.New(pipe.GetBody().(string)))
	ch.When(Header("status").IsEqualTo("403"))
	ch.To("set://?prop=body", errors.New("403 Forbidden"))
	ch.When(Header("status").IsEqualTo("500"))
	ch.To("transform://?format=json", getTicketErrorResponse(), `{{.errorMessage}}`, tmpl.GetFuncMaps())
	ch.To("set://?prop=body", errors.New(pipe.GetBody().(string)))
	ch.Otherwise()
	ch.To("logseq://?type=request&url="+url, b.log).To("print://?msg=${body}").To("set://?prop=body", errors.New("integration error"))
	switch t := pipe.GetBody().(type) {
	case string:
		return t, nil
	case error:
		return "", t
	}
	return "", nil
}

func (b bankItau) RegisterBoleto(input *models.BoletoRequest) (models.BoletoResponse, error) {
	timing := metrics.GetTimingMetrics()
	itauURL := config.Get().URLRegisterBoletoItau
	fromResponse := getResponseItau()
	fromResponseError := getResponseErrorItau()
	toAPI := getAPIResponseItau()
	inputTemplate := getRequestItau()
	exec := NewFlow().From("message://?source=inline", input, inputTemplate, tmpl.GetFuncMaps())
	exec.To("logseq://?type=request&url="+itauURL, b.log)
	duration := util.Duration(func() {
		exec.To(itauURL, map[string]string{"method": "POST", "insecureSkipVerify": "true"})
	})
	timing.Push("itau-register-boleto-time", duration.Seconds())
	exec.To("logseq://?type=response&url="+itauURL, b.log)
	ch := exec.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=json", fromResponse, toAPI, tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))
	ch.When(Header("status").IsEqualTo("401"))
	ch.To("transform://?format=json", fromResponse, toAPI, tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))
	ch.When(Header("status").IsEqualTo("400"))
	ch.To("transform://?format=json", fromResponseError, toAPI, tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))
	ch.Otherwise()
	ch.To("logseq://?type=response&url="+itauURL, b.log).To("apierro://")
	switch t := exec.GetBody().(type) {
	case *models.BoletoResponse:
		return *t, nil
	case error:
		return models.BoletoResponse{}, t
	}
	return models.BoletoResponse{}, models.NewInternalServerError("Internal error", "MP500")
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

package pfisa

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

type bankPfisa struct {
	validate *models.Validator
	log      *log.Log
}

func New() bankPfisa {
	b := bankPfisa{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)

	return b
}

func (b bankPfisa) Log() *log.Log {
	return b.log
}

func (b bankPfisa) GetToken(boleto *models.BoletoRequest) (string, error) {
	
	timing := metrics.GetTimingMetrics()
	pipe := NewFlow()
	url := config.Get().URLPfisaToken
	pipe.From("message://?source=inline", boleto, getRequestToken(), tmpl.GetFuncMaps())
	pipe.To("logseq://?type=request&url="+url, b.log)
	duration := util.Duration(func() {
		pipe.To(url, map[string]string{"method": "POST", "insecureSkipVerify": "true", "timeout": config.Get().TimeoutToken})
	})
	timing.Push("itau-get-tocket-boleto-time", duration.Seconds())
	pipe.To("logseq://?type=response&url="+url, b.log)
	ch := pipe.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=json", getTokenResponse(), `{{.access_token}}`, tmpl.GetFuncMaps())
	ch.When(Header("status").IsEqualTo("401"))
	ch.To("transform://?format=json", getTokenResponse(), `{{.error_description}}`, tmpl.GetFuncMaps())
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

func (b bankPfisa) RegisterBoleto(input *models.BoletoRequest) (models.BoletoResponse, error) {
	timing := metrics.GetTimingMetrics()
	pfisaURL := config.Get().URLPfisaRegister
	fromResponse := getResponsePfisa()
	toAPI := getAPIResponsePfisa()
	inputTemplate := getRequestPfisa()
	exec := NewFlow().From("message://?source=inline", input, inputTemplate, tmpl.GetFuncMaps())
	exec.To("logseq://?type=request&url="+pfisaURL, b.log)
	duration := util.Duration(func() {
		exec.To(pfisaURL, map[string]string{"method": "POST", "insecureSkipVerify": "true", "timeout": config.Get().TimeoutRegister})
	})
	timing.Push("pfisa-register-boleto-time", duration.Seconds())
	exec.To("logseq://?type=response&url="+pfisaURL, b.log)

	ch := exec.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=json", fromResponse, toAPI, tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))

	ch.Otherwise()
	ch.To("logseq://?type=response&url="+pfisaURL, b.log).To("apierro://")

	switch t := exec.GetBody().(type) {
	case *models.BoletoResponse:
		return *t, nil
	case error:
		return models.BoletoResponse{}, t
	}
	return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
}

func (b bankPfisa) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
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

func (b bankPfisa) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

//GetBankNumber retorna o codigo do banco
func (b bankPfisa) GetBankNumber() models.BankNumber {
	return models.Pfisa
}

func (b bankPfisa) GetBankNameIntegration() string {
	return "Pfisa"
}
package bb

import (
	"errors"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"

	"github.com/mundipagg/boleto-api/validations"
)

type bankBB struct {
	validate *models.Validator
	log      *log.Log
}

//Cria uma nova instância do objeto que implementa os serviços do Banco do Brasil e configura os validadores que serão utilizados
func New() bankBB {
	b := bankBB{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}
	b.validate.Push(bbValidateAccountAndDigit)
	b.validate.Push(bbValidateAgencyAndDigit)
	b.validate.Push(bbValidateOurNumber)
	b.validate.Push(bbValidateWalletVariation)
	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(bbValidateTitleInstructions)
	b.validate.Push(bbValidateTitleDocumentNumber)
	return b
}

//Log retorna a referencia do log
func (b bankBB) Log() *log.Log {
	return b.log
}

func (b *bankBB) login(boleto *models.BoletoRequest) (string, error) {
	type errorAuth struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	timing := metrics.GetTimingMetrics()
	r := flow.NewFlow()
	url := config.Get().URLBBToken
	from, resp := GetBBAuthLetters()
	bod := r.From("message://?source=inline", boleto, from, tmpl.GetFuncMaps())
	r = r.To("logseq://?type=request&url="+url, b.log)
	duration := util.Duration(func() {
		bod = bod.To(url, map[string]string{"method": "POST", "insecureSkipVerify": "true"})
	})
	timing.Push("bb-login-time", duration.Seconds())
	r = r.To("logseq://?type=response&url="+url, b.log)
	ch := bod.Choice().When(flow.Header("status").IsEqualTo("200")).To("transform://?format=json", resp, `{{.authToken}}`)
	ch = ch.Otherwise().To("unmarshall://?format=json", new(errorAuth))
	result := bod.GetBody()
	switch t := result.(type) {
	case string:
		return t, nil
	case error:
		return "", t
	case *errorAuth:
		return "", errors.New(t.ErrorDescription)
	}
	return "", errors.New("Saída inválida")
}

//ProcessBoleto faz o processamento de registro de boleto
func (b bankBB) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	tok, err := b.login(boleto)
	if err != nil {
		return models.BoletoResponse{}, err
	}
	boleto.Authentication.AuthorizationToken = tok
	return b.RegisterBoleto(boleto)
}

func (b bankBB) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	r := flow.NewFlow()
	url := config.Get().URLBBRegisterBoleto
	from := getRequest()
	timing := metrics.GetTimingMetrics()
	r = r.From("message://?source=inline", boleto, from, tmpl.GetFuncMaps())
	r.To("logseq://?type=request&url="+url, b.log)
	duration := util.Duration(func() {
		r.To(url, map[string]string{"method": "POST", "insecureSkipVerify": "true"})
	})
	timing.Push("bb-register-boleto-time", duration.Seconds())
	r.To("logseq://?type=response&url="+url, b.log)
	ch := r.Choice()
	ch.When(flow.Header("status").IsEqualTo("200"))
	ch.To("transform://?format=xml", getResponseBB(), getAPIResponse(), tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))
	ch.Otherwise()
	ch.To("logseq://?type=response&url="+url, b.log).To("apierro://")
	switch t := r.GetBody().(type) {
	case *models.BoletoResponse:
		return *t, nil
	case models.BoletoResponse:
		return t, nil
	case error:
		return models.BoletoResponse{}, t
	default:
		return models.BoletoResponse{}, errors.New("Unexpected Type")
	}

}

func (b bankBB) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

//GetBankNumber retorna o codigo do banco
func (b bankBB) GetBankNumber() models.BankNumber {
	return models.BancoDoBrasil
}

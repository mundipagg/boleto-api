package pefisa

import (
	"errors"
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
	pipe.To("logseq://?type=request&url="+url, b.log)

	duration := util.Duration(func() {
		pipe.To(url, map[string]string{"method": "POST", "insecureSkipVerify": "true", "timeout": config.Get().TimeoutToken})
	})
	timing.Push("pefisa-get-token-boleto-time", duration.Seconds())
	pipe.To("logseq://?type=response&url="+url, b.log)
	ch := pipe.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=json", getTokenResponse(), `{{.access_token}}`, tmpl.GetFuncMaps())

	ch.When(Header("status").IsEqualTo("401"))
	ch.To("transform://?format=json", getTokenErrorResponse(), `{{.error_description}}`, tmpl.GetFuncMaps())
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

func (b bankPefisa) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	timing := metrics.GetTimingMetrics()
	pefisaURL := config.Get().URLPefisaRegister

	exec := NewFlow().From("message://?source=inline", boleto, getRequestPefisa(), tmpl.GetFuncMaps())
	exec.To("logseq://?type=request&url="+pefisaURL, b.log)

	var response string
	var status int
	var err error
	duration := util.Duration(func() {
		response, status, err = b.sendRequest(exec.GetBody().(string), boleto.Authentication.AuthorizationToken)
	})
	if err != nil {
		return models.BoletoResponse{}, err
	}

	timing.Push("pefisa-register-boleto-time", duration.Seconds())
	exec.To("set://?prop=header", map[string]string{"status": strconv.Itoa(status)})
	exec.To("set://?prop=body", response)
	exec.To("logseq://?type=response&url="+pefisaURL, b.log)

	if status == 200 || status == 401 {
		exec.To("set://?prop=body", response)
	} else {
		dataError := util.ParseJSON(response, new(models.ArrayDataError)).(*models.ArrayDataError)
		exec.To("set://?prop=body", util.Stringify(dataError.Error[0]))
	}

	ch := exec.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=json", getResponsePefisa(), getAPIResponsePefisa(), tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))

	ch.When(Header("status").IsEqualTo("400"))
	ch.To("transform://?format=json", getResponseErrorPefisaArray(), getAPIResponsePefisa(), tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))

	ch.When(Header("status").IsEqualTo("401"))
	ch.To("transform://?format=json", getResponseErrorPefisa(), getAPIResponsePefisa(), tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))

	ch.Otherwise()
	ch.To("logseq://?type=response&url="+pefisaURL, b.log).To("apierro://")

	switch t := exec.GetBody().(type) {
	case *models.BoletoResponse:
		return *t, nil
	case error:
		return models.BoletoResponse{}, t
	}
	return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
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

func (b bankPefisa) sendRequest(body string, token string) (string, int, error) {
	serviceURL := config.Get().URLPefisaRegister

	h := map[string]string{"Authorization": "Bearer " + token, "Content-Type": "application/json"}
	return util.Post(serviceURL, body, config.Get().TimeoutRegister, h)
}

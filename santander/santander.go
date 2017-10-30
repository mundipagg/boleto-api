package santander

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	. "github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

type bankSantander struct {
	validate  *models.Validator
	log       *log.Log
	transport *http.Transport
}

//New Create a new Santander Integration Instance
func New() bankSantander {
	b := bankSantander{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}
	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(santanderValidateAgreementNumber)

	t, err := util.BuildTLSTransport(config.Get().CertBoletoPathCrt, config.Get().CertBoletoPathKey, config.Get().CertBoletoPathCa)
	if err != nil {
		fmt.Println(err.Error())
	}
	b.transport = t

	return b
}

//Log retorna a referencia do log
func (b bankSantander) Log() *log.Log {
	return b.log
}
func (b bankSantander) GetTicket(boleto *models.BoletoRequest) (string, error) {
	timing := metrics.GetTimingMetrics()
	boleto.Title.OurNumber = calculateOurNumber(boleto)
	pipe := NewFlow()
	url := config.Get().URLTicketSantander
	tlsURL := strings.Replace(config.Get().URLTicketSantander, "https", "tls", 1)
	pipe.From("message://?source=inline", boleto, getRequestTicket(), tmpl.GetFuncMaps())
	pipe.To("logseq://?type=request&url="+url, b.log)
	duration := util.Duration(func() {
		pipe.To(tlsURL, b.transport)
	})
	timing.Push("santander-get-ticket-boleto-time", duration.Seconds())
	pipe.To("logseq://?type=response&url="+url, b.log)
	ch := pipe.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=xml", getTicketResponse(), `{{.returnCode}}:::{{.ticket}}`, tmpl.GetFuncMaps())
	ch.When(Header("status").IsEqualTo("403"))
	ch.To("set://?prop=body", errors.New("403 Forbidden"))
	ch.Otherwise()
	ch.To("logseq://?type=request&url="+url, b.log).To("set://?prop=body", errors.New("integration error"))
	switch t := pipe.GetBody().(type) {
	case string:
		items := pipe.GetBody().(string)
		parts := strings.Split(items, ":::")
		returnCode, ticket := parts[0], parts[1]
		return ticket, checkError(returnCode)
	case error:
		return "", t
	}
	return "", nil
}

func (b bankSantander) RegisterBoleto(input *models.BoletoRequest) (models.BoletoResponse, error) {
	timing := metrics.GetTimingMetrics()
	serviceURL := config.Get().URLRegisterBoletoSantander
	fromResponse := getResponseSantander()
	toAPI := getAPIResponseSantander()
	inputTemplate := getRequestSantander()
	santanderURL := strings.Replace(serviceURL, "https", "tls", 1)

	exec := NewFlow().From("message://?source=inline", input, inputTemplate, tmpl.GetFuncMaps())
	exec.To("logseq://?type=request&url="+serviceURL, b.log)
	duration := util.Duration(func() {
		exec.To(santanderURL, b.transport, map[string]string{"method": "POST", "insecureSkipVerify": "true"})
	})
	timing.Push("santander-register-boleto-time", duration.Seconds())
	exec.To("logseq://?type=response&url="+serviceURL, b.log)
	ch := exec.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=xml", fromResponse, toAPI, tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))
	ch.Otherwise()
	ch.To("logseq://?type=response&url="+serviceURL, b.log).To("apierro://")
	switch t := exec.GetBody().(type) {
	case *models.BoletoResponse:
		return *t, nil
	case error:
		return models.BoletoResponse{}, t
	}
	return models.BoletoResponse{}, models.NewInternalServerError("Internal error", "MP500")
}
func (b bankSantander) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
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

func (b bankSantander) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

//GetBankNumber retorna o codigo do banco
func (b bankSantander) GetBankNumber() models.BankNumber {
	return models.Santander
}

func calculateOurNumber(boleto *models.BoletoRequest) uint {
	ourNumberWithDigit := strconv.Itoa(int(boleto.Title.OurNumber)) + util.OurNumberDv(strconv.Itoa(int(boleto.Title.OurNumber)))
	value, _ := strconv.Atoi(ourNumberWithDigit)
	return uint(value)
}

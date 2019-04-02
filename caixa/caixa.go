package caixa

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mundipagg/boleto-api/metrics"

	"github.com/PMoneda/flow"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

var o = &sync.Once{}
var m map[string]string

type bankCaixa struct {
	validate *models.Validator
	log      *log.Log
}

func New() bankCaixa {
	b := bankCaixa{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}
	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(caixaValidateAgency)
	b.validate.Push(validadeOurNumber)
	b.validate.Push(caixaValidateBoletoType)
	return b
}

//Log retorna a referencia do log
func (b bankCaixa) Log() *log.Log {
	return b.log
}
func (b bankCaixa) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {

	boleto.Title.BoletoType, boleto.Title.BoletoTypeCode = getBoletoType(boleto)

	r := flow.NewFlow()
	urlCaixa := config.Get().URLCaixaRegisterBoleto
	from := getResponseCaixa()
	to := getAPIResponseCaixa()

	bod := r.From("message://?source=inline", boleto, getRequestCaixa(), tmpl.GetFuncMaps())
	bod = bod.To("logseq://?type=request&url="+urlCaixa, b.log)
	duration := util.Duration(func() {
		bod = bod.To(urlCaixa, map[string]string{"method": "POST", "insecureSkipVerify": "true", "timeout": config.Get().TimeoutDefault})
	})
	metrics.PushTimingMetric("caixa-register-time", duration.Seconds())
	bod = bod.To("logseq://?type=response&url="+urlCaixa, b.log)
	ch := bod.Choice()
	ch = ch.When(flow.Header("status").IsEqualTo("200"))
	ch = ch.To("transform://?format=xml", from, to, tmpl.GetFuncMaps())
	ch = ch.Otherwise()
	ch = ch.To("logseq://?type=response&url="+urlCaixa, b.log).To("apierro://")

	switch t := bod.GetBody().(type) {
	case string:
		response := util.ParseJSON(t, new(models.BoletoResponse)).(*models.BoletoResponse)
		return *response, nil
	case models.BoletoResponse:
		return t, nil
	}
	return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
}
func (b bankCaixa) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}

	boleto.Title.OurNumber = b.FormatOurNumber(boleto.Title.OurNumber)

	checkSum := b.getCheckSumCode(*boleto)

	boleto.Authentication.AuthorizationToken = b.getAuthToken(checkSum)
	return b.RegisterBoleto(boleto)
}

func (b bankCaixa) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

func (b bankCaixa) FormatOurNumber(ourNumber uint) uint {

	if ourNumber != 0 {
		ourNumberFormatted := 14000000000000000 + ourNumber

		return ourNumberFormatted
	}

	return ourNumber
}

//getCheckSumCode Código do Cedente (7 posições) + Nosso Número (17 posições) + Data de Vencimento (DDMMAAAA) + Valor (15 posições) + CPF/CNPJ (14 Posições)
func (b bankCaixa) getCheckSumCode(boleto models.BoletoRequest) string {

	return fmt.Sprintf("%07d%017d%s%015d%014s",
		boleto.Agreement.AgreementNumber,
		boleto.Title.OurNumber,
		boleto.Title.ExpireDateTime.Format("02012006"),
		boleto.Title.AmountInCents,
		boleto.Recipient.Document.Number)
}

func (b bankCaixa) getAuthToken(info string) string {
	return util.Sha256(info, "base64")
}

//GetBankNumber retorna o codigo do banco
func (b bankCaixa) GetBankNumber() models.BankNumber {
	return models.Caixa
}

func (b bankCaixa) GetBankNameIntegration() string {
	return "Caixa"
}

func caixaBoletoTypes() map[string]string {
	o.Do(func() {
		m = make(map[string]string)

		m["OUT"] = "99" //Duplicata Mercantil p/ Indicação
	})
	return m
}

func getBoletoType(boleto *models.BoletoRequest) (bt string, btc string) {
	if len(boleto.Title.BoletoType) < 1 {
		return "OUT", "99"
	}
	btm := caixaBoletoTypes()

	if btm[strings.ToUpper(boleto.Title.BoletoType)] == "" {
		return "OUT", "99"
	}

	return boleto.Title.BoletoType, btm[strings.ToUpper(boleto.Title.BoletoType)]
}

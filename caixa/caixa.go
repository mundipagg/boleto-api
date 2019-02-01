package caixa

import (
	"strconv"
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
	return b
}

//Log retorna a referencia do log
func (b bankCaixa) Log() *log.Log {
	return b.log
}
func (b bankCaixa) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	
	timing := metrics.GetTimingMetrics()
	r := flow.NewFlow()
	urlCaixa := config.Get().URLCaixaRegisterBoleto

	bod := r.From("message://?source=inline", boleto, getRequestCaixa(), tmpl.GetFuncMaps())
	b.log.RequestCustom(bod.GetBody().(string), nil, map[string]string{"URL" : urlCaixa})

	var response string
	var status int
	var err error

	duration := util.Duration(func() {
		response, status, err = util.Post(urlCaixa, bod.GetBody().(string), config.Get().TimeoutRegister, bod.GetHeader())
	})

	timing.Push("caixa-register-time", duration.Seconds())
	b.log.ResponseCustom(response, map[string]string{"ContentStatusCode": strconv.Itoa(status), "URL" : urlCaixa, "Duration" : util.ConvertDuration(duration)})

	if err != nil {
		return models.BoletoResponse{}, err
	}
	
	var result interface{}

	if status == 200 {
		result = tmpl.TransformFromXML(response, getResponseCaixa(), getAPIResponseCaixa(), new(models.BoletoResponse))
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

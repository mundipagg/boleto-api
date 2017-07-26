package caixa

import (
	"fmt"

	"github.com/PMoneda/flow"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/letters"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
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
	b.validate.Push(baseValidateAmountInCents)
	b.validate.Push(baseValidateExpireDate)
	b.validate.Push(baseValidateBuyerDocumentNumber)
	b.validate.Push(baseValidateRecipientDocumentNumber)
	b.validate.Push(caixaValidateAgency)
	b.validate.Push(validateInstructions)
	return b
}

//Log retorna a referencia do log
func (b bankCaixa) Log() *log.Log {
	return b.log
}
func (b bankCaixa) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	r := flow.NewFlow()
	urlCaixa := config.Get().URLCaixaRegisterBoleto
	from := flow.Transform(letters.GetResponseTemplateCaixa())
	to := flow.Transform(letters.GetRegisterBoletoAPIResponseTmpl(models.Caixa))
	bod := r.From("message://?source=inline", boleto, letters.GetRegisterBoletoCaixaTmpl(), tmpl.GetFuncMaps())
	bod = bod.To("logseq://?type=request&url="+urlCaixa, b.log)
	//bod = bod.To("print://?msg=${body}")
	bod = bod.To(urlCaixa, map[string]string{"method": "POST", "insecureSkipVerify": "true"})
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
	return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Erro interno")
}
func (b bankCaixa) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	checkSum := b.getCheckSumCode(*boleto)
	//fmt.Println(checkSum)
	boleto.Authentication.AuthorizationToken = b.getAuthToken(checkSum)
	return b.RegisterBoleto(boleto)
}

func (b bankCaixa) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
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
	return util.Sha256(info)
}

//GetBankNumber retorna o codigo do banco
func (b bankCaixa) GetBankNumber() models.BankNumber {
	return models.Caixa
}

package bradescoShopFacil

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

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

type bankBradescoShopFacil struct {
	validate *models.Validator
	log      *log.Log
}

//barcode struct for bradescoShopFacil
type barcode struct {
	bankCode      string
	currencyCode  string
	dateDueFactor string
	value         string
	agency        string
	wallet        string
	ourNumber     string
	account       string
	zero          string
}

//New creates a new BradescoShopFacil instance
func New() bankBradescoShopFacil {
	b := bankBradescoShopFacil{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}
	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(bradescoShopFacilValidateAgency)
	b.validate.Push(bradescoShopFacilValidateAccount)
	b.validate.Push(bradescoShopFacilValidateWallet)
	b.validate.Push(bradescoShopFacilValidateAuth)
	b.validate.Push(bradescoShopFacilValidateAgreement)
	b.validate.Push(bradescoShopFacilBoletoTypeValidate)
	return b
}

//Log retorna a referencia do log
func (b bankBradescoShopFacil) Log() *log.Log {
	return b.log
}

func (b bankBradescoShopFacil) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	boleto.Title.BoletoType, boleto.Title.BoletoTypeCode = getBoletoType(boleto)
	r := flow.NewFlow()
	serviceURL := config.Get().URLBradescoShopFacil
	from := getResponseBradescoShopFacil()
	to := getAPIResponseBradescoShopFacil()
	bod := r.From("message://?source=inline", boleto, getRequestBradescoShopFacil(), tmpl.GetFuncMaps())
	bod.To("logseq://?type=request&url="+serviceURL, b.log)
	duration := util.Duration(func() {
		bod.To(serviceURL, map[string]string{"method": "POST", "insecureSkipVerify": "true", "timeout": config.Get().TimeoutDefault})
	})
	metrics.PushTimingMetric("bradesco-shopfacil-register-boleto-online", duration.Seconds())
	bod.To("logseq://?type=response&url="+serviceURL, b.log)
	ch := bod.Choice()
	ch.When(flow.Header("status").IsEqualTo("201"))
	ch.To("transform://?format=json", from, to, tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))
	ch.When(flow.Header("status").IsEqualTo("200"))
	ch.To("transform://?format=json", from, to, tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))
	ch.Otherwise()
	ch.To("logseq://?type=response&url="+serviceURL, b.log).To("apierro://")
	switch t := bod.GetBody().(type) {
	case *models.BoletoResponse:
		if !t.HasErrors() {
			t.BarCodeNumber = getBarcode(*boleto).toString()
		}
		return *t, nil
	case error:
		return models.BoletoResponse{}, t
	}
	return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
}

func (b bankBradescoShopFacil) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	return b.RegisterBoleto(boleto)
}

func (b bankBradescoShopFacil) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

func (b bankBradescoShopFacil) GetBankNumber() models.BankNumber {
	return models.Bradesco
}

func getBarcode(boleto models.BoletoRequest) (bc barcode) {
	bc.bankCode = fmt.Sprintf("%d", models.Bradesco)
	bc.currencyCode = fmt.Sprintf("%d", models.Real)
	bc.account = fmt.Sprintf("%07s", boleto.Agreement.Account)
	bc.agency = fmt.Sprintf("%04s", boleto.Agreement.Agency)
	bc.dateDueFactor, _ = dateDueFactor(boleto.Title.ExpireDateTime)
	bc.ourNumber = fmt.Sprintf("%011d", boleto.Title.OurNumber)
	bc.value = fmt.Sprintf("%010d", boleto.Title.AmountInCents)
	bc.wallet = fmt.Sprintf("%02d", boleto.Agreement.Wallet)
	bc.zero = "0"
	return
}

func (bc barcode) toString() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s", bc.bankCode, bc.currencyCode, bc.calcCheckDigit(), bc.dateDueFactor, bc.value, bc.agency, bc.wallet, bc.ourNumber, bc.account, bc.zero)
}

func (bc barcode) calcCheckDigit() string {
	prevCode := fmt.Sprintf("%s%s%s%s%s%s%s%s%s", bc.bankCode, bc.currencyCode, bc.dateDueFactor, bc.value, bc.agency, bc.wallet, bc.ourNumber, bc.account, bc.zero)
	return util.BarcodeDv(prevCode)
}

func dateDueFactor(dateDue time.Time) (string, error) {
	var dateDueFixed = time.Date(1997, 10, 7, 0, 0, 0, 0, time.UTC)
	dif := dateDue.Sub(dateDueFixed)
	factor := int(dif.Hours() / 24)
	if factor <= 0 {
		return "", errors.New("DateDue must be in the future")
	}
	return fmt.Sprintf("%04d", factor), nil
}

func (b bankBradescoShopFacil) GetBankNameIntegration() string {
	return "BradescoShopFacil"
}

func bradescoShopFacilBoletoTypes() map[string]string {

	o.Do(func() {
		m = make(map[string]string)

		m["DM"] = "01"  //Duplicata Mercantil
		m["NP"] = "02"  //Nota promissória
		m["RC"] = "05"  //Recibo
		m["DS"] = "12"  //Duplicata de serviço
		m["OUT"] = "99" //Outros
	})
	return m
}

func getBoletoType(boleto *models.BoletoRequest) (bt string, btc string) {
	if len(boleto.Title.BoletoType) < 1 {
		return "DM", "01"
	}
	btm := bradescoShopFacilBoletoTypes()

	if btm[strings.ToUpper(boleto.Title.BoletoType)] == "" {
		return "DM", "01"
	}

	return boleto.Title.BoletoType, btm[strings.ToUpper(boleto.Title.BoletoType)]

}

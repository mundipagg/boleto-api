package bradescoShopFacil

import (
	"strconv"
	"errors"
	"fmt"
	"time"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

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
	return b
}

//Log retorna a referencia do log
func (b bankBradescoShopFacil) Log() *log.Log {
	return b.log
}

func (b bankBradescoShopFacil) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	timing := metrics.GetTimingMetrics()
	serviceURL := config.Get().URLBradescoShopFacil

	bod := flow.NewFlow().From("message://?source=inline", boleto, getRequestBradescoShopFacil(), tmpl.GetFuncMaps())
	b.log.RequestCustom(bod.GetBody().(string), bod.GetHeader(), map[string]string{"URL" : serviceURL})

	var response string
	var status int
	var err error

	duration := util.Duration(func() {
		response, status, err = util.Post(serviceURL, bod.GetBody().(string), config.Get().TimeoutRegister, bod.GetHeader())
	})
	
	timing.Push("bradesco-shopfacil-register-boleto-online", duration.Seconds())
	b.log.ResponseCustom(response, map[string]string{"ContentStatusCode": strconv.Itoa(status), "URL" : serviceURL, "Duration" : util.ConvertDuration(duration)})

	if err != nil {
		return models.BoletoResponse{}, err
	}

	var result interface{}

	if status == 201 || status == 200{
		result = tmpl.TransformFromJSON(response, getResponseBradescoShopFacil(), getAPIResponseBradescoShopFacil(), new(models.BoletoResponse))
	} 

	switch t := result.(type) {
	case *models.BoletoResponse:
		if !t.HasErrors() {
			t.BarCodeNumber = getBarcode(*boleto).toString()
		}
		return *t, nil
	case error:
		return models.BoletoResponse{}, t
	default:
		return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
	}
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

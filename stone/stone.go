package stone

import (
	"io/ioutil"
	"strings"

	openBank "github.com/stone-co/go-stone-openbank"
	openBankTypes "github.com/stone-co/go-stone-openbank/types"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

type bankStone struct {
	validate *models.Validator
	log      *log.Log
}

func New() bankStone {
	b := bankStone{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	return b
}

func (b bankStone) Log() *log.Log {
	return b.log
}

func (b bankStone) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	var response models.BoletoResponse

	paymentInvoiceRequest := openBankTypes.PaymentInvoiceInput{
		AccountID:      boleto.Authentication.Username,
		Amount:         int(boleto.Title.AmountInCents),
		ExpirationDate: boleto.Title.ExpireDate,
		LimitDate:      boleto.Title.ExpireDate,
		InvoiceType:    boleto.Title.BoletoType,
		Payer: openBankTypes.PaymentInvoicePayerInput{
			Document:  boleto.Buyer.Document.Number,
			LegalName: boleto.Buyer.Name,
		},
	}

	client, err := b.Authenticate()
	if err != nil {
		return response, err
	}

	resp := &openBank.Response{}
	paymentInvoiceResponse := &openBankTypes.PaymentInvoice{}
	duration := util.Duration(func() {
		paymentInvoiceResponse, resp, err = client.PaymentInvoice.PaymentInvoice(paymentInvoiceRequest, boleto.RequestKey)
	})
	if err != nil {
		return response, err
	}
	metrics.PushTimingMetric("stone-register-boleto-time", duration.Seconds())

	boleto.Recipient.Name = paymentInvoiceResponse.Beneficiary.LegalName
	boleto.Recipient.Document.Number = paymentInvoiceResponse.Beneficiary.Document
	boleto.Recipient.Document.Type = strings.ToUpper(paymentInvoiceResponse.Beneficiary.DocumentType)
	boleto.Agreement.Account = paymentInvoiceResponse.Beneficiary.AccountCode
	boleto.Agreement.Agency = paymentInvoiceResponse.Beneficiary.BranchCode

	response.StatusCode = resp.StatusCode
	response.ID = paymentInvoiceResponse.ID
	response.DigitableLine = paymentInvoiceResponse.WritableLine
	response.BarCodeNumber = paymentInvoiceResponse.Barcode
	response.OurNumber = paymentInvoiceResponse.OurNumber

	return response, nil
}

func (b bankStone) Authenticate() (*openBank.Client, error) {
	clientID := config.Get().StoneClientID
	privKeyPath := config.Get().StonePrivateKeyPath
	pemPrivKey, _ := ioutil.ReadFile(privKeyPath)

	client, err := openBank.NewClient(openBank.WithClientID(clientID),openBank.WithPEMPrivateKey(pemPrivKey))
	if err != nil {
		return client, err
	}

	if config.Get().Environment != "Production" {
		client.ApplyOpts(openBank.UseSandbox(), openBank.EnableDebug())
	}

	err = client.Authenticate()
	if err != nil {
		return client, err
	}

	return client, nil
}

func (b bankStone) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	return b.RegisterBoleto(boleto)
}

func (b bankStone) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return nil
}

func (b bankStone) GetBankNumber() models.BankNumber {
	return models.Stone
}

func (b bankStone) GetBankNameIntegration() string {
	return "Stone"
}

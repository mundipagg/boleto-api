package stone

import (
	"fmt"
	"testing"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "bill_of_exchange"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "bill_of_exchange"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "bill_of_exchange"},
}

var boletoResponseFailParameters = []test.Parameter{
	{Input: newStubBoletoRequestStone().WithAccessKey("").Build(), Expected: models.ErrorResponse{Code: `MP400`, Message: `o campo AccessKey não pode ser vazio`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(200).Build(), Expected: models.ErrorResponse{Code: "MPOurNumberFail", Message: "our number was not returned by the bank"}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(401).Build(), Expected: models.ErrorResponse{Code: `srn:error:unauthenticated`, Message: `srn:error:unauthenticated`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(403).Build(), Expected: models.ErrorResponse{Code: `srn:error:unauthorized`, Message: `srn:error:unauthorized`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(409).Build(), Expected: models.ErrorResponse{Code: `srn:error:conflict`, Message: `srn:error:conflict`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(422).Build(), Expected: models.ErrorResponse{Code: `srn:error:product_not_enabled`, Message: `barcode_payment_invoice_bill_of_exchange is not ena bled on this account`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4001).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:is invalid,path:[customer,document]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4002).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:can&#39;t be blank,path:[customer,legal_name]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4003).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:not allowed,path:[amount]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4004).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:is invalid,path:[receiver,document]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4005).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:is invalid,path:[account_id]},{error:not allowed,path:[amount]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(504).Build(), Expected: models.ErrorResponse{Code: `MPTimeout`, Message: `Post http://localhost:9099/stone/registrarBoleto: context deadline exceeded`}},
}

func Test_GetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}
func Test_TemplateRequestStone_WhenBuyerIsPerson_ParseSuccessful(t *testing.T) {
	var result map[string]interface{}
	f := flow.NewFlow()
	input := newStubBoletoRequestStone().WithBoletoType("DM").Build()

	body := fmt.Sprintf("%v", f.From("message://?source=inline", input, templateRequest, tmpl.GetFuncMaps()).GetBody())
	util.FromJSON(body, &result)

	assert.Equal(t, result["account_id"], input.Authentication.AccessKey)
	assert.Equal(t, uint64(result["amount"].(float64)), input.Title.AmountInCents)
	assert.Equal(t, result["expiration_date"], input.Title.ExpireDate)
	assert.Equal(t, result["invoice_type"], input.Title.BoletoTypeCode)
	assert.Equal(t, result["customer"].(map[string]interface{})["document"], input.Buyer.Document.Number)
	assert.Equal(t, result["customer"].(map[string]interface{})["legal_name"], input.Buyer.Name)
	assert.Equal(t, result["customer"].(map[string]interface{})["trade_name"], nil)
}
func Test_TemplateRequestStone_WhenBuyerIsCompany_ParseSuccessful(t *testing.T) {
	var result map[string]interface{}
	f := flow.NewFlow()
	input := newStubBoletoRequestStone().WithDocument("12123123000112", "CNPJ").WithBoletoType("DM").Build()

	body := fmt.Sprintf("%v", f.From("message://?source=inline", input, templateRequest, tmpl.GetFuncMaps()).GetBody())
	util.FromJSON(body, &result)

	assert.Equal(t, result["account_id"], input.Authentication.AccessKey)
	assert.Equal(t, uint64(result["amount"].(float64)), input.Title.AmountInCents)
	assert.Equal(t, result["expiration_date"], input.Title.ExpireDate)
	assert.Equal(t, result["invoice_type"], input.Title.BoletoTypeCode)
	assert.Equal(t, result["customer"].(map[string]interface{})["document"], input.Buyer.Document.Number)
	assert.Equal(t, result["customer"].(map[string]interface{})["legal_name"], input.Buyer.Name)
	assert.Equal(t, result["customer"].(map[string]interface{})["trade_name"], input.Buyer.Name)
}

func Test_ProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9099")

	input := newStubBoletoRequestStone().WithAmountInCents(201).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func Test_ProcessBoleto_WhenServiceRespondsUnsuccessful_ShouldHasErrorResponse(t *testing.T) {
	bank := New()
	mock.StartMockService("9099")

	for _, fact := range boletoResponseFailParameters {
		request := fact.Input.(*models.BoletoRequest)
		response, _ := bank.ProcessBoleto(request)

		test.AssertProcessBoletoFailed(t, response)
		assert.Equal(t, fact.Expected.(models.ErrorResponse).Code, response.Errors[0].Code)
		assert.Equal(t, fact.Expected.(models.ErrorResponse).Message, response.Errors[0].Message)
	}
}

func Test_GetBankNumber(t *testing.T) {
	bank := New()

	result := bank.GetBankNumber()

	assert.Equal(t, models.Stone, int(result))
}

func Test_GetBankNameIntegration(t *testing.T) {
	bank := New()

	result := bank.GetBankNameIntegration()

	assert.Equal(t, "Stone", result)
}

func Test_GetBankLog(t *testing.T) {
	bank := New()

	result := bank.Log()

	assert.NotNil(t, result)
}

func Test_bankStone_ProcessBoleto(t *testing.T) {
	mock.StartMockService("9093")

	bankInst := New()

	type args struct {
		request *models.BoletoRequest
	}
	tests := []struct {
		name    string
		b       bankStone
		args    args
		want    models.BoletoResponse
		wantErr bool
	}{
		{
			name: "StoneEmptyAccessKeyRequest",
			b:    bankInst,
			args: args{
				request: successRequest,
			},
			want: models.BoletoResponse{
				StatusCode: 0,
				Errors: []models.ErrorResponse{
					{
						Code:    "MP400",
						Message: "o campo AccessKey não pode ser vazio",
					},
				},
				ID:            "",
				DigitableLine: "",
				BarCodeNumber: "",
				OurNumber:     "",
				Links:         []models.Link{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.b.ProcessBoleto(tt.args.request)
			assert.Greater(t, len(got.Errors), 0)
			err := got.Errors[0]
			assert.Equal(t, err.Code, "MP400")
			assert.Equal(t, err.Message, "o campo AccessKey não pode ser vazio")
		})
	}
}

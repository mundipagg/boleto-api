package stone

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
)

const day = time.Hour * 24

type stubBoletoRequestStone struct {
	test.StubBoletoRequest
}

//newStubBoletoRequestStone Cria um novo StubBoletoRequest com valores default validáveis para Stone
func newStubBoletoRequestStone() *stubBoletoRequestStone {
	expirationDate := time.Now().Add(5 * day)

	base := test.NewStubBoletoRequest(models.Stone)
	s := &stubBoletoRequestStone{
		StubBoletoRequest: *base,
	}

	s.Authentication = models.Authentication{
		Username:           "VsKkTASTTdri0",
		Password:           "Tkms6VwoPdjLWFCLOLhYt_KbV2hIvdWqmNKQX7XOVclTnigKXmn6CqQMf2UxhVoo",
		AuthorizationToken: "",
		AccessKey:          "946b50ce-ed5d-45ab-8c86-ce3baf90a73a",
	}

	s.Title = models.Title{
		ExpireDateTime: expirationDate,
		ExpireDate:     expirationDate.Format("2006-01-02"),
		AmountInCents:  201,
		Instructions:   "Sr. Caixa, favor não receber após vencimento",
		DocumentNumber: "999999999999999",
	}

	s.Recipient = models.Recipient{
		Document: models.Document{
			Type:   "CNPJ",
			Number: "12123123000112",
		},
	}

	s.Buyer = models.Buyer{
		Name:  "Nome do Comprador (Cliente)",
		Email: "",
		Document: models.Document{
			Type:   "CPF",
			Number: "39734022059",
		},
		Address: models.Address{
			Street:     "Logradouro do Comprador",
			Number:     "1000",
			Complement: "Casa 01",
			ZipCode:    "15050466",
			City:       "Cidade do Comprador",
			District:   "Bairro do Comprador",
			StateCode:  "SP",
		},
	}

	return s
}

func (s *stubBoletoRequestStone) WithBoletoType(bt string) *stubBoletoRequestStone {
	switch bt {
	case "DM":
		s.Title.BoletoType, s.Title.BoletoTypeCode = bt, "bill_of_exchange"
	default:
		s.Title.BoletoType = bt
	}
	return s
}

func (s *stubBoletoRequestStone) WithDocument(number string, doctype string) *stubBoletoRequestStone {
	s.Buyer.Document.Type = doctype
	s.Buyer.Document.Number = number
	return s
}

func (s *stubBoletoRequestStone) WithAccessKey(key string) *stubBoletoRequestStone {
	s.Authentication.AccessKey = key
	return s
}

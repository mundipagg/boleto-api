package caixa

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
)

const day = time.Hour * 24

var expirationDate = time.Now().Add(5 * day)

type stubBoletoRequestCaixa struct {
	test.StubBoletoRequest
}

//newStubBoletoRequestCaixa Cria um novo StubBoletoRequest com valores default validáveis para Caixa
func newStubBoletoRequestCaixa() *stubBoletoRequestCaixa {
	base := test.NewStubBoletoRequest(models.Caixa)
	s := &stubBoletoRequestCaixa{
		StubBoletoRequest: *base,
	}

	s.Agreement = models.Agreement{
		AgreementNumber: 123456,
		Agency:          "1234",
	}

	s.Title = models.Title{
		ExpireDateTime: expirationDate,
		ExpireDate:     expirationDate.Format("2006-01-02"),
		OurNumber:      12345678901234,
		AmountInCents:  200,
		DocumentNumber: "1234567890A",
		Instructions:   "Campo de instruções -  max 40 caracteres",
		BoletoType:     "OUT",
		BoletoTypeCode: "99",
	}

	s.Recipient = models.Recipient{
		Document: models.Document{
			Type:   "CNPJ",
			Number: "12123123000112",
		},
	}

	s.Buyer = models.Buyer{
		Name: "Willian Jadson Bezerra Menezes Tupinambá",
		Document: models.Document{
			Type:   "CPF",
			Number: "12312312312",
		},
		Address: models.Address{
			Street:     "Rua da Assunção de Sá",
			Number:     "123",
			Complement: "Seção A, s 02",
			ZipCode:    "20520051",
			City:       "Belém do Pará",
			District:   "Açaí",
			StateCode:  "PA",
		},
	}
	return s
}

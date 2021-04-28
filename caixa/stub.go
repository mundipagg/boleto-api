package caixa

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
)

const day = time.Hour * 24

var expirationDate = time.Now().Add(5 * day)

type stubBoletoRequestCaixa struct {
	test.BuilderBoletoRequest
	authentication models.Authentication
	agreement      models.Agreement
	title          models.Title
	recipient      models.Recipient
	buyer          models.Buyer
}

//newStubBoletoRequestCaixa Cria um novo StubBoletoRequest com valores default validáveis para Caixa
func newStubBoletoRequestCaixa() *stubBoletoRequestCaixa {
	s := &stubBoletoRequestCaixa{
		BuilderBoletoRequest: test.NewBuilderBoletoRequest(),
	}

	s.authentication = models.Authentication{}

	s.agreement = models.Agreement{
		AgreementNumber: 123456,
		Agency:          "1234",
	}

	s.title = models.Title{
		ExpireDateTime: expirationDate,
		ExpireDate:     expirationDate.Format("2006-01-02"),
		OurNumber:      12345678901234,
		AmountInCents:  200,
		DocumentNumber: "1234567890A",
		Instructions:   "Campo de instruções -  max 40 caracteres",
		BoletoType:     "OUT",
		BoletoTypeCode: "99",
	}

	s.recipient = models.Recipient{
		Document: models.Document{
			Type:   "CNPJ",
			Number: "12123123000112",
		},
	}

	s.buyer = models.Buyer{
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

func (s *stubBoletoRequestCaixa) withAmountIsCents(amount uint64) *stubBoletoRequestCaixa {
	s.title.AmountInCents = amount
	return s
}

func (s *stubBoletoRequestCaixa) withOurNumber(ourNumber uint) *stubBoletoRequestCaixa {
	s.title.OurNumber = ourNumber
	return s
}

func (s *stubBoletoRequestCaixa) withRecipientDocumentNumber(docNumber string) *stubBoletoRequestCaixa {
	s.recipient.Document.Number = docNumber
	return s
}

func (s *stubBoletoRequestCaixa) withAgreementNumber(number uint) *stubBoletoRequestCaixa {
	s.agreement.AgreementNumber = number
	return s
}

func (s *stubBoletoRequestCaixa) withExpirationDate(expiredAt time.Time) *stubBoletoRequestCaixa {
	s.title.ExpireDateTime = expiredAt
	s.title.ExpireDate = expiredAt.Format("2006-01-02")
	return s
}

func (s *stubBoletoRequestCaixa) withStrictRules() *stubBoletoRequestCaixa {
	s.title.Rules = &models.Rules{
		AcceptDivergentAmount: false,
		MaxDaysToPayPastDue:   0,
	}
	return s
}

func (s *stubBoletoRequestCaixa) withFlexRules() *stubBoletoRequestCaixa {
	s.title.Rules = &models.Rules{
		AcceptDivergentAmount: true,
		MaxDaysToPayPastDue:   60,
	}
	return s
}

func (s *stubBoletoRequestCaixa) Build() *models.BoletoRequest {
	s.SetAuthentication(s.authentication)
	s.SetAgreement(s.agreement)
	s.SetTitle(s.title)
	s.SetRecipient(s.recipient)
	s.SetBuyer(s.buyer)
	return s.BoletoRequest()
}

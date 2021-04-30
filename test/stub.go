package test

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
)

//StubBoletoRequest Stub base para criação de BoletoRequest
type StubBoletoRequest struct {
	BuilderBoletoRequest
	Authentication models.Authentication
	Agreement      models.Agreement
	Title          models.Title
	Recipient      models.Recipient
	Buyer          models.Buyer
	bank           models.BankNumber
}

func NewStubBoletoRequest(bank models.BankNumber) *StubBoletoRequest {
	s := &StubBoletoRequest{
		BuilderBoletoRequest: NewBuilderBoletoRequest(),
	}

	s.bank = bank

	s.Authentication = models.Authentication{}
	s.Agreement = models.Agreement{}
	s.Title = models.Title{}
	s.Recipient = models.Recipient{}
	s.Buyer = models.Buyer{}

	return s
}

func (s *StubBoletoRequest) WithAgreementNumber(number uint) *StubBoletoRequest {
	s.Agreement.AgreementNumber = number
	return s
}

func (s *StubBoletoRequest) WithAmountInCents(amount uint64) *StubBoletoRequest {
	s.Title.AmountInCents = amount
	return s
}

func (s *StubBoletoRequest) WithOurNumber(ourNumber uint) *StubBoletoRequest {
	s.Title.OurNumber = ourNumber
	return s
}

func (s *StubBoletoRequest) WithExpirationDate(expiredAt time.Time) *StubBoletoRequest {
	s.Title.ExpireDateTime = expiredAt
	s.Title.ExpireDate = expiredAt.Format("2006-01-02")
	return s
}

func (s *StubBoletoRequest) WithAcceptDivergentAmount(accepted bool) *StubBoletoRequest {
	if !s.Title.HasRules() {
		s.Title.Rules = &models.Rules{}
	}

	s.Title.Rules.AcceptDivergentAmount = accepted
	return s
}

func (s *StubBoletoRequest) WithRecipientDocumentNumber(docNumber string) *StubBoletoRequest {
	s.Recipient.Document.Number = docNumber
	return s
}

func (s *StubBoletoRequest) Build() *models.BoletoRequest {
	s.SetAuthentication(s.Authentication)
	s.SetAgreement(s.Agreement)
	s.SetTitle(s.Title)
	s.SetRecipient(s.Recipient)
	s.SetBuyer(s.Buyer)
	s.SetBank(s.bank)
	return s.BoletoRequest()
}

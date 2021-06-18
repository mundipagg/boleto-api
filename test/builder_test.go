package test

import (
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func Test_BoletoBuilder_WhenCreateAndSetBoletoRequest_ReturnBoletoRequestSuccessful(t *testing.T) {

	authentication := models.Authentication{
		Username: "username",
	}
	agreement := models.Agreement{
		Agency:  "1234",
		Account: "123456",
	}
	title := models.Title{
		OurNumber: 1234567890,
	}
	recipient := models.Recipient{
		Name: "Recebedor",
	}
	buyer := models.Buyer{
		Name: "Comprador",
	}

	b := NewBuilderBoletoRequest()
	b.SetBank(models.Caixa)
	b.SetAuthentication(authentication)
	b.SetAgreement(agreement)
	b.SetTitle(title)
	b.SetRecipient(recipient)
	b.SetBuyer(buyer)

	r := b.BoletoRequest()

	assert.Equal(t, models.Caixa, int(r.BankNumber))
	assert.Equal(t, authentication, r.Authentication)
	assert.Equal(t, agreement, r.Agreement)
	assert.Equal(t, title, r.Title)
	assert.Equal(t, recipient, r.Recipient)
	assert.Equal(t, buyer, r.Buyer)
	assert.NotEqual(t, "00000000-0000-0000-0000-000000000000", r.RequestKey)

}

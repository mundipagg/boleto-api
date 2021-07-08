package db_test

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func TestGetTokenByClientIDAndIssuerBank(t *testing.T) {
	mock.StartMockService("9093")
	token := models.NewToken("9979b005-5k40-41c4-976e-3cec24f8006s", "SarumanBank", "Palantir")

	err := deleteTokenByIssuerBank(token.IssuerBank)
	assert.Nil(t, err)

	err = db.SaveToken(token)
	assert.Nil(t, err)

	got, err := db.GetTokenByClientIDAndIssuerBank("9979b005-5k40-41c4-976e-3cec24f8006s", "SarumanBank")
	assert.Nil(t, err)
	assert.Equal(t, "Palantir", got.AccessToken)

	deleteTokenByIssuerBank(token.IssuerBank)
}

func TestGetTokenByClientIDAndIssuerBankWithExpiratedToken(t *testing.T) {
	mock.StartMockService("9093")
	token := models.NewToken("5999h015-5l40-41c4-976e-2cec24f8006s", "Alatar", "The Blue")
	token.CreatedAt = token.CreatedAt.Add(-14 * time.Minute)

	err := deleteTokenByIssuerBank(token.IssuerBank)
	assert.Nil(t, err)

	err = db.SaveToken(token)
	assert.Nil(t, err)

	got, err := db.GetTokenByClientIDAndIssuerBank("5999h015-5l40-41c4-976e-2cec24f8006s", "Alatar")
	assert.Nil(t, err)
	assert.Equal(t, "", got.AccessToken)

	deleteTokenByIssuerBank(token.IssuerBank)
}

func TestSaveToken(t *testing.T) {
	mock.StartMockService("9093")
	token := models.NewToken("6275b002-5e20-67y1-716t-4aec24f8004w", "OlorinBank", "Mellon")

	err := deleteTokenByIssuerBank(token.IssuerBank)
	assert.Nil(t, err)

	err = db.SaveToken(token)
	assert.Nil(t, err)

	got, err := db.GetTokenByClientIDAndIssuerBank("6275b002-5e20-67y1-716t-4aec24f8004w", token.IssuerBank)
	assert.Nil(t, err)
	assert.Equal(t, "Mellon", got.AccessToken)

	deleteTokenByIssuerBank(token.IssuerBank)
}

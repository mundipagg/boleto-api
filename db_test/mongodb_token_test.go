package db_test

import (
	"testing"

	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAccessTokenByClientIDAndOrigin(t *testing.T) {
	mock.StartMockService("9093")
	token := models.NewToken("9979b005-5k40-41c4-976e-3cec24f8006s", "SarumanBank", "Palantir")

	err := deleteTokenByOrigin(token.Origin)
	assert.Nil(t, err)

	err = db.SaveAccessToken(token)
	assert.Nil(t, err)

	got, err := db.GetAccessTokenByClientIDAndOrigin("9979b005-5k40-41c4-976e-3cec24f8006s", "SarumanBank")
	assert.Nil(t, err)
	assert.Equal(t, "Palantir", got.AccessToken)

	deleteTokenByOrigin(token.Origin)
}

func TestSaveAccessToken(t *testing.T) {
	mock.StartMockService("9093")
	token := models.NewToken("6275b002-5e20-67y1-716t-4aec24f8004w", "OlorinBank", "Mellon")

	err := deleteTokenByOrigin(token.Origin)
	assert.Nil(t, err)

	err = db.SaveAccessToken(token)
	assert.Nil(t, err)

	got, err := db.GetAccessTokenByClientIDAndOrigin("6275b002-5e20-67y1-716t-4aec24f8004w", token.Origin)
	assert.Nil(t, err)
	assert.Equal(t, "Mellon", got.AccessToken)

	deleteTokenByOrigin(token.Origin)
}

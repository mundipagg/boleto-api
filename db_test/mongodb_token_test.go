package db_test

import (
	"testing"

	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAccessTokenByOrigin(t *testing.T) {
	mock.StartMockService("9093")
	token := models.NewToken("SarumanBank", "Palantir")

	err := deleteTokenByOrigin(token.Origin)
	assert.Nil(t, err)

	err = db.SaveAccessToken(token)
	assert.Nil(t, err)

	got, err := db.GetAccessTokenByOrigin("SarumanBank")
	assert.Nil(t, err)
	assert.Equal(t, "Palantir", got.AccessToken)

	deleteTokenByOrigin(token.Origin)
}

func TestSaveAccessToken(t *testing.T) {
	mock.StartMockService("9093")
	token := models.NewToken("OlorinBank", "Mellon")

	err := deleteTokenByOrigin(token.Origin)
	assert.Nil(t, err)

	err = db.SaveAccessToken(token)
	assert.Nil(t, err)

	got, err := db.GetAccessTokenByOrigin(token.Origin)
	assert.Nil(t, err)
	assert.Equal(t, "Mellon", got.AccessToken)

	deleteTokenByOrigin(token.Origin)
}

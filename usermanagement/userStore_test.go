package usermanagement

import (
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func Test_GetUser_WhenUserExists_ReturnUser(t *testing.T) {
	k := "key"
	c := models.NewCredentials(k, "pass")

	addUser(k, c)

	result, exists := GetUser(k)

	assert.True(t, exists)
	assert.Equal(t, c, result)
}

func Test_GetUser_WhenUserNotExists_ReturnNil(t *testing.T) {
	k := "key"
	c := models.NewCredentials(k, "pass")

	addUser(k, c)

	result, exists := GetUser("anotherKey")

	assert.False(t, exists)
	assert.Nil(t, result)
}

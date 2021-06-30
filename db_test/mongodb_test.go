package db_test

import (
	"testing"

	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateMongo(t *testing.T) {
	mock.StartMockService("9093")
	conn, err := db.CreateMongo()

	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestCreateMongoMustNoBeRecreated(t *testing.T) {
	mock.StartMockService("9093")
	conn1, err := db.CreateMongo()

	assert.Nil(t, err)
	assert.NotNil(t, conn1)

	conn2, err := db.CreateMongo()

	assert.Nil(t, err)
	assert.NotNil(t, conn2)

	assert.Equal(t, conn1, conn2)
}

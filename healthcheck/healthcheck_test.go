package healthcheck

import (
	"testing"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/stretchr/testify/assert"
)

func Test_EnsureDependecies_WithSucess(t *testing.T) {
	mock.StartMockService("9088")
	result := EnsureDependencies()
	assert.True(t, result)
}

package util

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestParameter struct {
	Input    interface{}
	Expected interface{}
}

var getTLSparameters = []TestParameter{
	{Input: "1.0", Expected: uint16(tls.VersionTLS10)},
	{Input: "1.1", Expected: uint16(tls.VersionTLS11)},
	{Input: "1.2", Expected: uint16(tls.VersionTLS12)},
	{Input: "1.3", Expected: uint16(tls.VersionTLS13)},
	{Input: "", Expected: uint16(tls.VersionTLS12)},
	{Input: " ", Expected: uint16(tls.VersionTLS12)},
}

func TestGetTLSVersion_WhenCall_ReturnTLSVersionSuccessul(t *testing.T) {
	for _, fact := range getTLSparameters {
		result := GetTLSVersion(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve retornar o TLS correto")
	}
}

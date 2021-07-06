package app

import (
	"testing"

	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/stretchr/testify/assert"
)

func Test_openBankSkFromBlob(t *testing.T) {
	mock.StartMockService("9093")

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Fetch sk from blob successfully",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk, err := openBankSkFromBlob()
			assert.False(t, (err != nil) != tt.wantErr)
			assert.NotNil(t, sk)
		})
	}
}

func Test_installCertificates(t *testing.T) {
	mock.StartMockService("9093")

	tests := []struct {
		name string
	}{
		{
			name: "Fetch sk from localStorage successfully",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installCertificates()
		})
		sk, err := certificate.GetCertificateFromStore(config.Get().AzureStorageOpenBankSkName)
		assert.Nil(t, err)
		assert.NotNil(t, sk)
	}
}

package certificate_test

import (
	"testing"

	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/stretchr/testify/assert"
)

func TestAzureBlob_Download(t *testing.T) {
	mock.StartMockService("9093")
	azureBlobInst, err := certificate.NewAzureBlob(
		config.Get().AzureStorageAccount,
		config.Get().AzureStorageAccessKey,
		config.Get().AzureStorageContainerName,
	)
	assert.Nil(t, err)

	type args struct {
		path     string
		filename string
	}
	tests := []struct {
		name    string
		ab      *certificate.AzureBlob
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Donwload successfully",
			ab:   azureBlobInst,
			args: args{
				path:     config.Get().AzureStorageOpenBankSkPath,
				filename: config.Get().AzureStorageOpenBankSkName,
			},
			want:    "secret",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ab.Download(tt.args.path, tt.args.filename)
			assert.False(t, (err != nil) != tt.wantErr)
			assert.NotNil(t, got)
		})
	}
}

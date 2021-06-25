package stonebank

import (
	"testing"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/stretchr/testify/assert"
)

func Test_accessToken(t *testing.T) {
	mock.StartMockService("9093")

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "Get AccessToken",
			want:    "xxx",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := accessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("accessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

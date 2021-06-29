package stone

import (
	"testing"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func Test_bankStone_ProcessBoleto(t *testing.T) {
	mock.StartMockService("9093")

	bankInst, err := New()
	assert.Nil(t, err)

	type args struct {
		request *models.BoletoRequest
	}
	tests := []struct {
		name    string
		b       stone
		args    args
		want    models.BoletoResponse
		wantErr bool
	}{
		{
			name: "StoneEmptyAccessKeyRequest",
			b:    bankInst,
			args: args{
				request: successRequest,
			},
			want: models.BoletoResponse{
				StatusCode: 0,
				Errors: []models.ErrorResponse{
					{
						Code:    "MP400",
						Message: "o campo AccessKey não pode ser vazio",
					},
				},
				ID:            "",
				DigitableLine: "",
				BarCodeNumber: "",
				OurNumber:     "",
				Links:         []models.Link{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.b.ProcessBoleto(tt.args.request)
			assert.Greater(t, len(got.Errors), 0)
			err := got.Errors[0]
			assert.Equal(t, err.Code, "MP400")
			assert.Equal(t, err.Message, "o campo AccessKey não pode ser vazio")
		})
	}
}

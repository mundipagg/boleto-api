package stonebank

import (
	"reflect"
	"testing"

	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
)

func Test_bankStoneBank_ProcessBoleto(t *testing.T) {
	mock.StartMockService("9093")

	type args struct {
		request *models.BoletoRequest
	}
	tests := []struct {
		name    string
		b       bankStoneBank
		args    args
		want    models.BoletoResponse
		wantErr bool
	}{
		{
			name: "BankStoneSuccessfullRequest",
			b: bankStoneBank{
				validate: &models.Validator{
					Rules: []models.Rule{},
				},
				log: log.CreateLog(),
			},
			args: args{
				request: successRequest,
			},
			want: models.BoletoResponse{
				StatusCode:    0,
				Errors:        []models.ErrorResponse{},
				ID:            "",
				DigitableLine: "",
				BarCodeNumber: "",
				OurNumber:     "",
				Links:         []models.Link{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.ProcessBoleto(tt.args.request)
			if (err != nil) != tt.wantErr {
				// t.Skip("Not implemented yet")
				t.Errorf("bankStoneBank.ProcessBoleto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				// t.Skip("Not implemented yet")
				t.Errorf("bankStoneBank.ProcessBoleto() = %v, want %v", got, tt.want)
			}
		})
	}
}

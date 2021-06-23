package stonebank

import (
	"reflect"
	"testing"
	"time"

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
			name: "BankStone",
			b: bankStoneBank{
				validate: &models.Validator{
					Rules: []models.Rule{},
				},
				log: log.CreateLog(),
			},
			args: args{
				request: &models.BoletoRequest{
					Authentication: models.Authentication{
						Username:           "VsKkTASTTdri0",
						Password:           "Tkms6VwoPdjLWFCLOLhYt_KbV2hIvdWqmNKQX7XOVclTnigKXmn6CqQMf2UxhVoo",
						AuthorizationToken: "",
						AccessKey:          "",
					},
					Agreement: models.Agreement{
						AgreementNumber: 0,
						Wallet:          109,
						WalletVariation: 0,
						Agency:          "2938",
						AgencyDigit:     "",
						Account:         "23195",
						AccountDigit:    "4",
					},
					Title: models.Title{
						CreateDate:     time.Now(),
						ExpireDateTime: time.Now().Add(5 * time.Hour * 24),
						ExpireDate:     time.Now().Format("2006-05-11"),
						AmountInCents:  200,
						OurNumber:      94726341,
						Instructions:   "Sr. Caixa, favor não receber após vencimento",
						DocumentNumber: "999999999999999",
						NSU:            "",
						BoletoType:     "BDP",
						BoletoTypeCode: "",
					},
					Recipient: models.Recipient{
						Name: "Nome do Recebedor (Loja)",
						Document: models.Document{
							Type:   "CNPJ",
							Number: "14068605000129",
						},
						Address: models.Address{
							Street:     "Logradouro do Recebedor",
							Number:     "1000",
							Complement: "Sala 01",
							ZipCode:    "00000000",
							City:       "Cidade do Recebedor",
							District:   "Bairro do Recebedor",
							StateCode:  "RJ",
						},
					},
					Buyer: models.Buyer{
						Name:  "Nome do Comprador (Cliente)",
						Email: "",
						Document: models.Document{
							Type:   "CPF",
							Number: "39734022059",
						},
						Address: models.Address{
							Street:     "Logradouro do Comprador",
							Number:     "1000",
							Complement: "Casa 01",
							ZipCode:    "15050466",
							City:       "Cidade do Comprador",
							District:   "Bairro do Comprador",
							StateCode:  "SP",
						},
					},
					BankNumber: 197,
					RequestKey: "d26039c8-1bf4-4b42-8fc9-7b0cf0534ebc",
				},
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
				t.Errorf("bankStoneBank.ProcessBoleto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bankStoneBank.ProcessBoleto() = %v, want %v", got, tt.want)
			}
		})
	}
}

package stone

import (
	"testing"

	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/stretchr/testify/assert"
)

func Test_authenticate(t *testing.T) {
	mock.StartMockService("9093")

	pkByte, err := generateTestPK()
	assert.Nil(t, err)
	certificate.SetCertificateOnStore(config.Get().AzureStorageOpenBankSkName, pkByte)

	type args struct {
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Successfully create access token",
			args: args{
				clientID: "94j3943-394329x",
			},
			want:    "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI4d3NUd3BhYTRJWUZIYWV5ZFRubnRoRC1UaVlCaU9kanNmOGx6RUlMR1hVIn0.eyJqdGkiOiIyZTlkNGZkMy0zN2M1LTRjOWUtYTJjYy1lMjQ1N2MxZDgyMWQiLCJleHAiOjE2MjQ4OTQ1NDQsIm5iZiI6MCwiaWF0IjoxNjI0ODkzNjQ0LCJpc3MiOiJodHRwczovL3NhbmRib3gtYWNjb3VudHMub3BlbmJhbmsuc3RvbmUuY29tLmJyL2F1dGgvcmVhbG1zL3N0b25lX2JhbmsiLCJzdWIiOiJkNDY0ZDg3MC1mYzc2LTRjZGMtYWM5OC1hNjcyYjYyOTdhOGYiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiIzMjc5YjAwNS01ZTQwLTQxYzEtOTk2ZS04Y2VjMjRmODAwNmIiLCJhdXRoX3RpbWUiOjAsInNlc3Npb25fc3RhdGUiOiJhM2MyYzY3OC0wOGIxLTRmNmQtYmQ2Yi0wNjgzMjQ1M2UzNmMiLCJhY3IiOiIxIiwic2NvcGUiOiJwYXltZW50YWNjb3VudDpwYXltZW50bGlua3M6d3JpdGUgcGF5bWVudGFjY291bnQ6Y29udGFjdDp3cml0ZSBwaXg6cGF5bWVudF9pbnZvaWNlIHBpeDpwYXltZW50IHBpeDplbnRyeV9jbGFpbSBwYXltZW50YWNjb3VudDpyZWFkIHBpeDplbnRyeSBwYXltZW50YWNjb3VudDp0cmFuc2ZlcnM6aW50ZXJuYWwgcGF5bWVudGFjY291bnQ6ZmVlczpyZWFkIHBheW1lbnRhY2NvdW50OnBheW1lbnRzIHN0b25lX3N1YmplY3RfaWQgcGF5bWVudGFjY291bnQ6Y29udGFjdDpyZWFkIHNpZ251cDpwYXltZW50YWNjb3VudCBwYXltZW50YWNjb3VudDpib2xldG9pc3N1YW5jZSBwYXltZW50YWNjb3VudDpwYXltZW50bGlua3M6cmVhZCBwYXltZW50YWNjb3VudDp0cmFuc2ZlcnM6ZXh0ZXJuYWwiLCJjbGllbnRJZCI6IjMyNzliMDA1LTVlNDAtNDFjMS05OTZlLThjZWMyNGY4MDA2YiIsImNsaWVudEhvc3QiOiIxMC4xMC4zLjE3MiIsInN0b25lX3N1YmplY3RfaWQiOiJhcHBsaWNhdGlvbjozMjc5YjAwNS01ZTQwLTQxYzEtOTk2ZS04Y2VjMjRmODAwNmIiLCJjbGllbnRBZGRyZXNzIjoiMTAuMTAuMy4xNzIifQ.JloXzaTUFW0IVDi191U_WujRLhIIPiZUZngDb1nbhHo9mclG176CIgdSsBPmoOZr35ry47JCLgEq5ZAos8Sts72kpi1BivvVq0rJn5_NrmSyb0zqMSK4sNYzbhBafK7U6wamUZCjDeJmQ_wBUDNvPxGC1gToreMFnhrbak0pQr_CWp9Csgkn-9QUvFFpTRkJ3fdca57YnKoGsEWJWMs8Suq6g097244EWHISlUtO1ZGt01mypDeU8g5Z_eYD8qdN_woUeCGL86QhDoH-V8Dl_NIwbsHGTm8iRDTqjRBid2XH6Cj0RAMH10EpTKI8buSBzJ872bKLoCwXQQUnWIYT4Q",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authenticate(tt.args.clientID)
			assert.Nil(t, err)
			assert.False(t, (err != nil) != tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

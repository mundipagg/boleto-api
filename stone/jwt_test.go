package stone

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/stretchr/testify/assert"
)

func Test_generateJWT(t *testing.T) {
	mock.StartMockService("9093")

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "generate with success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateJWT()
			assert.Nil(t, err)
			assert.Equal(t, (err != nil), tt.wantErr)
			assert.NotEmpty(t, got)
		})
	}
}

func Test_generateJTIFromTime(t *testing.T) {
	layout := "2006-01-02T15:04:05.000Z"
	tStr := "2021-06-24T19:54:26.371Z"
	expTime, _ := time.Parse(layout, tStr)
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "jwt generations",
			args: args{expTime},
			want: "20210624195426371",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, generateJTIFromTime(tt.args.t), tt.want)
		})
	}
}

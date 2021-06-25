package stonebank

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
		want    string
		wantErr bool
	}{
		{
			name:    "generate with success",
			want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjQ1NzM4MzksIm5iZiI6MTYyNDU2MTE5MywiYXVkIjoiaHR0cHM6Ly9zYW5kYm94LWFjY291bnRzLm9wZW5iYW5rLnN0b25lLmNvbS5ici9hdXRoL3JlYWxtcy9zdG9uZV9iYW5rIiwicmVhbG0iOiJzdG9uZV9iYW5rIiwic3ViIjoiMzI3OWIwMDUtNWU0MC00MWMxLTk5NmUtOGNlYzI0ZjgwMDZiIiwiY2xpZW50SWQiOiIzMjc5YjAwNS01ZTQwLTQxYzEtOTk2ZS04Y2VjMjRmODAwNmIiLCJpYXQiOjE2MjQ1NjExOTMsImp0aSI6IjIwMjEwNjIzMTQwNzM0MTIzLjg5Y2QzNTcifQ.dqIwPr85Y-Bqr0ucmXlBM-ddo_Fj11i1ps3UgI_9Hr_G3XL5y-IJYtFC9H4BEB6eHMaAkF5YhkNLKr1yZUHfXqQrNRgc6KHEImwkpMR1VV7kFZGK7MuLHT7tuEFK0z9Jbw_INmKAZml9rX3HoaV0yA2_altQR8PhDZ6aaf_gGhhD9b2kFyoXYu3dTAFUmkoB4HPZICux0Fu-hQfNKDqk9KfoYyN-Be93XYG6aYjcunIaeTlhPRmi7yje85Itrb4NyRmfsebNMv-csTqNmEUQYS7nkSrPZpqiR5ke0BicQfeKUMdeilt4uoF1xxBTcijzN15CJeWyD76ZGhRoTI6BuJOuTwrC73oq2zVXD34V7_7GVX4Ivl1-3bkdiY2Hs5F7XtbAkhsNf3bTv9ymQbvOfKad3rx81hXRM0rKuJ95nCP_EL_9hzl7crfiJAn7dhEVK0qlr9wr-sK7ObbIGVCISN5fN7bJwnSXH675UKhuxyuotFNLN7Wy9O9FyeLlwKSr5u8ThYSMvMOxmeXcd1j3sx8qBKbO--Hlo2m5QZow_rw76gLRPXIKFg7KB0aJfbsHKsLE-hv9D4Fez9EmG1qLxfJvp4OX8lMmSjI4SYv_K3mhbcqbIYkyztOP0v2tLmMdE6jrmj_0atkVgOkph6xFpuqfZ-L6UJiM-HjF-2503GE",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateJWT()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateJWT() = %v, want %v", got, tt.want)
			}
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

package stonebank

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

type bankStoneBank struct {
	validate *models.Validator
	log      *log.Log
}

//New Create a new Santander Integration Instance
func New() (bankStoneBank, error) {
	var err error
	b := bankStoneBank{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	if err != nil {
		return bankStoneBank{}, err
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)

	return b, nil
}

func (b bankStoneBank) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	b.accessToken()
	// if ticket, err := b.accessToken(); err != nil {
	// 	return models.BoletoResponse{Errors: errs}, err
	// } else {
	// 	boleto.Authentication.AuthorizationToken = ticket
	// }
	return b.RegisterBoleto(boleto)
}

func (b bankStoneBank) RegisterBoleto(request *models.BoletoRequest) (models.BoletoResponse, error) {
	b.log.Info(fmt.Sprintf("StoneBank Register Boleto %v", request))
	return models.BoletoResponse{}, nil
}

func (b bankStoneBank) ValidateBoleto(request *models.BoletoRequest) models.Errors {
	b.log.Info(fmt.Sprintf("StoneBank ValidateBoleto %v", request))
	return nil
}

func (b bankStoneBank) GetBankNumber() models.BankNumber {
	b.log.Info("StoneBank GetBankNumber")
	return 197
}

func (b bankStoneBank) GetBankNameIntegration() string {
	b.log.Info("StoneBank GetBankNameIntegration")
	return "stonebank"
}

func (b bankStoneBank) Log() *log.Log {
	b.log.Info("StoneBank Log")
	return b.log
}

func (b bankStoneBank) accessToken() string {
	httpClient := &util.HTTPClient{}
	endpoint := "https://sandbox-accounts.openbank.stone.com.br/auth/realms/stone_bank/protocol/openid-connect/token"

	jwt := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjQ1NzM4MzksIm5iZiI6MTYyNDU2MTE5MywiYXVkIjoiaHR0cHM6Ly9zYW5kYm94LWFjY291bnRzLm9wZW5iYW5rLnN0b25lLmNvbS5ici9hdXRoL3JlYWxtcy9zdG9uZV9iYW5rIiwicmVhbG0iOiJzdG9uZV9iYW5rIiwic3ViIjoiMzI3OWIwMDUtNWU0MC00MWMxLTk5NmUtOGNlYzI0ZjgwMDZiIiwiY2xpZW50SWQiOiIzMjc5YjAwNS01ZTQwLTQxYzEtOTk2ZS04Y2VjMjRmODAwNmIiLCJpYXQiOjE2MjQ1NjExOTMsImp0aSI6IjIwMjEwNjIzMTQwNzM0MTIzLjg5Y2QzNTcifQ.dqIwPr85Y-Bqr0ucmXlBM-ddo_Fj11i1ps3UgI_9Hr_G3XL5y-IJYtFC9H4BEB6eHMaAkF5YhkNLKr1yZUHfXqQrNRgc6KHEImwkpMR1VV7kFZGK7MuLHT7tuEFK0z9Jbw_INmKAZml9rX3HoaV0yA2_altQR8PhDZ6aaf_gGhhD9b2kFyoXYu3dTAFUmkoB4HPZICux0Fu-hQfNKDqk9KfoYyN-Be93XYG6aYjcunIaeTlhPRmi7yje85Itrb4NyRmfsebNMv-csTqNmEUQYS7nkSrPZpqiR5ke0BicQfeKUMdeilt4uoF1xxBTcijzN15CJeWyD76ZGhRoTI6BuJOuTwrC73oq2zVXD34V7_7GVX4Ivl1-3bkdiY2Hs5F7XtbAkhsNf3bTv9ymQbvOfKad3rx81hXRM0rKuJ95nCP_EL_9hzl7crfiJAn7dhEVK0qlr9wr-sK7ObbIGVCISN5fN7bJwnSXH675UKhuxyuotFNLN7Wy9O9FyeLlwKSr5u8ThYSMvMOxmeXcd1j3sx8qBKbO--Hlo2m5QZow_rw76gLRPXIKFg7KB0aJfbsHKsLE-hv9D4Fez9EmG1qLxfJvp4OX8lMmSjI4SYv_K3mhbcqbIYkyztOP0v2tLmMdE6jrmj_0atkVgOkph6xFpuqfZ-L6UJiM-HjF-2503GE"
	m := map[string]string{
		"client_id":             "3279b005-5e40-41c1-996e-8cec24f8006b",
		"grant_type":            "client_credentials",
		"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
		"client_assertion":      jwt,
	}

	resp, err := httpClient.PostFormURLEncoded(endpoint, m)
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	var r Response
	err = json.Unmarshal(responseBody, &r)

	fmt.Printf("------> Stone response: [%v]\n", r)
	fmt.Printf("------> Error: [%s]\n", err)
	fmt.Printf("------> Stone response: [%s]\n", string(r.AccessToken))

	return r.AccessToken
}

type Request struct {
	ClientID            string `json:"client_id"`
	GrantType           string `json:"client_credentials"`
	ClientAssertionType string `json:"client_assertion_type"`
	ClientAssertion     string `json:"client_assertion"`
}

type Response struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenCreatedAt int    `json:"refresh_expires_in"`
	TokenType             string `json:"token_type"`
	NotBeforePolicy       int    `json:"not-before-policy"`
	SessionState          string `json:"session_state"`
	Scope                 string `json:"scope"`
}

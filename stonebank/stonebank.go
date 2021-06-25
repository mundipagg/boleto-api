package stonebank

import (
	"fmt"

	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
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
	accToken := b.accessToken()
	fmt.Printf("AccessToken [%s]\n", accToken)

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
	token, err := accessToken()
	if err != nil {
		return ""
	}

	return token
	// httpClient := &util.HTTPClient{}
	// endpoint := "https://sandbox-accounts.openbank.stone.com.br/auth/realms/stone_bank/protocol/openid-connect/token"

	// // jwt := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjQ1ODE1MzUsIm5iZiI6MTYyNDU2ODg4NSwiYXVkIjoiaHR0cHM6Ly9zYW5kYm94LWFjY291bnRzLm9wZW5iYW5rLnN0b25lLmNvbS5ici9hdXRoL3JlYWxtcy9zdG9uZV9iYW5rIiwicmVhbG0iOiJzdG9uZV9iYW5rIiwic3ViIjoiMzI3OWIwMDUtNWU0MC00MWMxLTk5NmUtOGNlYzI0ZjgwMDZiIiwiY2xpZW50SWQiOiIzMjc5YjAwNS01ZTQwLTQxYzEtOTk2ZS04Y2VjMjRmODAwNmIiLCJpYXQiOjE2MjQ1Njg4ODUsImp0aSI6IjIwMjEwNjIzMTQwNzM0MTIzLjg5Y2QzNTgifQ.1zSmt2x9bs_Vkig4M9YU8G4xvz52cgZfErNNCcxtlqkVIVbGRZJmodCdTsh4PljT8s2roF7g9s-yT6quqInyvkYWEMC_Oye46aaJBvvz0p7VzbFEniXIl0lkwCwL_y9Q0kpDt_fFDPFSCgmV4kcrOhq4WFqTqkDYVGGtk5g8sMY93N0AgHcmk0wFg7q1M-4qUvMgOXYjs9A-bQ8e9ZfmJOZ8b5Df3X1wEKlvLqQmj58qLQJBL1qJLF6o9mSaphmmX2E3ceELqBzDqE3qECCL3qnHrc7XOdOgHrWDu31hpwfyJgN73jyIRLQgepWN9T4PKMCJJL0-5OiQKHd3kmq6UnlKo4-zaBNKIpbLE6KC1S1H-s5wNW7iyhn3cO6Xh0htgWusnqfxShWufkovxG3Smbn2DCACttAoygJRnZcHenka_1rEzu_q_-AOOfCDwsFH8nz1qsAdJ8ZCh9PTawuW02qsIj316KfSZ_UH-md0n3oAErQOsX57oXWQXIIhf-OJD0j5RAkI8wLazWvQp-Rj8mCuOAJpH9w89LWlFqdzRVfxBzTqgtbUxcW7eldgxpRbKLt2EHQpHeo_mNe7mx8eUvlTYWBtMJznUxqt-_PP89hy4-loXp4xOEQGBnYXmAViyTUDQT8m2N-2ihMfFwjmZeh8zdTIdjEX2CCzM7TcZA8"
	// jwt, err := generateJWT()
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// 	return ""
	// }

	// m := map[string]string{
	// 	"client_id":             "3279b005-5e40-41c1-996e-8cec24f8006b",
	// 	"grant_type":            "client_credentials",
	// 	"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
	// 	"client_assertion":      jwt,
	// }

	// resp, err := httpClient.PostFormURLEncoded(endpoint, m)
	// defer resp.Body.Close()

	// responseBody, err := ioutil.ReadAll(resp.Body)
	// var r Response
	// err = json.Unmarshal(responseBody, &r)

	// fmt.Printf("------> Stone response: [%v]\n", r)
	// fmt.Printf("------> Error: [%s]\n", err)
	// fmt.Printf("------> Stone response: [%s]\n", string(r.AccessToken))

	// return r.AccessToken
}

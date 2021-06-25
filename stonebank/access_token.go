package stonebank

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

var (
	l                  = log.CreateLog()
	HttpClient         = &util.HTTPClient{}
	AccessTokenPayload = map[string]string{
		"client_id":             "3279b005-5e40-41c1-996e-8cec24f8006b",
		"grant_type":            "client_credentials",
		"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
		"client_assertion":      "",
	}
)

const (
	TokenOrigin = "stonebank"
	endpoint    = "https://sandbox-accounts.openbank.stone.com.br/auth/realms/stone_bank/protocol/openid-connect/token"
)

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

// TODO: handle concurrency situations
// TODO: receives client_id
func accessToken() (string, error) {
	tk := fetchAccessTokenFromStorage()
	if tk != "" {
		return tk, nil
	}

	// request a new access token
	tk, err := requestAccessToken()
	if err != nil {
		return "", err
	}

	// saves the new access token
	token := models.NewToken(TokenOrigin, tk)
	db.SaveAccessToken(token)

	// return token
	return tk, nil
}

func fetchAccessTokenFromStorage() string {
	token, err := db.GetAccessTokenByOrigin(TokenOrigin)
	if err != nil {
		l.Error(err.Error(), "Error fetching stonebank access token")
		return ""
	}

	return token.AccessToken
}

func requestAccessToken() (string, error) {
	jwt, err := generateJWT()
	if err != nil {
		l.Error(err.Error(), "Error generating jwt")
		return "", err
	}

	AccessTokenPayload["client_assertion"] = jwt
	resp, err := HttpClient.PostFormURLEncoded(endpoint, AccessTokenPayload)
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	var r Response
	err = json.Unmarshal(responseBody, &r)

	return r.AccessToken, nil
}

package stonebank

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

var (
	l                  = log.CreateLog()
	HttpClient         = &util.HTTPClient{}
	mu                 sync.Mutex
	AccessTokenPayload = map[string]string{
		"client_id":             "",
		"grant_type":            "client_credentials",
		"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
		"client_assertion":      "",
	}
)

const (
	TokenOrigin     = "stonebank"
	BadRequestError = "status code 400"
)

type AcessTokenRequest struct {
	ClientID            string `json:"client_id"`
	GrantType           string `json:"client_credentials"`
	ClientAssertionType string `json:"client_assertion_type"`
	ClientAssertion     string `json:"client_assertion"`
}

type AccessTokenResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenCreatedAt int    `json:"refresh_expires_in"`
	TokenType             string `json:"token_type"`
	NotBeforePolicy       int    `json:"not-before-policy"`
	SessionState          string `json:"session_state"`
	Scope                 string `json:"scope"`
}

func accessToken(clientID string) (string, error) {
	tk := fetchAccessTokenFromStorage(clientID)
	if tk != "" {
		return tk, nil
	}

	return requestAndSaveAccessToken(clientID)
}

func requestAndSaveAccessToken(clientID string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	tk, err := requestAccessTokenWithRetryOnBadRequest(clientID)
	if err != nil {
		l.Error(err.Error(), "Error fetching stonebank access token")
		return "", err
	}

	// saves the new access token
	token := models.NewToken(clientID, TokenOrigin, tk)
	db.SaveAccessToken(token)

	// return token
	return tk, nil
}

func fetchAccessTokenFromStorage(clientID string) string {
	token, err := db.GetAccessTokenByClientIDAndOrigin(clientID, TokenOrigin)
	if err != nil {
		l.Error(err.Error(), "Error fetching stonebank access token")
		return ""
	}

	return token.AccessToken
}

// requestAccessTokenWithRetryOnBadRequest encapsulates logic for retry access token request once again
// in bad request status code. That's because duplicated jti returns this mencioned status code
func requestAccessTokenWithRetryOnBadRequest(clientID string) (string, error) {
	var tk string
	var err error

	if tk, err = requestAccessToken(clientID); err != nil {
		if !strings.Contains(err.Error(), BadRequestError) {
			return "", err
		}
		return requestAccessToken(clientID)
	}

	return tk, nil
}

func requestAccessToken(clientID string) (string, error) {
	jwt, err := generateJWT()
	if err != nil {
		l.Error(err.Error(), "Error generating jwt")
		return "", err
	}

	AccessTokenPayload["client_assertion"] = jwt
	AccessTokenPayload["client_id"] = clientID
	resp, err := HttpClient.PostFormURLEncoded(config.Get().URLStoneBankToken, AccessTokenPayload)
	defer resp.Body.Close()

	if err != nil {
		l.Error(err.Error(), "Error requesting a new stonebank access token")
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Request for clientID %s returns status code %d", clientID, resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Error(err.Error(), "Error reading stonebank access token response")
		return "", err
	}

	var r AccessTokenResponse
	err = json.Unmarshal(responseBody, &r)
	if err != nil {
		l.Error(err.Error(), "Error unmarshaling stonebank access token response")
		return "", err
	}

	return r.AccessToken, nil
}

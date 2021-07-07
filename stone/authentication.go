package stone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

var (
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
	issuerBank       = "stone"
	BadRequestError  = "status code 400"
	DocumentNotFound = "mongo: no documents in result"
)

type AuthResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenCreatedAt int    `json:"refresh_expires_in"`
	TokenType             string `json:"token_type"`
	NotBeforePolicy       int    `json:"not-before-policy"`
	SessionState          string `json:"session_state"`
	Scope                 string `json:"scope"`
}

func authenticate(clientID string) (string, error) {
	tk, err := fetchTokenFromStorage(clientID)
	if err != nil {
		return "", err
	}

	if tk != "" {
		return tk, nil
	}

	return authenticateAndSaveToken(clientID)
}

func authenticateAndSaveToken(clientID string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	tk, err := AuthenticationWithRetryOnBadRequest()
	if err != nil {
		return "", err
	}

	token := models.NewToken(clientID, issuerBank, tk)
	db.SaveToken(token)

	return tk, nil
}

func fetchTokenFromStorage(clientID string) (string, error) {
	token, err := db.GetTokenByClientIDAndIssuerBank(clientID, issuerBank)
	if err != nil {
		if err.Error() == DocumentNotFound {
			return "", nil
		}
		return "", err
	}

	return token.AccessToken, nil
}

// AuthenticationWithRetryOnBadRequest encapsulates logic for retry access token request once again
// in bad request status code. That's because duplicated jti returns this mencioned status code
func AuthenticationWithRetryOnBadRequest() (string, error) {
	var tk string
	var err error

	if tk, err = doAuthentication(); err != nil {
		if !strings.Contains(err.Error(), BadRequestError) {
			return "", err
		}
		return doAuthentication()
	}

	return tk, nil
}

func doAuthentication() (string, error) {
	jwt, err := generateJWT()
	if err != nil {
		return "", err
	}

	AccessTokenPayload["client_assertion"] = jwt
	AccessTokenPayload["client_id"] = config.Get().StoneClientID
	resp, err := HttpClient.PostFormURLEncoded(config.Get().URLStoneToken, AccessTokenPayload)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("stone authentication returns status code %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r AuthResponse
	err = json.Unmarshal(responseBody, &r)
	if err != nil {
		return "", err
	}

	return r.AccessToken, nil
}

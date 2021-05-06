package middleware

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/mundipagg/boleto-api/usermanagement"
	"github.com/stretchr/testify/assert"
)

func Test_Authentication_WhenEmptyCredentials_ReturnUnauthorized(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/authentication", Authentication)
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer([]byte(`{"type":"without_credentials"}`)))

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP401","message":"Unauthorized"}]}`, w.Body.String())
}

func Test_Authentication_WhenInvalidCredentials_ReturnUnauthorized(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/authentication", Authentication)
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer([]byte(`{"type":"invalid_credentials"}`)))
	req.SetBasicAuth("invalid_user", "invalid_pass")

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP401","message":"Unauthorized"}]}`, w.Body.String())
}

func Test_Authentication_WhenInvalidPassword_ReturnUnauthorized(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/authentication", Authentication)
	user, _ := usermanagement.LoadMockUserCredentials()
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer([]byte(`{"type":"valid_credentials"}`)))
	req.SetBasicAuth(user, "invalid pass")

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP401","message":"Unauthorized"}]}`, w.Body.String())
}

func Test_Authentication_WhenValidCredentials_AuthorizedRequestSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/authentication", Authentication)
	user, pass := usermanagement.LoadMockUserCredentials()
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer([]byte(`{"type":"valid_credentials"}`)))
	req.SetBasicAuth(user, pass)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
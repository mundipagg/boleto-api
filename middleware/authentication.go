package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/usermanagement"
)

const (
	serviceUserKey = "serviceuser"
)

//Authentication Middleware de autenticação para registro de boleto
func Authentication(c *gin.Context) {

	cred := getHeaderCredentials(c)

	if cred == nil {
		c.AbortWithStatusJSON(401, models.GetBoletoResponseError("MP401", "Unauthorized"))
		return
	}

	if !hasValidCredentials(cred) {
		c.AbortWithStatusJSON(401, models.GetBoletoResponseError("MP401", "Unauthorized"))
		return
	}

	c.Set(serviceUserKey, cred.Username)
}

func hasValidCredentials(c *models.Credentials) bool {
	u, hasUser := usermanagement.GetUser(c.UserKey)

	if !hasUser {
		return false
	}

	user := u.(models.Credentials)

	if user.UserKey == c.UserKey && user.Password == c.Password {
		c.Username = user.Username
		return true
	}

	return false
}

func getHeaderCredentials(c *gin.Context) *models.Credentials {
	userkey, pass, hasAuth := c.Request.BasicAuth()
	if userkey == "" || pass == "" || !hasAuth {
		return nil
	}
	return models.NewCredentials(userkey, pass)
}

package mock

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func getTokenPefisa(c *gin.Context) {
	b, _ := ioutil.ReadAll(c.Request.Body)
	const tok = `{
		"access_token": "9a7a4813-63e5-4d9a-95cc-30e3800de95e",
		"token_type": "bearer",
		"expires_in": 86399,
		"scope": "app"
	}`

	const tokError = `
	{
		"error": "unauthorized",
		"error_description": "SD-756: clientId e/ou secret inv&iquest;lido(s)"
	}`

	if strings.Contains(string(b), `grant_type=client_credentials`) {
		c.Data(200, "application/json", []byte(tok))
	} else {

		c.Data(401, "application/json", []byte(tokError))
	}

}

func registerPefisa(c *gin.Context) {
	b, _ := ioutil.ReadAll(c.Request.Body)
	const resp = `
	{
		"data": {
			"codigoBarras": "17496772400001300000000000002670000000101364",
			"linhaDigitavel": "17490.00004   00002.670008   00001.013648   6   77240000130000",
			"idTitulo": 382565
		}
	}`

	const respError = `
	{
		"error" : [ {
		  "code" : "COB-2344",
		  "message" : "Inser¿¿o do T¿tulo negada, pois a data de emiss¿o(05/12/2018@DataProximoDiaUtil=22/11/2018) ¿ posterior a D+1(@#DataProximoDiaUtil@#) da data de refer¿ncia do sistema.",
		  "action" : "Verificar a data de emiss¿o do titulo."
		} ]
	}`

	if strings.Contains(string(b), `"valorTitulo": "200"`) {
		c.Data(200, "application/json", []byte(resp))
	} else {
		c.Data(400, "application/json", []byte(respError))
	}

}

package mock

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func getTokenPfisa(c *gin.Context) {
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

func registerPfisa(c *gin.Context) {
	b, _ := ioutil.ReadAll(c.Request.Body)
	const resp = `{
		"data": {
			"codigoBarras": "17496772400001300000000000002670000000101364",
			"linhaDigitavel": "17490.00004   00002.670008   00001.013648   6   77240000130000",
			"idTitulo": 382565
		}`

	if strings.Contains(string(b), `"valor_cobrado": "0000000000000200"`) {
		c.Data(200, "text/json", []byte(resp))
	} else {
		c.Data(400, "text/json", []byte(`
			{
				"error": [
					{
						"code": "COB-1757",
						"message": "Valor incorreto ao inserir titulo controlado. Seu Numero: 00000000025"
					}
				]
			}
		`))
	}

}

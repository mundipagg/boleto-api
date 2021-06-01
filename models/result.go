package models

import "github.com/gin-gonic/gin"

//GetBoletoResult Centraliza as informações da operação GetBoleto
type GetBoletoResult struct {
	Id                                string
	Format                            string
	PrivateKey                        string
	URI                               string
	BoletoSource                      string
	TotalElapsedTimeInMilliseconds    int64
	CacheElapsedTimeInMilliseconds    int64
	DatabaseElapsedTimeInMilliseconds int64
	ErrorResponse                     BoletoResponse
	LogSeverity                       string
}

func NewGetBoletoResult(c *gin.Context) *GetBoletoResult {
	g := new(GetBoletoResult)
	g.Id = c.Query("id")
	g.Format = c.Query("fmt")
	g.PrivateKey = c.Query("pk")
	g.URI = c.Request.RequestURI
	g.BoletoSource = "none"
	return g
}

//HasValidKeys Verifica se as chaves básicas para buscar um boleto estão presentes
func (g *GetBoletoResult) HasValidKeys() bool {
	return g.Id != "" && g.PrivateKey != ""
}

//SetErrorResponse Insere as informações de erro para resposta
func (g *GetBoletoResult) SetErrorResponse(c *gin.Context, err ErrorResponse, statusCode int) {
	g.ErrorResponse = BoletoResponse{
		Errors: NewErrors(),
	}
	g.ErrorResponse.Errors.Append(err.Code, err.Message)

	if statusCode > 499 {
		c.JSON(statusCode, ErrorResponseToClient())
	} else {
		c.JSON(statusCode, g.ErrorResponse)
	}
}

func ErrorResponseToClient() BoletoResponse {
	resp := BoletoResponse{
		Errors: NewErrors(),
	}
	resp.Errors.Append("MP500", "Internal Error")
	return resp
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/models"
)

//ValidateRegisterV1 Middleware de validação das requisições de registro de boleto na rota V1
func ValidateRegisterV1(c *gin.Context) {
	rules := getBoletoFromContext(c).Title.Rules

	if rules != nil {
		c.AbortWithStatusJSON(400, models.NewSingleErrorCollection("MP400", "title.rules not available in this version"))
		return
	}
}

//ValidateRegisterV2 Middleware de validação das requisições de registro de boleto na rota V2
func ValidateRegisterV2(c *gin.Context) {
	r := getBoletoFromContext(c).Title.Rules
	bn := getBankFromContext(c).GetBankNumber()

	if r != nil && bn != models.Caixa {
		c.AbortWithStatusJSON(400, models.NewSingleErrorCollection("MP400", "title.rules not available for this bank"))
		return
	}
}

func getBoletoFromContext(c *gin.Context) models.BoletoRequest {
	var exists bool
	var boleto interface{}
	if boleto, exists = c.Get(boletoKey); exists {
		return boleto.(models.BoletoRequest)
	}
	return models.BoletoRequest{}
}

func getBankFromContext(c *gin.Context) bank.Bank {
	var exists bool
	var banking interface{}
	if banking, exists = c.Get(bankKey); exists {
		return banking.(bank.Bank)
	}
	return nil
}

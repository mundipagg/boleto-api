package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/models"
)

//validateRegisterV1 Middleware de validação das requisições de registro de boleto na rota V1
func validateRegisterV1(c *gin.Context) {
	rules := getBoletoFromContext(c).Title.Rules
	bn := getBankFromContext(c).GetBankNumber()

	if rules != nil {
		c.AbortWithStatusJSON(400, models.NewSingleErrorCollection("MP400", "title.rules not available in this version"))
		return
	}

	if bn == models.Stone {
		c.AbortWithStatusJSON(400, models.NewSingleErrorCollection("MP400", "bank Stone not available in this version"))
		return
	}
}

//validateRegisterV2 Middleware de validação das requisições de registro de boleto na rota V2
func validateRegisterV2(c *gin.Context) {
	r := getBoletoFromContext(c).Title.Rules
	bn := getBankFromContext(c).GetBankNumber()

	if r != nil && bn != models.Caixa {
		c.AbortWithStatusJSON(400, models.NewSingleErrorCollection("MP400", "title.rules not available for this bank"))
		return
	}
}

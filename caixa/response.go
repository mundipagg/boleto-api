package caixa

import "github.com/mundipagg/boleto-api/models"

//Response focado sna integracao com a Caixa
const registerBoletoResponseCaixa = `{
    {{if (eq .returnCode "1")}}
        "Errors":[{
            "Code":"{{trim .returnCode}}",
            "Message":"{{trim .returnMessage}}"
        }]
    {{else}}
        "DigitableLine": "{{fmtDigitableLine (trim .digitableLine)}}",
        "BarCodeNumber": "{{trim .barcodeNumber}}",
        "Links": [{
            "href":"{{trim .url}}",
            "rel": "pdf",
            "method":"GET"
        }],        
        "OurNumber": "{{trim .ourNumber}}"
    {{end}}
}
`

//GetRegisterBoletoAPIResponseTmpl retorna o template de resposta para a Api
func GetRegisterBoletoAPIResponseTmpl(bank models.BankNumber) string {
	return registerBoletoResponseCaixa
}

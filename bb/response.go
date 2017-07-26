package bb

import "github.com/mundipagg/boleto-api/models"

const registerBoletoResponseBB = `{
    {{if (hasErrorTags . "errorCode")}}
        "Errors": [
            {                    
                "Code": "{{trim .errorCode}}",
                "Message": "{{trim .errorMessage}}"
            }
        ]
    {{else}}
        "DigitableLine": "{{fmtDigitableLine (trim .digitableLine)}}",
        "BarCodeNumber": "{{trim .barcodeNumber}}"
    {{end}}
}
`

//GetRegisterBoletoAPIResponseTmpl retorna o template de resposta para a Api
func GetRegisterBoletoAPIResponseTmpl(bank models.BankNumber) string {
	return registerBoletoResponseBB
}

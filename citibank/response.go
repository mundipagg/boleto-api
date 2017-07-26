package citibank

import "github.com/mundipagg/boleto-api/models"

const registerBoletoResponseCiti = `{
    {{if (hasErrorTags . "errorCode" "errorMessage" "exception")}}
        {{if (hasErrorTags . "exception")}}
            "Errors": [
                {                    
                    "Code": "{{trim .returnCode}}",
                    "Message": "{{trim .returnMessage}}"
                }
            ]
        {{else}}
            "Errors": [
                {                    
                    "Code": "{{trim .errorCode}}",
                    "Message": "{{trim .errorMessage}}"
                }
            ]
        {{end}}
        
    {{else}}
        {{if .digitableLine}}
            "DigitableLine": "{{fmtDigitableLine (trim .digitableLine)}}",
        {{end}}
        {{if .barcodeNumber}}
            "BarCodeNumber": "{{trim .barcodeNumber}}"
        {{end}}
        
    {{end}}
}
`

//GetRegisterBoletoAPIResponseTmpl retorna o template de resposta para a Api
func GetRegisterBoletoAPIResponseTmpl(bank models.BankNumber) string {
	return registerBoletoResponseCiti
}

package pfisa

const registerBoletoResponsePfisa = `{
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

const boletoResponsePfisa = `
{
	"data": {
        "codigoBarras": "{{barcodeNumber}}",
        "linhaDigitavel": "{{replace digitableLine " " ""}}",        
    }
}
`

const boletoResponseErrorPfisa = `
{
	"error": [
		{
			"code": "{{errorCode}}",
			"message": "{{errorMessage}}",				
		}
	]
}
`

func getResponsePfisa() string {
	return boletoResponsePfisa
}

func getAPIResponsePfisa() string {
	return registerBoletoResponsePfisa
}

func getResponseErrorPfisa() string {
	return boletoResponseErrorPfisa
}

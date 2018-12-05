package pefisa

const registerBoletoResponsePefisa = `
{
    {{if (hasErrorTags . "errorCode")}}
        "Errors": [
            {
                "Code": "{{trim .errorCode}}",
                "Message": "{{trim .errorMessage}}"
            }
        ]
    {{else}}
	    "DigitableLine": "{{fmtDigitableLine (replace (replace .digitableLine "." "") " " "") }}",
	    "BarCodeNumber": "{{trim .barCodeNumber}}"
    {{end}}
}
`

const boletoResponsePefisa = `
{
	"data": {
        "codigoBarras": "{{barCodeNumber}}",
		"linhaDigitavel": "{{digitableLine}}"
    }
}
`

const boletoResponseErrorPefisa = `
{
	"error": [
		{
			"code": "{{errorCode}}",
			"message": "{{errorMessage}}"
		}
	]
}
`

func getResponsePefisa() string {
	return boletoResponsePefisa
}

func getAPIResponsePefisa() string {
	return registerBoletoResponsePefisa
}

func getResponseErrorPefisa() string {
	return boletoResponseErrorPefisa
}

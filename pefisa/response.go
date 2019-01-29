package pefisa

const registerBoletoResponsePefisa = `
{
    {{if (hasErrorTags . "errorCode") | (hasErrorTags . "errorMessage")}}
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

const boletoResponseErrorPefisaArray = `
{
	"code": "{{errorCode}}",
	"message": "{{errorMessage}}"
}
`

const boletoResponseErrorPefisa = `
{
	"error": "{{errorCode}}",
	"error_description": "{{errorMessage}}"
}
`

func getResponsePefisa() string {
	return boletoResponsePefisa
}

func getAPIResponsePefisa() string {
	return registerBoletoResponsePefisa
}

func getResponseErrorPefisaArray() string {
	return boletoResponseErrorPefisaArray
}

func getResponseErrorPefisa() string {
	return boletoResponseErrorPefisa
}

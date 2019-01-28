package itau

const registerBoletoResponseItau = `{
    {{if (hasErrorTags . "errorCode") | (hasErrorTags . "errorMessage")}}
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

const boletoResponseItau = `
{		
	"codigo_barras": "{{barcodeNumber}}",
	"numero_linha_digitavel": "{{digitableLine}}"	
}
`

const boletoResponseErrorItau = `
{
    "codigo":"{{errorCode}}",
    "mensagem":"{{errorMessage}}"	
}
`

func getResponseItau() string {
	return boletoResponseItau
}

func getAPIResponseItau() string {
	return registerBoletoResponseItau
}

func getResponseErrorItau() string {
	return boletoResponseErrorItau
}

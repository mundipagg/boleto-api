package citibank

const registerBoletoResponseCiti = `
	{
		{{if eq .returnCode "0"}}
		   "DigitableLine": "{{fmtDigitableLine (trim .digitableLine)}}",
		   "BarCodeNumber": "{{trim .barcodeNumber}}"
		{{else}}
		 "Errors": [
			{
				"Code": "{{trim .returnCode}}",
				"Message": "{{trim .returnMessage}}"
			}
		 ]
		{{end}}
    }
`

func getAPIResponseCiti() string {
	return registerBoletoResponseCiti
}

package santander

var apiResponse = `
{
    {{if (eq .errorCode "20")}}
        "Errors": [
            {                    
                "Code": "{{trim .errorCode}}",
                "Message": "{{trim .message | formatSingleLine}}"
            }
        ]
    {{else}}
        "DigitableLine": "{{fmtDigitableLine (trim .digitableLine)}}",
        "BarcodeNumber": "{{trim .barcodeNumber}}"        
    {{end}}
}
`

func getAPIResponseSantander() string {
	return apiResponse
}

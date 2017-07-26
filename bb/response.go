package bb

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

func getAPIResponse() string {
	return registerBoletoResponseBB
}

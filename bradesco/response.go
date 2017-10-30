package bradesco

var apiResponse = `
{
	{{if eq .returnCode "0"}}
       "DigitableLine": "{{.digitableLine}}",
		"Links": [{
			"href":"{{.url}}",
			"rel": "pdf",
			"method":"GET"
		}]
    {{else}}
     "Errors": [
		{
			"Code": "{{.returnCode}}",
			"Message": "{{.returnMessage}}"
		}
        ]
    {{end}}
}
`

func getAPIResponseBradesco() string {
	return apiResponse
}

package stone

const templateRequest = `
## Content-Type:application/json
## Authorization:Bearer {{.Authentication.AuthorizationToken}}
{
    "account_id": "{{.Authentication.AccessKey}}",
    "amount": {{.Title.AmountInCents}},
    "expiration_date": "{{.Title.ExpireDate}}",
    "invoice_type": "{{.Title.BoletoTypeCode}}",
    "customer": {
        "document": "{{.Buyer.Document.Number}}",
        "legal_name": "{{.Buyer.Name}}",
	{{if eq .Buyer.Document.Type "CNPJ"}}
        "trade_name": "{{.Buyer.Name}}"
	{{else}}
		"trade_name": null
	{{end}}
    }
}`

const templateResponse = `
{
    "barcode": "{{barCodeNumber}}",
    "our_number": "{{ourNumber}}",
    "writable_line": "{{digitableLine}}"
}
`

const templateError = `
{
    "reason": "{{messageError}}",
    "type": "{{errorCode}}"
}
`

const templateAPI = `
{
    {{if (hasErrorTags . "errorCode") | (hasErrorTags . "messageError")}}
    "Errors": [
        {
        {{if (hasErrorTags . "errorCode")}}
            "Code": "{{trim .errorCode}}",
        {{end}}
        {{if (eq .messageError "{}")}}
            "Message": "{{trim .errorCode}}"
        {{else}}
            "Message": "{{trim .messageError}}"
        {{end}}
        }
    ]
    {{else}}
        "DigitableLine": "{{fmtDigitableLine .digitableLine}}",
        "BarCodeNumber": "{{.barCodeNumber}}",
        "OurNumber": "{{.ourNumber}}"
    {{end}}
}
`

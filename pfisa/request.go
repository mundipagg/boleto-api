package pfisa

const registerPfisa = `

## Authorization: Bearer {{.Authentication.AuthorizationToken}}
## Content-Type: application/json



{
    "idBeneficiario": "{{.Agreement.AgreementNumber}}",
    "carteira": "{{.Agreement.Wallet}}",
    "nossoNumero": "{{padLeft .Title.OurNumber 20}}",
    "seuNumero": "{{truncate .Title.DocumentNumber 10}}",
    "usoEmpresa": "A",
    "tipoTitulo": 1,
    "valorTitulo": "{{.Title.AmountInCents}}",
    "dataDocumento": "{{enDate (today) "-"}}",
    "dataVencimento": "{{enDate .Title.ExpireDate "-"}}",  
    "emitente": {
        "nome": "{{.Recipient.Name}}",
        {{if (eq .Recipient.Document.Type "CNPJ")}}
        "tipo": "J",        
        {{else}}
        "tipo": "F",
        {{end}}        
        "cnpjCpf": "{{extractNumbers .Recipient.Document.Number}}",
        "endereco": "{{truncate .Recipient.Address.Street 40}}",
        "cidade": "{{truncate .Recipient.Address.City 60}}",
        "cep": "{{truncate .Recipient.Address.ZipCode 8}}",
        "uf": "{{truncate .Recipient.Address.StateCode 2}}",
        "bairro": "{{truncate .Recipient.Address.District 65}}",        
    },
    "pagador": {
        "nome": "{{truncate .Buyer.Name 40}}",
        {{if (eq .Recipient.Document.Type "CNPJ")}}
        "tipo": "J",
        {{else}}
        "tipo": "F",
        {{end}}        
        "cnpjCpf": "{{extractNumbers .Recipient.Document.Number}}",
        "endereco": "Avenida Afonso Arinos de Melo Franco",
        "cidade": "{{truncate .Buyer.Address.City 20}}",
        "cep": "{{truncate (extractNumbers .Buyer.Address.ZipCode) 8",
        "uf": "{{truncate .Buyer.Address.StateCode 2}}",
        "bairro": "{{truncate .Buyer.Address.District 65}}",
        
    },
    "mensagens": [
        "{{truncate title.instructions 80}}",
    ]
}
`

const pfisaGetTokenRequest = `
## Authorization:Basic {{base64 (concat .Authentication.Username ":" .Authentication.Password)}}
## Content-Type: application/x-www-form-urlencoded
grant_type=client_credentials`

const tokenResponse = `{	
	"access_token": "{{access_token}}",
}`

const tokenErrorResponse = `{    
	"error_description": "{{errorMessage}}",
}`

func getRequestToken() string {
	return pfisaGetTokenRequest
}

func getTokenResponse() string {
	return tokenResponse
}

func getTokenErrorResponse() string {
	return tokenErrorResponse
}

func getRequestPfisa() string {
	return registerPfisa
}

package bradescoNetEmpresa

const registerBradescoNetEmpresa = `
## Content-Type:application/json
{
    {{if (eq .Recipient.Document.Type "CNPJ")}}
        "nuCPFCNPJ": "{{splitValues (extractNumbers .Recipient.Document.Number) 0 8}}",    
    {{else}}
         "nuCPFCNPJ": "{{splitValues (extractNumbers .Recipient.Document.Number) 0 9}}",	
	{{end}}
    
	{{if (eq .Recipient.Document.Type "CNPJ")}}
        "filialCPFCNPJ": "{{splitValues (extractNumbers .Recipient.Document.Number) 8 12}}",    
    {{else}}
            "filialCPFCNPJ": "0",	
	{{end}}
	
    {{if (eq .Recipient.Document.Type "CNPJ")}}
        "ctrlCPFCNPJ": "{{splitValues (extractNumbers .Recipient.Document.Number) 12 14}}",
    {{else}}
        "ctrlCPFCNPJ": "{{splitValues (extractNumbers .Recipient.Document.Number) 9 11}}",	
    {{end}}	
    "cdTipoAcesso": "2",
    "clubBanco": "2269651",
    "cdTipoContrato": "48",    
    "idProduto": "{{padLeft (toString16 .Agreement.Wallet) "0" 2}}",
    "nuNegociacao": "{{.Agreement.Agency}}0000000{{.Agreement.Account}}",
    "cdBanco": "237",    
    "tpRegistro": "1",    
    "nuTitulo": "{{.Title.OurNumber}}",
    "nuCliente": "{{.Title.DocumentNumber}}",
	"dtEmissaoTitulo": "{{brDateDelimiterTime .Title.CreateDate "."}}",
    "dtVencimentoTitulo": "{{brDateDelimiter .Title.ExpireDate "."}}",
    "tpVencimento": "0",
    "vlNominalTitulo": "{{.Title.AmountInCents}}",
    "cdEspecieTitulo": "{{ .Title.BoletoTypeCode}}",
    "nomePagador": "{{truncate .Buyer.Name 70}}",
    "logradouroPagador": "{{truncate .Buyer.Address.Street 40}}",
    "nuLogradouroPagador": "{{truncate .Buyer.Address.Number 10}}",
    "complementoLogradouroPagador": "{{truncate .Buyer.Address.Complement 15}}",
	{{ $length := len (extractNumbers .Buyer.Address.ZipCode)}}
	{{ $zipCode := (extractNumbers .Buyer.Address.ZipCode)}}
	{{ if ge $length 5}}
		"cepPagador": "{{splitValues $zipCode 0 5}}",
	{{ end }}

	{{ if ge $length 8}}
		"complementoCepPagador": "{{splitValues $zipCode 5 8}}",
	{{ end }}

    "bairroPagador": "{{truncate .Buyer.Address.District 40}}",
    "municipioPagador": "{{truncate .Buyer.Address.City 30}}",
    "ufPagador": "{{truncate .Buyer.Address.StateCode 2}}",
    {{if (eq .Buyer.Document.Type "CPF")}}
    	"cdIndCpfcnpjPagador": "1",
    {{else}}
        "cdIndCpfcnpjPagador": "2",
    {{end}}
    "nuCpfcnpjPagador": "{{extractNumbers .Buyer.Document.Number}}",
    "endEletronicoPagador": "{{truncate .Buyer.Email 70}}",    
}
`

const reponseBradescoNetEmpresaXml = `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <ns2:registrarTituloResponse xmlns:ns2="http://ws.registrotitulo.ibpj.web.bradesco.com.br/">
			<return>{{contentJson}}</return>
        </ns2:registrarTituloResponse>
    </soapenv:Body>
</soapenv:Envelope>
`

const responseBradescoNetEmpresaJson = `{{.contentJson}}`

const reponseBradescoNetEmpresa = `
{
    "cdErro": "{{returnCode}}",
    "msgErro": "{{returnMessage}}",    
    "linhaDigitavel": "{{digitableLine}}"
}
`

func getRequestBradescoNetEmpresa() string {
	return registerBradescoNetEmpresa
}

func getResponseBradescoNetEmpresaXml() string {
	return reponseBradescoNetEmpresaXml
}

func getResponseBradescoNetEmpresaJson() string {
	return responseBradescoNetEmpresaJson
}

func getResponseBradescoNetEmpresa() string {
	return reponseBradescoNetEmpresa
}

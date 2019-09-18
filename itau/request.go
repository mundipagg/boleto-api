package itau

const registerItau = `
## Accept: application/vnd.itau
## access_token: {{.Authentication.AuthorizationToken}}
## itau-chave: {{.Authentication.AccessKey}}
## identificador: {{.Recipient.Document.Number}}
## Content-Type: application/json

{
    "tipo_ambiente": {{itauEnv}},
    "tipo_registro": 1,
    "tipo_cobranca": 1,
    "tipo_produto": "00006",
    "subproduto": "00008",
    "beneficiario": {
        "cpf_cnpj_beneficiario": "{{extractNumbers .Recipient.Document.Number}}",
        "agencia_beneficiario": "{{padLeft .Agreement.Agency "0" 4}}",
        "conta_beneficiario": "{{padLeft .Agreement.Account "0" 7}}",
        "digito_verificador_conta_beneficiario": "{{.Agreement.AccountDigit}}"
    },
    "identificador_titulo_empresa": "{{unescapeHtmlString (truncate .Recipient.Name 25)}}",
    "uso_banco": "",
    "titulo_aceite": "S",
    "pagador": {
        "cpf_cnpj_pagador": "{{extractNumbers .Buyer.Document.Number}}",
        "nome_pagador": "{{unescapeHtmlString (truncate .Buyer.Name 30)}}",
        "logradouro_pagador": "{{unescapeHtmlString (truncate (concat .Buyer.Address.Street " " .Buyer.Address.Number " " .Buyer.Address.Complement) 40) }}",        
        "bairro_pagador": "{{unescapeHtmlString (truncate .Buyer.Address.District 15)}}",
        "cidade_pagador": "{{unescapeHtmlString (truncate .Buyer.Address.City 20)}}",
        "uf_pagador": "{{truncate .Buyer.Address.StateCode 2}}",
        "cep_pagador": "{{truncate (extractNumbers .Buyer.Address.ZipCode) 8}}",
        "grupo_email_pagador": [
            {
                "email_pagador": ""
            }
        ]
    },
    "tipo_carteira_titulo": "{{.Agreement.Wallet}}",
    "moeda": {
        "codigo_moeda_cnab": "09",
        "quantidade_moeda": ""
    },
    "nosso_numero": "{{padLeft (toString .Title.OurNumber) "0" 8}}",
    "digito_verificador_nosso_numero": "{{mod10ItauDv (padLeft (toString .Title.OurNumber) "0" 8) (padLeft .Agreement.Agency "0" 4) (padLeft .Agreement.Account "0" 7) .Agreement.Wallet}}",
    "codigo_barras": "",
    "data_vencimento": "{{enDate .Title.ExpireDateTime "-"}}",
    "valor_cobrado": "{{padLeft (toString64 .Title.AmountInCents) "0" 16}}",
    "seu_numero": "{{padLeft .Title.DocumentNumber "0" 10}}",
    "especie": "{{ .Title.BoletoTypeCode}}",
    "data_emissao": "{{enDate (today) "-"}}",    
    "tipo_pagamento": 1,
    "indicador_pagamento_parcial": "false",
    "quantidade_pagamento_parcial": "0",
    "quantidade_parcelas": "0",
    "instrucao_cobranca_1": "",
    "quantidade_dias_1": "",
    "data_instrucao_1": "",
    "instrucao_cobranca_2": "",
    "quantidade_dias_2": "",
    "data_instrucao_2": "",
    "instrucao_cobranca_3": "",
    "quantidade_dias_3": "",
    "data_instrucao_3": "",
    "valor_abatimento": "",
    "juros": {
        "data_juros": "",
        "tipo_juros": 5,
        "valor_juros": "",
        "percentual_juros": ""
    },
    "multa": {
        "data_multa": "",
        "tipo_multa": 3,
        "valor_multa": "",
        "percentual_multa": ""
    },    
    "grupo_desconto": [{
        "data_desconto": "",
        "tipo_desconto": 0,
        "valor_desconto": "",
        "percentual_desconto": ""
    }],    
    "recebimento_divergente": {
        "tipo_autorizacao_recebimento": "3",
        "tipo_valor_percentual_recebimento": "",
        "valor_minimo_recebimento": "",
        "percentual_minimo_recebimento": "",
        "valor_maximo_recebimento": "",
        "percentual_maximo_recebimento": ""
    },
    "grupo_rateio": []
}

`

const itauGetTicketRequest = `## Authorization:Basic {{base64 (concat .Authentication.Username ":" .Authentication.Password)}}
## Content-Type: application/x-www-form-urlencoded
scope=readonly&grant_type=client_credentials&clientId={{.Authentication.Username}}&clientSecret={{.Authentication.Password}}`

const ticketResponse = `{
    "codigo":"{{errorCode}}",
    "mensagem":"{{errorMessage}}",
    "access_token": "{{access_token}}",
    "Message":"{{errorMessage}}"
}`

const ticketErrorResponse = `{
    "Message":"{{errorMessage}}"
}`

func getRequestTicket() string {
	return itauGetTicketRequest
}

func getTicketResponse() string {
	return ticketResponse
}

func getTicketErrorResponse() string {
	return ticketErrorResponse
}

func getRequestItau() string {
	return registerItau
}

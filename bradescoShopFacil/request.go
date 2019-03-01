package bradescoShopFacil

const registerBradescoShopFacil = `
## Content-Type:application/json
## Authorization:Basic {{base64 (concat .Authentication.Username ":" .Authentication.Password)}}
{
    "merchant_id": "{{.Authentication.Username}}",
    "meio_pagamento": "300",
    "pedido": {
        "numero": "{{escape .Title.DocumentNumber}}",
        "valor": {{.Title.AmountInCents}},
        "descricao": ""
    },
    "comprador": {
        "nome": "{{escape .Buyer.Name}}",
        "documento": "{{escape .Buyer.Document.Number}}",
        "endereco": {
            "cep": "{{extractNumbers .Buyer.Address.ZipCode}}",
            "logradouro": "{{escape .Buyer.Address.Street}}",
            "numero": "{{escape .Buyer.Address.Number}}",
            "complemento": "{{escape .Buyer.Address.Complement}}",
            "bairro": "{{escape .Buyer.Address.District}}",
            "cidade": "{{escape .Buyer.Address.City}}",
            "uf": "{{escape .Buyer.Address.StateCode}}"
        },
        "ip": "",
        "user_agent": ""
    },
    "boleto": {
        "beneficiario": "{{escape .Recipient.Name}}",
        "carteira": "{{.Agreement.Wallet}}",
        "nosso_numero": "{{padLeft (toString .Title.OurNumber) "0" 11}}",
        "data_emissao": "{{enDate today "-"}}",
        "data_vencimento": "{{enDate .Title.ExpireDateTime "-"}}",
        "valor_titulo": {{.Title.AmountInCents}},
        "url_logotipo": "",
        "mensagem_cabecalho": "",
        "tipo_renderizacao": "1",
        "instrucoes": {
            "instrucao_linha_1": "{{escape .Title.Instructions}}"
        },
        "registro": {
            "agencia_pagador": "",
            "razao_conta_pagador": "",
            "conta_pagador": "",
            "controle_participante": "",
            "aplicar_multa": false,
            "valor_percentual_multa": 0,
            "valor_desconto_bonificacao": 0,
            "debito_automatico": false,
            "rateio_credito": false,
            "endereco_debito_automatico": "2",
            "tipo_ocorrencia": "02",
            "especie_titulo": "01",
            "primeira_instrucao": "00",
            "segunda_instrucao": "00",
            "valor_juros_mora": 0,
            "data_limite_concessao_desconto": null,
            "valor_desconto": 0,
            "valor_iof": 0,
            "valor_abatimento": 0,
            {{if (eq .Buyer.Document.Type "CPF")}}
            	"tipo_inscricao_pagador": "01",
			{{else}}
            	"tipo_inscricao_pagador": "02",
			{{end}}
            "sequencia_registro": ""
        }
    },
    "token_request_confirmacao_pagamento": ""
}
`

const responseBradescoShopFacil = `
{
    "boleto": {
        "linha_digitavel": "{{digitableLine}}",
        "url_acesso": "{{url}}"
    },
    "status": {
        "codigo": "{{returnCode}}",
        "mensagem": "{{returnMessage}}"
    }
}
`

func getRequestBradescoShopFacil() string {
	return registerBradescoShopFacil
}

func getResponseBradescoShopFacil() string {
	return responseBradescoShopFacil
}

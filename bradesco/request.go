package bradesco

const registerBradesco = `
## Content-Type:application/json
## Authorization:Basic {{base64 (concat .Authentication.Username ":" .Authentication.Password)}}
{
    "merchant_id": "{{.Authentication.Username}}",
    "meio_pagamento": "300",
    "pedido": {
        "numero": "{{.Title.DocumentNumber}}",
        "valor": {{.Title.AmountInCents}},
        "descricao": ""
    },
    "comprador": {
        "nome": "{{.Buyer.Name}}",
        "documento": "{{.Buyer.Document.Number}}",
        "endereco": {
            "cep": "{{.Buyer.Address.ZipCode}}",
            "logradouro": "{{.Buyer.Address.Street}}",
            "numero": "{{.Buyer.Address.Number}}",
            "complemento": "{{.Buyer.Address.Complement}}",
            "bairro": "{{.Buyer.Address.District}}",
            "cidade": "{{.Buyer.Address.City}}",
            "uf": "{{.Buyer.Address.StateCode}}"
        },
        "ip": "",
        "user_agent": ""
    },
    "boleto": {
        "beneficiario": "{{.Recipient.Name}}",
        "carteira": "{{.Agreement.Wallet}}",
        "nosso_numero": "{{padLeft (toString .Title.OurNumber) "0" 11}}",
        "data_emissao": "{{enDate today "-"}}",
        "data_vencimento": "{{enDate .Title.ExpireDateTime "-"}}",
        "valor_titulo": {{.Title.AmountInCents}},
        "url_logotipo": "",
        "mensagem_cabecalho": "",
        "tipo_renderizacao": "1",
        "instrucoes": {
            "instrucao_linha_1": "{{.Title.Instructions}}"
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

const responseBradesco = `
{
    "boleto": {
        "linha_digitavel_formatada": "{{digitableLine}}",
        "url_acesso": "{{url}}"
    },
    "status": {
        "codigo": "{{returnCode}}",
        "mensagem": "{{returnMessage}}"
    }
}
`

func getRequestBradesco() string {
	return registerBradesco
}

func getResponseBradesco() string {
	return responseBradesco
}

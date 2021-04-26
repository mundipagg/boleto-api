package caixa

var expectedBasicTitleRequestFields = []string{
	`<CODIGO_BENEFICIARIO>0123456</CODIGO_BENEFICIARIO>`,
	`<NOSSO_NUMERO>12345678901234</NOSSO_NUMERO>`,
	`<VALOR>2.00</VALOR>`,
	`<TIPO_ESPECIE>99</TIPO_ESPECIE>`,
	`<MENSAGEM>Campo de instrucoes -  max 40 caracteres</MENSAGEM>`,
}

var expectedBuyerRequestFields = []string{
	`<NOME>Willian Jadson Bezerra Menezes Tupinamba</NOME>`,
	`<LOGRADOURO>Rua da Assuncao de Sa 123 Secao A, s 02</LOGRADOURO>`,
	`<BAIRRO>Acai</BAIRRO>`,
	`<CIDADE>Belem do Para</CIDADE>`,
	`<UF>PA</UF>`,
	`<CEP>20520051</CEP>`,
}

var expectedStrictRulesFieldsV2 = []string{
	`<PAGAMENTO>`,
	`<TIPO>NAO_ACEITA_VALOR_DIVERGENTE</TIPO>`,
	`</PAGAMENTO>`,
}

var expectedFlexRulesFieldsV2 = []string{
	`<PAGAMENTO>`,
	`<TIPO>ACEITA_QUALQUER_VALOR</TIPO>`,
	`<NUMERO_DIAS>60</NUMERO_DIAS>`,
	`</PAGAMENTO>`,
}

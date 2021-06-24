package caixa

const responseFromCaixa = `
<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <manutencaocobrancabancaria:SERVICO_SAIDA xmlns:manutencaocobrancabancaria="http://caixa.gov.br/sibar/manutencao_cobranca_bancaria/boleto/externo" xmlns:sibar_base="http://caixa.gov.br/sibar">
            <sibar_base:HEADER>
                <OPERACAO>{{operation}}</OPERACAO>
                <DATA_HORA>{{datetime}}</DATA_HORA>
            </sibar_base:HEADER>
            <DADOS>                
                <CONTROLE_NEGOCIAL>
                    <ORIGEM_RETORNO>SIGCB</ORIGEM_RETORNO>
                    <COD_RETORNO>{{returnCode}}</COD_RETORNO>
                    <MENSAGENS>
                        <RETORNO>{{returnMessage}}</RETORNO>
                    </MENSAGENS>
                </CONTROLE_NEGOCIAL>
                <INCLUI_BOLETO>
                    <EXCECAO>{{exception}}</EXCECAO>
                    <CODIGO_BARRAS>{{barcodeNumber}}</CODIGO_BARRAS>
                    <LINHA_DIGITAVEL>{{digitableLine}}</LINHA_DIGITAVEL>
                    <NOSSO_NUMERO>{{ourNumber}}</NOSSO_NUMERO>
                    <URL>{{url}}</URL>
                </INCLUI_BOLETO>
            </DADOS>
        </manutencaocobrancabancaria:SERVICO_SAIDA>
    </soapenv:Body>
</soapenv:Envelope>
`

const requestToCaixa = `

## SOAPAction:IncluiBoleto
## Content-Type:text/xml

<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ext="http://caixa.gov.br/sibar/manutencao_cobranca_bancaria/boleto/externo" xmlns:sib="http://caixa.gov.br/sibar">
<soapenv:Body>
<ext:SERVICO_ENTRADA >
         <sib:HEADER>
            <VERSAO>1.0</VERSAO>
            <AUTENTICACAO>{{unscape .Authentication.AuthorizationToken}}</AUTENTICACAO>
            <USUARIO_SERVICO>{{caixaEnv}}</USUARIO_SERVICO>
            <OPERACAO>INCLUI_BOLETO</OPERACAO>
            <SISTEMA_ORIGEM>SIGCB</SISTEMA_ORIGEM>
            <UNIDADE>{{.Agreement.Agency}}</UNIDADE>
            <DATA_HORA>{{fullDate today}}</DATA_HORA>
            </sib:HEADER>
         <DADOS>
            <INCLUI_BOLETO>
              <CODIGO_BENEFICIARIO>{{padLeft (toString .Agreement.AgreementNumber) "0" 7}}</CODIGO_BENEFICIARIO>
               <TITULO>
                  <NOSSO_NUMERO>{{toString .Title.OurNumber}}</NOSSO_NUMERO>
                  <NUMERO_DOCUMENTO>{{truncateOnly .Title.DocumentNumber 11}}</NUMERO_DOCUMENTO>
                  <DATA_VENCIMENTO>{{enDate .Title.ExpireDateTime "-"}}</DATA_VENCIMENTO>
                  <VALOR>{{toFloatStr .Title.AmountInCents}}</VALOR>
                  <TIPO_ESPECIE>{{.Title.BoletoTypeCode}}</TIPO_ESPECIE>
                  <FLAG_ACEITE>S</FLAG_ACEITE>
                  <DATA_EMISSAO>{{enDate today "-"}}</DATA_EMISSAO>
                  <JUROS_MORA>
                     <TIPO>ISENTO</TIPO>
                     <VALOR>0</VALOR>
                  </JUROS_MORA>
                  <VALOR_ABATIMENTO>0</VALOR_ABATIMENTO>
                  <POS_VENCIMENTO>
                     <ACAO>DEVOLVER</ACAO>
                  {{if .Title.HasRules}}       
                     <NUMERO_DIAS>{{.Title.Rules.MaxDaysToPayPastDue}}</NUMERO_DIAS>
                  {{else}}
                     <NUMERO_DIAS>0</NUMERO_DIAS>
                  {{end}}
                  </POS_VENCIMENTO>
                  <CODIGO_MOEDA>9</CODIGO_MOEDA>
                  <PAGADOR>
                  {{if eq .Buyer.Document.Type "CPF"}}
					 	   <CPF>{{.Buyer.Document.Number}}</CPF>
                     <NOME>{{truncateOnly (clearStringCaixa .Buyer.Name) 40}}</NOME>
                  {{else}}
					 	   <CNPJ>{{.Buyer.Document.Number}}</CNPJ>
                     <RAZAO_SOCIAL>{{truncateOnly (clearStringCaixa .Buyer.Name) 40}}</RAZAO_SOCIAL>
					   {{end}}
                     <ENDERECO>
                        <LOGRADOURO>{{truncateOnly (joinSpace (clearStringCaixa .Buyer.Address.Street) (clearStringCaixa .Buyer.Address.Number) (clearStringCaixa .Buyer.Address.Complement)) 40}}</LOGRADOURO>
                        <BAIRRO>{{truncateOnly (clearStringCaixa .Buyer.Address.District) 15}}</BAIRRO>
                        <CIDADE>{{truncateOnly (clearStringCaixa .Buyer.Address.City) 15}}</CIDADE>
                        <UF>{{truncateOnly (clearStringCaixa .Buyer.Address.StateCode) 2}}</UF>
                        <CEP>{{truncateOnly (clearStringCaixa .Buyer.Address.ZipCode) 8}}</CEP>
                     </ENDERECO>
                  </PAGADOR>
                  <FICHA_COMPENSACAO>
                     <MENSAGENS>
                        <MENSAGEM>{{truncateOnly (clearStringCaixa .Title.Instructions) 40}}</MENSAGEM>
                     </MENSAGENS>
                  </FICHA_COMPENSACAO>
                  <RECIBO_PAGADOR>
                     <MENSAGENS>
                        <MENSAGEM>{{truncateOnly (clearStringCaixa .Title.Instructions) 40}}</MENSAGEM>
                     </MENSAGENS>
                  </RECIBO_PAGADOR>
               {{if .Title.HasRules}} 
                  <PAGAMENTO>
                     <QUANTIDADE_PERMITIDA>1</QUANTIDADE_PERMITIDA>   
                  {{if .Title.Rules.AcceptDivergentAmount}}
                     <TIPO>ACEITA_QUALQUER_VALOR</TIPO>
                  {{else}}
                     <TIPO>NAO_ACEITA_VALOR_DIVERGENTE</TIPO>
                  {{end}}
                     <VALOR_MINIMO>0.00</VALOR_MINIMO>
                     <VALOR_MAXIMO>0.00</VALOR_MAXIMO>
                  </PAGAMENTO>
               {{end}}                 
               </TITULO>
            </INCLUI_BOLETO>
           </DADOS>
      </ext:SERVICO_ENTRADA>
</soapenv:Body>
</soapenv:Envelope>
`

func getRequestCaixa() string {
	return requestToCaixa
}

func getResponseCaixa() string {
	return responseFromCaixa
}

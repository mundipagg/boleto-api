package caixa

const responseCaixa = `
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

const incluiBoleto = `

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
                  <NUMERO_DOCUMENTO>{{truncate .Title.DocumentNumber 11}}</NUMERO_DOCUMENTO>
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
                    <NUMERO_DIAS>0</NUMERO_DIAS>
                  </POS_VENCIMENTO>
                  <CODIGO_MOEDA>9</CODIGO_MOEDA>
                  <PAGADOR>
                     {{if eq .Buyer.Document.Type "CPF"}}
					 	<CPF>{{.Buyer.Document.Number}}</CPF>
                     	<NOME>{{truncate .Buyer.Name 40}}</NOME>
                     {{else}}
					 	<CNPJ>{{.Buyer.Document.Number}}</CNPJ>
                     	<RAZAO_SOCIAL>{{truncate .Buyer.Name 40}}</RAZAO_SOCIAL>
					 {{end}}
                     <ENDERECO>
                     <LOGRADOURO>{{clearString (truncateManyFields 40 .Buyer.Address.Street .Buyer.Address.Number .Buyer.Address.Complement)}}</LOGRADOURO>
                        <BAIRRO>{{clearString (truncate .Buyer.Address.District 15)}}</BAIRRO>
                        <CIDADE>{{clearString (truncate .Buyer.Address.City 15)}}</CIDADE>
                        <UF>{{clearString (truncate .Buyer.Address.StateCode 2)}}</UF>
                        <CEP>{{clearString (truncate .Buyer.Address.ZipCode 8)}}</CEP>
                     </ENDERECO>
                  </PAGADOR>
                  <FICHA_COMPENSACAO>
                     <MENSAGENS>
                        <MENSAGEM>{{clearString (truncate .Title.Instructions 40)}}</MENSAGEM>
                        </MENSAGENS>
                  </FICHA_COMPENSACAO>
                  <RECIBO_PAGADOR>
                     <MENSAGENS>
                        <MENSAGEM>{{clearString (truncate .Title.Instructions 40)}}</MENSAGEM>
                     </MENSAGENS>
                  </RECIBO_PAGADOR>                 
               </TITULO>
            </INCLUI_BOLETO>
           </DADOS>
      </ext:SERVICO_ENTRADA>
</soapenv:Body>
</soapenv:Envelope>
`

func getRequestCaixa() string {
	return incluiBoleto
}

func getResponseCaixa() string {
	return responseCaixa
}

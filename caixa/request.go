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

//LvWr1op5Ayibn6jsCQ3/2bW4KwThVAlLK5ftxABlq20=
//
const incluiBoleto = `

## SOAPAction:IncluiBoleto
## Content-Type:text/xml

<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ext="http://caixa.gov.br/sibar/manutencao_cobranca_bancaria/boleto/externo" xmlns:sib="http://caixa.gov.br/sibar">
<soapenv:Body>
<ext:SERVICO_ENTRADA >
         <sib:HEADER>
            <VERSAO>1.0</VERSAO>
            <AUTENTICACAO>{{unscape .Authentication.AuthorizationToken}}</AUTENTICACAO>
            <USUARIO_SERVICO>SGCBS01D</USUARIO_SERVICO>
            <OPERACAO>INCLUI_BOLETO</OPERACAO>
            <SISTEMA_ORIGEM>SIGCB</SISTEMA_ORIGEM>
            <UNIDADE>{{.Agreement.Agency}}</UNIDADE>
            <DATA_HORA>{{fullDate today}}</DATA_HORA>
            </sib:HEADER>
         <DADOS>
            <INCLUI_BOLETO>
              <CODIGO_BENEFICIARIO>{{padLeft (toString .Agreement.AgreementNumber) "0" 7}}</CODIGO_BENEFICIARIO>
               <TITULO>
                  <NOSSO_NUMERO>{{padLeft (toString .Title.OurNumber) "0"  17 }}</NOSSO_NUMERO>
                  <NUMERO_DOCUMENTO>{{.Title.DocumentNumber}}</NUMERO_DOCUMENTO>
                  <DATA_VENCIMENTO>{{enDate .Title.ExpireDateTime "-"}}</DATA_VENCIMENTO>
                  <VALOR>{{toFloatStr .Title.AmountInCents}}</VALOR>
                  <TIPO_ESPECIE>99</TIPO_ESPECIE>
                  <FLAG_ACEITE>S</FLAG_ACEITE>
                  <DATA_EMISSAO>2017-05-16</DATA_EMISSAO>
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
                     	<NOME>{{.Buyer.Name}}</NOME>
                     {{else}}
					 	<CNPJ>{{.Buyer.Document.Number}}</CNPJ>
                     	<RAZAO_SOCIAL>{{.Buyer.Name}}</RAZAO_SOCIAL>
					 {{end}}
                     <ENDERECO>
                        <LOGRADOURO>{{.Buyer.Address.Street}} {{.Buyer.Address.Number}} {{.Buyer.Address.Complement}}</LOGRADOURO>
                        <BAIRRO>{{.Buyer.Address.District}}</BAIRRO>
                        <CIDADE>{{.Buyer.Address.City}}</CIDADE>
                        <UF>{{.Buyer.Address.StateCode}}</UF>
                        <CEP>{{.Buyer.Address.ZipCode}}</CEP>
                     </ENDERECO>
                  </PAGADOR>                  
                  <FICHA_COMPENSACAO>
                     <MENSAGENS>
                        <MENSAGEM>{{.Title.Instructions}}</MENSAGEM>                        
                        </MENSAGENS>
                  </FICHA_COMPENSACAO>
                  <RECIBO_PAGADOR>
                     <MENSAGENS>
                        <MENSAGEM>{{.Title.Instructions}}</MENSAGEM>
                     </MENSAGENS>
                  </RECIBO_PAGADOR>                 
               </TITULO>
            </INCLUI_BOLETO>
           </DADOS>
      </ext:SERVICO_ENTRADA>
</soapenv:Body>
</soapenv:Envelope>
`

//GetRegisterBoletoCaixaTmpl retorna o padrão de registro de boleto da Caixa
func GetRegisterBoletoCaixaTmpl() string {
	return incluiBoleto
}

//GetResponseTemplateCaixa retorna o template de mensagem da Caixa
func GetResponseTemplateCaixa() string {
	return responseCaixa
}

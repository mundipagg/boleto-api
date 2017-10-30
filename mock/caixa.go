package mock

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerBoletoCaixa(c *gin.Context) {
	sData := `
<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Body>
      <manutencaocobrancabancaria:SERVICO_SAIDA xmlns:manutencaocobrancabancaria="http://caixa.gov.br/sibar/manutencao_cobranca_bancaria/boleto/externo" xmlns:sibar_base="http://caixa.gov.br/sibar">
         <sibar_base:HEADER>
            <VERSAO>1.0</VERSAO>
            <AUTENTICACAO>LvWr1op5Ayibn6jsCQ3/2bW4KwThVAlLK5ftxABlq20=</AUTENTICACAO>
            <USUARIO_SERVICO>SGCBS01D</USUARIO_SERVICO>
            <OPERACAO>INCLUI_BOLETO</OPERACAO>
            <SISTEMA_ORIGEM>SIGCB</SISTEMA_ORIGEM>
            <UNIDADE>1679</UNIDADE>
            <DATA_HORA>20170718150257</DATA_HORA>
         </sibar_base:HEADER>
         <COD_RETORNO>00</COD_RETORNO>
         <ORIGEM_RETORNO>MANUTENCAO_COBRANCA_BANCARIA</ORIGEM_RETORNO>
         <MSG_RETORNO />
         <DADOS>
            <CONTROLE_NEGOCIAL>
               <ORIGEM_RETORNO>SIGCB</ORIGEM_RETORNO>
               <COD_RETORNO>0</COD_RETORNO>
               <MENSAGENS>
                  <RETORNO>(0) OPERACAO EFETUADA</RETORNO>
               </MENSAGENS>
            </CONTROLE_NEGOCIAL>
            <INCLUI_BOLETO>
               <CODIGO_BARRAS>10493726700000010002006561000100040992226984</CODIGO_BARRAS>
               <LINHA_DIGITAVEL>10492006506100010004209922269841372670000001000</LINHA_DIGITAVEL>
               <NOSSO_NUMERO>14000000099222698</NOSSO_NUMERO>
               <URL>https://200.201.168.67:8010/ecobranca/SIGCB/imprimir/0200656/14000000099222698</URL>
            </INCLUI_BOLETO>
         </DADOS>
      </manutencaocobrancabancaria:SERVICO_SAIDA>
   </soapenv:Body>
</soapenv:Envelope>
	`

	sDataErr := `
<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <manutencaocobrancabancaria:SERVICO_SAIDA xmlns:manutencaocobrancabancaria="http://caixa.gov.br/sibar/manutencao_cobranca_bancaria/boleto/externo" xmlns:sibar_base="http://caixa.gov.br/sibar">
            <sibar_base:HEADER>
                <VERSAO>1.0</VERSAO>
                <AUTENTICACAO>LvWr1op5Ayibn6jsCQ3/2bW4KwThVAlLK5ftxABlq201=</AUTENTICACAO>
                <USUARIO_SERVICO>SGCBS01D</USUARIO_SERVICO>
                <OPERACAO>INCLUI_BOLETO</OPERACAO>
                <SISTEMA_ORIGEM>SIGCB</SISTEMA_ORIGEM>
                <UNIDADE>1679</UNIDADE>
                <DATA_HORA>20170503102800</DATA_HORA>
            </sibar_base:HEADER>
            <COD_RETORNO>00</COD_RETORNO>
            <ORIGEM_RETORNO>MANUTENCAO_COBRANCA_BANCARIA</ORIGEM_RETORNO>
            <MSG_RETORNO></MSG_RETORNO>
            <DADOS>
                <CONTROLE_NEGOCIAL>
                    <ORIGEM_RETORNO>SIGCB</ORIGEM_RETORNO>
                    <COD_RETORNO>1</COD_RETORNO>
                    <MENSAGENS>
                        <RETORNO>(54) OPERACAO NAO PERMITIDA - HASH DIVERGENTE</RETORNO>
                    </MENSAGENS>
                </CONTROLE_NEGOCIAL>
            </DADOS>
        </manutencaocobrancabancaria:SERVICO_SAIDA>
    </soapenv:Body>
</soapenv:Envelope>
	`
	d, _ := ioutil.ReadAll(c.Request.Body)
	xml := string(d)
	//time.Sleep(400 * time.Millisecond)
	if strings.Contains(xml, "<VALOR>5.04</VALOR>") {
		c.AbortWithError(504, errors.New("Teste de Erro"))
	} else if strings.Contains(xml, "<VALOR>2.00</VALOR>") {
		c.Data(200, "text/xml", []byte(sData))
	} else {
		c.Data(200, "text/xml", []byte(sDataErr))
	}

}

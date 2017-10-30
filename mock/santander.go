package mock

import (
	"io/ioutil"
	"net/http/httputil"

	"github.com/gin-gonic/gin"

	"strings"

	"errors"
)

/*

0 Ticket validado ok
1 Erro, dados de entrada inválidos
2 Erro interno de criptografia
3 Erro, Ticket já utilizado anteriormente
4 Erro, Ticket gerado para outro sistema
5 Erro, Ticket expirado
6 Erro interno (dados)
7 Erro interno (timestamp)
*/
func getTicket(c *gin.Context) {
	const tok = `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"> 
	<soapenv:Body> 
	<dlwmin:createResponse xmlns:dlwmin="http://impl.webservice.dl.app.bsbr.altec.com/"> 
	<TicketResponse> 
	<retCode>0</retCode> 
	<ticket>OK</ticket> 
	</TicketResponse> 
	</dlwmin:createResponse> 
	</soapenv:Body> 
	</soapenv:Envelope>
	`

	const tokErr = `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"> 
	<soapenv:Body> 
	<dlwmin:createResponse xmlns:dlwmin="http://impl.webservice.dl.app.bsbr.altec.com/"> 
	<TicketResponse> 
	<retCode>3</retCode> 
	<ticket>NOK</ticket> 
	</TicketResponse> 
	</dlwmin:createResponse> 
	</soapenv:Body> 
	</soapenv:Envelope>
	`
	data, _ := httputil.DumpRequest(c.Request, true)
	str := string(data)
	if strings.Contains(str, "<value>3</value>") {
		c.Data(200, "text/json", []byte(tokErr))
	} else if strings.Contains(str, "<value>500</value>") {
		c.AbortWithError(500, errors.New("internal error"))
	} else {
		c.Data(200, "text/json", []byte(tok))
	}

}

func registerBoletoSantander(c *gin.Context) {
	const tokErr = `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Body>
      <dlwmin:registraTituloResponse xmlns:dlwmin="http://impl.webservice.ymb.app.bsbr.altec.com/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
         <return xmlns:ns2="http://impl.webservice.ymb.app.bsbr.altec.com/">
            <codcede/>
            <convenio>
               <codBanco/>
               <codConv/>
            </convenio>
            <descricaoErro>erro santander</descricaoErro>
            <dtNsu>05082017</dtNsu>
            <estacao>?</estacao>
            <nsu>?</nsu>
            <pagador>
               <bairro/>
               <cep/>
               <cidade/>
               <ender/>
               <nome/>
               <numDoc/>
               <tpDoc/>
               <uf/>
            </pagador>
            <situacao>20</situacao>
            <titulo>
               <aceito>N</aceito>
               <cdBarra>23192839182391829382132131238219893123</cdBarra>
               <dtEmissao/>
               <dtEntr/>
               <dtLimiDesc/>
               <dtVencto/>
               <especie/>
               <linDig>21321312382198931232132131238219893123</linDig>
               <mensagem/>
               <nossoNumero>313123131231231</nossoNumero>
               <pcJuro/>
               <pcMulta/>
               <qtDiasBaixa/>
               <qtDiasMulta/>
               <qtDiasProtesto/>
               <seuNumero/>
               <tpDesc/>
               <tpProtesto/>
               <vlAbatimento/>
               <vlDesc/>
               <vlNominal/>
            </titulo>
            <tpAmbiente>T</tpAmbiente>
         </return>
      </dlwmin:registraTituloResponse>
   </soapenv:Body>
</soapenv:Envelope>
	`

	const tok = `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Body>
      <dlwmin:registraTituloResponse xmlns:dlwmin="http://impl.webservice.ymb.app.bsbr.altec.com/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
         <return xmlns:ns2="http://impl.webservice.ymb.app.bsbr.altec.com/">
            <codcede/>
            <convenio>
               <codBanco/>
               <codConv/>
            </convenio>
            <descricaoErro>00000 - Título registrado em cobrança</descricaoErro>
            <dtNsu>05082017</dtNsu>
            <estacao>?</estacao>
            <nsu>?</nsu>
            <pagador>
               <bairro/>
               <cep/>
               <cidade/>
               <ender/>
               <nome/>
               <numDoc/>
               <tpDoc/>
               <uf/>
            </pagador>
            <situacao>00</situacao>
            <titulo>
               <aceito>N</aceito>
               <cdBarra>23192839182391829382132131238219893123</cdBarra>
               <dtEmissao/>
               <dtEntr/>
               <dtLimiDesc/>
               <dtVencto/>
               <especie/>
               <linDig>21321312382198931232132131238219893123</linDig>
               <mensagem/>
               <nossoNumero>313123131231231</nossoNumero>
               <pcJuro/>
               <pcMulta/>
               <qtDiasBaixa/>
               <qtDiasMulta/>
               <qtDiasProtesto/>
               <seuNumero/>
               <tpDesc/>
               <tpProtesto/>
               <vlAbatimento/>
               <vlDesc/>
               <vlNominal/>
            </titulo>
            <tpAmbiente>T</tpAmbiente>
         </return>
      </dlwmin:registraTituloResponse>
   </soapenv:Body>
</soapenv:Envelope>
	`

	body, _ := ioutil.ReadAll(c.Request.Body)
	str := string(body)
	if strings.Contains(str, "<ticket>OK</ticket>") {
		c.Data(200, "text/json", []byte(tok))
	} else {
		c.Data(200, "text/json", []byte(tokErr))
	}

}

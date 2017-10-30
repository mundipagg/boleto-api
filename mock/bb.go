package mock

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func authBB(c *gin.Context) {

	const tok = `{
					"access_token":"Yemfehvhs9PmpKLcJgSLzjfjSxHj4QREdKcrhHbM_oivhlYXmOIPn5j2Tp6FdZFmbIzCxVN-SdxnyoGT7cE5xg.AcQSHv1xiN0uaZ-efTafZIWPNenudJn9eU54TUAfiR0ff8RRE9thqgdRb2gm9t_uTREmOsBOz9jvQySnsPBqbfptOqz9-O_63c-LQq2ogxbu7iet-6te8V28gfjOVePnr87yIK8ueATW9ulb7jytYRYJd7CuZXF3PyD763tI2ykX-PNm2LAClqpAU-WAORQ_2OSLo5ElwPS_MgVAqvXm_n1PX0wPazW-YlwSvoYr9pYabiBAOCW4KkiZva0hRver7AMWlkP2t2M_wttG6wv1V8szty2Lb9oyGDL-cdfdit4rHFgpXp9dzG3qFS5qWxtO5tnQc3sVBBybDNINOmMlOaxKrsrNtQl5ncELy6jozyrPS-Yb3JhlvVaj3IDHq599bd30G8JMjDsGSc9wuws7Ws9tuUbOTiS-d2TaGOqTjKs.Mvco1yM5ErKRswd11TBjFYr-zX0kor4y1EDyq52M_ew874ifju5PaU-G5btc2zMhUcltTIVY1sJlTf2rqXKaPQ",
					"token_type":"Bearer",
					"refresh_token":"eyJhbGciOiJBMTI4S1ciLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0.4PdGzdqXluMJ67StpSmi5Ds5rWUXiLKvhYvZh_HR8DAjBt361RaGdw.Gcg48k3Omleobjs-c5J_mw.CtxZAiHOf_oA3c4uPKzgGesG6V-Y9QzFhJh8ww262jI-GQL2S6YqWe1ucrJ9oY_hrST05Y1ns7rTZJkGluDBscNtE3mIuv-WkCykHUDlor2gevZjxUApj98mJIKeFqfaeIGnXZpyeQBpPXAcCIELIjUN4CAWm99ed72DCcCWiPbO3v2smSQVLX04ESKqTbnjRyHQLiHGm8jP4PnOFIafdBrnRSfhsqIggJCZYNfIC1aRIrDnTSDiTBdx1vEruLOCFIOv9z4pqySPbImzC3Uxv9UxNDKvEa11TGoVYlnAx62_8d7pFAC8IeDwXNuaRzFklyDWZCMNtFl0pEB1bqh3mN6QdeQE2sfsoMhyif9iXqcFnUJvFAu4Oj981M_Vyh2GW7VTAvs67sw27xvCS1diJZGNLR_O09WEjn529MZGyT_4oWqmlVTb-a6dflFWwdI3DhsusgvT6pK_ja-eIXq5pw.o50PzlpZnNnks17cNsaKog","expires_in":1200
				}`
	time.Sleep(500 * time.Millisecond)
	c.Data(200, "text/json", []byte(tok))
}

func registerBoletoBB(c *gin.Context) {

	sData := `
		<?xml version="1.0" encoding="UTF-8"?>
		<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
		<SOAP-ENV:Body>
			<ns0:resposta xmlns:ns0="http://www.tibco.com/schemas/bws_registro_cbr/Recursos/XSD/Schema.xsd">
				<ns0:siglaSistemaMensagem />
				<ns0:codigoRetornoPrograma>0</ns0:codigoRetornoPrograma>
				<ns0:nomeProgramaErro />
				<ns0:textoMensagemErro />
				<ns0:numeroPosicaoErroPrograma>0</ns0:numeroPosicaoErroPrograma>
				<ns0:codigoTipoRetornoPrograma>0</ns0:codigoTipoRetornoPrograma>
				<ns0:textoNumeroTituloCobrancaBb>00010140510000066673</ns0:textoNumeroTituloCobrancaBb>
				<ns0:numeroCarteiraCobranca>17</ns0:numeroCarteiraCobranca>
				<ns0:numeroVariacaoCarteiraCobranca>19</ns0:numeroVariacaoCarteiraCobranca>
				<ns0:codigoPrefixoDependenciaBeneficiario>3851</ns0:codigoPrefixoDependenciaBeneficiario>
				<ns0:numeroContaCorrenteBeneficiario>8570</ns0:numeroContaCorrenteBeneficiario>
				<ns0:codigoCliente>932131545</ns0:codigoCliente>
				<ns0:linhaDigitavel>00190000090101405100500066673179971340000010000</ns0:linhaDigitavel>
				<ns0:codigoBarraNumerico>00199713400000100000000001014051000006667317</ns0:codigoBarraNumerico>
				<ns0:codigoTipoEnderecoBeneficiario>0</ns0:codigoTipoEnderecoBeneficiario>
				<ns0:nomeLogradouroBeneficiario>Cliente nao informado.</ns0:nomeLogradouroBeneficiario>
				<ns0:nomeBairroBeneficiario />
				<ns0:nomeMunicipioBeneficiario />
				<ns0:codigoMunicipioBeneficiario>0</ns0:codigoMunicipioBeneficiario>
				<ns0:siglaUfBeneficiario />
				<ns0:codigoCepBeneficiario>0</ns0:codigoCepBeneficiario>
				<ns0:indicadorComprovacaoBeneficiario />
				<ns0:numeroContratoCobranca>17414296</ns0:numeroContratoCobranca>
			</ns0:resposta>
		</SOAP-ENV:Body>
		</SOAP-ENV:Envelope>
	`

	sDataErr := `
		<?xml version="1.0" encoding="UTF-8"?>
		<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
		<SOAP-ENV:Body>
			<ns0:resposta xmlns:ns0="http://www.tibco.com/schemas/bws_registro_cbr/Recursos/XSD/Schema.xsd">
				<ns0:siglaSistemaMensagem />
				<ns0:codigoRetornoPrograma>ER500</ns0:codigoRetornoPrograma>
				<ns0:nomeProgramaErro>Nome Programa ERRO</ns0:nomeProgramaErro>
				<ns0:textoMensagemErro>Falha ao registrar Boleto</ns0:textoMensagemErro>
				<ns0:numeroPosicaoErroPrograma>0</ns0:numeroPosicaoErroPrograma>
				<ns0:codigoTipoRetornoPrograma>0</ns0:codigoTipoRetornoPrograma>
				<ns0:textoNumeroTituloCobrancaBb>00010140510000066673</ns0:textoNumeroTituloCobrancaBb>
				<ns0:numeroCarteiraCobranca>17</ns0:numeroCarteiraCobranca>
				<ns0:numeroVariacaoCarteiraCobranca>19</ns0:numeroVariacaoCarteiraCobranca>
				<ns0:codigoPrefixoDependenciaBeneficiario>3851</ns0:codigoPrefixoDependenciaBeneficiario>
				<ns0:numeroContaCorrenteBeneficiario>8570</ns0:numeroContaCorrenteBeneficiario>
				<ns0:codigoCliente>932131545</ns0:codigoCliente>
				<ns0:linhaDigitavel>00190000090101405100500066673179971340000010000</ns0:linhaDigitavel>
				<ns0:codigoBarraNumerico>00199713400000100000000001014051000006667317</ns0:codigoBarraNumerico>
				<ns0:codigoTipoEnderecoBeneficiario>0</ns0:codigoTipoEnderecoBeneficiario>
				<ns0:nomeLogradouroBeneficiario>Cliente nao informado.</ns0:nomeLogradouroBeneficiario>
				<ns0:nomeBairroBeneficiario />
				<ns0:nomeMunicipioBeneficiario />
				<ns0:codigoMunicipioBeneficiario>0</ns0:codigoMunicipioBeneficiario>
				<ns0:siglaUfBeneficiario />
				<ns0:codigoCepBeneficiario>0</ns0:codigoCepBeneficiario>
				<ns0:indicadorComprovacaoBeneficiario />
				<ns0:numeroContratoCobranca>17414296</ns0:numeroContratoCobranca>
			</ns0:resposta>
		</SOAP-ENV:Body>
		</SOAP-ENV:Envelope>
	`

	d, _ := ioutil.ReadAll(c.Request.Body)
	xml := string(d)
	//time.Sleep(700 * time.Millisecond)
	if strings.Contains(xml, "<sch:valorOriginalTitulo>2.00</sch:valorOriginalTitulo>") {
		c.Data(200, "text/xml", []byte(sData))
	} else {
		c.Data(200, "text/xml", []byte(sDataErr))
	}
}

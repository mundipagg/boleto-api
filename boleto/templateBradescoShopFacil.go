package boleto

const templateBoletoBradescoShopFacil = `
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <style>

        @media print
        {    
            .no-print, .no-print *
            {
                display: none !important;
            }
        }

        body {
            font-family: "Arial";
    		background-color: #fff;
            font-size:0.7em;
        }
        .left {
    		margin: auto;		
    		width: 216mm;
    	}
        .document {
            margin: auto auto;
            width: 216mm;            
        }

        .headerBtn {
            margin: auto auto;
            width: 216mm;
            background-color: #fff;
            display: none;
        }

        table {
            width: 100%;
            position: relative;
            border-collapse: collapse;
        }

        .boletoNumber {
            width: 66%;
            font-weight: bold;
            font-size:0.9em;
        }

        .center {
            text-align: center;
        }

        .right {
            text-align: right;
            right: 20px;
        }

        td {
            position: relative;
        }

        .title {
            position: absolute;
            left: 0px;
            top: 0px;
            font-size:0.65em;
            font-weight: bold;
        }

        .text {
             font-size:0.7em;
        }

        p.content {
            padding: 0px;
            width: 100%;
            margin: 0px;
            font-size:0.7em;
        }

        .sideBorders {
            border-left: 1px solid black;
            border-right: 1px solid black;
        }

        hr {
            size: 1;
            border: 1px dashed;
    		width: 216mm;
    		margin-top: 9mm;
        	margin-bottom: 9mm;
        }

        br {
            content: " ";
            display: block;
            margin: 12px 0;
            line-height: 12px;
        }

        .print {
            /* TODO(dbeam): reconcile this with overlay.css' .default-button. */
            background-color: rgb(77, 144, 254);
            background-image: linear-gradient(to bottom, rgb(77, 144, 254), rgb(71, 135, 237));
            border: 1px solid rgb(48, 121, 237);
            color: #fff;
            text-shadow: 0 1px rgba(0, 0, 0, 0.1);
        }

        .btnDefault {
            font-kerning: none;
            font-weight: bold;
        }

        .btnDefault:not(:focus):not(:disabled) {
            border-color: #808080;
        }

        button {
            border: 1px;
            padding: 5px;
            line-height: 20px;
        }

        span.iconFont {
            font-size: 20px;
        }

        span.align {
            display: inline-block;
            vertical-align: middle;
        }

        label {
            -moz-user-select: -moz-none;
            -khtml-user-select: none;
            -webkit-user-select: none;            
            -ms-user-select: none;
            user-select: none;
        }

        i[class*=icss-]{position:relative;display:inline-block;font-style:normal;background-color:currentColor;-webkit-box-sizing:border-box;box-sizing:border-box;vertical-align:middle}i[class*=icss-]:after,i[class*=icss-]:before{content:"";border-width:0;position:absolute;-webkit-box-sizing:border-box;box-sizing:border-box}i.icss-print{width:.68em;height:1em;border-style:solid;border-color:currentcolor;border-width:.07em;-webkit-border-radius:.05em;border-radius:.05em;background-color:transparent;margin:0 .17em}i.icss-print:before{width:1em;height:.4em;border-width:.07em .21em 0;border-style:solid;border-color:currentColor currentcolor transparent;-webkit-border-radius:.05em .05em 0 0;border-radius:.05em .05em 0 0;top:.25em;left:50%;-webkit-transform:translateX(-50%);-ms-transform:translateX(-50%);transform:translateX(-50%);background-image:-webkit-gradient(linear,left top,left bottom,color-stop(20%,transparent),color-stop(20%,currentcolor),color-stop(60%,currentcolor),color-stop(60%,transparent));background-image:-webkit-linear-gradient(transparent 20%,currentcolor 20%,currentcolor 60%,transparent 60%);background-image:-o-linear-gradient(transparent 20%,currentcolor 20%,currentcolor 60%,transparent 60%);background-image:linear-gradient(transparent 20%,currentcolor 20%,currentcolor 60%,transparent 60%)}i.icss-print:after{width:.45em;height:.065em;background-color:currentColor;left:50%;-webkit-transform:translateX(-50%);-ms-transform:translateX(-50%);transform:translateX(-50%);top:.6em;-webkit-box-shadow:0 .12em,-.1em -.28em 0 .05em;box-shadow:0 .12em,-.1em -.28em 0 .05em}i.icss-files{width:.75em;height:.95em;background-color:transparent;border:.05em solid transparent;border-width:0 .05em .05em 0;-webkit-box-shadow:inset 0 0 0 .065em,.13em .11em 0 -.05em;box-shadow:inset 0 0 0 .065em,.13em .11em 0 -.05em;-webkit-border-radius:0 .3em 0 0;border-radius:0 .3em 0 0;margin:0 .17em .05em .1em}i.icss-files:before{border-style:solid;border-width:.2em;top:.037em;left:.25em;-webkit-border-radius:.1em;border-radius:.1em;border-color:transparent currentColor transparent transparent;-webkit-transform:rotate(-45deg);-ms-transform:rotate(-45deg);transform:rotate(-45deg)}
        </style>

        <script type="text/javascript">
            window.onload = function getUrlParams() {
                var url_string = window.location.href;
                var url = new URL(url_string);
            
                var fmt = url.searchParams.get("fmt");
            
                if(fmt === "html") {
                    document.getElementById("headerBtn").style.display = "block";
                } 
            }
        </script>    
</head>

<body>

	<br/>
    <div class="headerBtn" id="headerBtn">
        <div style="text-align:right;">
            <button class="no-print btnDefault print" onclick="window.print()">
                <i class="icss-print"></i>
                <span class="align">&nbspImprimir</span>
            </button>
            <button class="no-print btnDefault print" onclick="window.location='./boleto?fmt=pdf&id={{bsonMongoToString .View.ID}}&pk={{.View.PublicKey}}'">
                <i class="icss-files"></i>
                <span class="align">&nbspGerar PDF</span>
            </button>            
        </div>
    </div>
    <br/>
    {{template "boletoForm" .}}
	<hr/>
    {{template "boletoForm" .}}
	<div class="left">
		<img style="margin-left:5mm;" id="barcode_{{printIfNotProduction .View.Barcode}}" src="data:image/png;base64,{{.Barcode64}}" alt="">
		<br/>
		</div>
    </div>
 </body>
</html>
`

const boletoFormBradescoShopFacil = `
{{define "boletoForm"}}
<div class="document">
        <table cellspacing="0" cellpadding="0">
            <tr class="topLine">
                <td class="bankLogo">
                    {{.ConfigBank.Logo}}					
                </td>
                <td class="sideBorders center"><span style="font-weight:bold;font-size:0.9em;">{{.View.BankNumber}}</span></td>
                <td class="boletoNumber center"><img src="data:image/png;base64,{{.DigitableLine}}" line="{{printIfNotProduction .View.DigitableLine}}"  /></td>
            </tr>
        </table>
        <table cellspacing="0" cellpadding="0" border="1">
            <tr>
                <td width="80%" colspan="6">
                    <span class="title">Local de Pagamento</span>
                    <br/>
                    <span class="text">ATÉ O VENCIMENTO EM QUALQUER BANCO OU CORRESPONDENTE NÃO BANCÁRIO</span>
                </td>
                <td width="20%">
                    <span class="title">Data de Vencimento</span>
                    <br/>
                    <br/>
                    <p class="content right text" style="font-weight:bold;" id="expire_date">{{.View.Boleto.Title.ExpireDateTime | brdate}}</p>
                </td>
            </tr>
            <tr>
                <td width="80%" colspan="6">
                    <span class="title">Nome do Beneficiário / CNPJ / CPF / Endereço:</span>
                    <br/>
                    <table border="0" style="border:none">
                        <tr>
                            <td width="60%"><span class="text" id="recipient_name">{{.View.Boleto.Recipient.Name}}</span></td>
                            <td><span class="text" id="recipient_document"><b>{{.View.Boleto.Recipient.Document.Type}}</b> {{fmtDoc .View.Boleto.Recipient.Document}}</span></td>
                        </tr>
                    </table>
                    <br/>
                    <span class="text" id="recipient_address">{{.View.Boleto.Recipient.Address.Street}}, 
                    {{.View.Boleto.Recipient.Address.Number}} - 
                    {{.View.Boleto.Recipient.Address.District}}, 
                    {{.View.Boleto.Recipient.Address.StateCode}} - 
                    {{.View.Boleto.Recipient.Address.ZipCode}}</span>
                </td>
                <td width="20%">
                    <span class="title">Agência/Código Beneficiário</span>
                    <br/>
                    <br/>
                    <p class="content right" id="agreement_agency_account">


                    {{padLeft .View.Boleto.Agreement.Agency "0" 4}}{{if eq .View.Boleto.Agreement.AgencyDigit ""}}{{else}}-{{.View.Boleto.Agreement.AgencyDigit}}{{end}}
                    /
                    {{padLeft .View.Boleto.Agreement.Account "0" 7}}{{if eq .View.Boleto.Agreement.AccountDigit ""}}{{else}}-{{.View.Boleto.Agreement.AccountDigit}}{{end}}
                    </p>
                </td>
            </tr>

            <tr>
                <td width="20%">
                    <span class="title">Data do Documento</span>
                    <br/>
                    <p class="content center" id="create_date">{{.View.Boleto.Title.CreateDate | brdate}}</p>
                </td>
                <td width="17%" colspan="2">
                    <span class="title">Num. do Documento</span>
                    <br/>
                    <p class="content center" id="boleto_document_number">{{.View.Boleto.Title.DocumentNumber}}</p>
                </td>
                <td width="10%">
                    <span class="title">Espécie doc</span>
                    <br/>
                    <p class="content center" id="configbank_especie_doc">{{.ConfigBank.EspecieDoc}}</p>
                </td>
                <td width="8%">
                    <span class="title">Aceite</span>
                    <br/>
                    <p class="content center" id="configbank_aceite" >{{.ConfigBank.Aceite}}</p>
                </td>
                <td>
                    <span class="title">Data Processamento</span>
                    <br/>
                    <p class="content center" id="process_date">{{.View.Boleto.Title.CreateDate | brdate}}</p>
                </td>
                <td width="30%">
                    <span class="title">Carteira/Nosso Número</span>
                    <br/>
                    <br/>
                    <p class="content right" id="ournumber">
                        {{padLeft (toString16 .View.Boleto.Agreement.Wallet) "0" 2}}
                        / 
                        {{padLeft (toString .View.Boleto.Title.OurNumber) "0" 11}}-{{mod11BradescoShopFacilDv (padLeft (toString .View.Boleto.Title.OurNumber) "0" 11) (padLeft (toString16 .View.Boleto.Agreement.Wallet) "0" 2)}}
                    </p>                   
                </td>
            </tr>

            <tr>
                {{if eq .View.BankNumber "033-7"}}
                <td width="29%" colspan="2">
                    <table>
                        <tr>                            
                            <td>
                                <span class="title">Carteira</span>
                                <br/>
                                <p class="content center" id="wallet">COBRANCA SIMPLES RCR</p>
                            </td>
                        </tr>
                    </table>                
                </td>
                {{else}}
                <td width="20%">
                    <span class="title">Uso do Banco</span>
                    <br/>
                    <p class="content center">&nbsp;</p>
                </td>                
                <td width="14%">
                    <table>
                        <tr>
                            {{if eq .View.BankNumber "237-2"}}
                                <td style="border-right: 1px solid #808080;" id="cel_cip">
                                    <span class="title">Cip</span>
                                    <br/>
                                    <p class="content center" id="wallet">865</p>
                                </td>
                            {{end}}

                            <td>
                                <span class="title">Carteira</span>
                                <br/>
                                <p class="content center" id="wallet">
                                    {{padLeft (toString16 .View.Boleto.Agreement.Wallet) "0" 2}}
                                </p>
                            </td>
                        </tr>
                    </table>
                    
                </td>
                {{end}}
                <td width="10%">
                    <span class="title">Espécie</span>
                    <br/>
                    <p class="content center">{{.ConfigBank.Moeda}}</p>
                </td>
                <td width="8%" colspan="2">
                    <span class="title">Quantidade</span>
                    <br/>
                    <p class="content center" id="configbank_quantidade">{{.ConfigBank.Quantidade}}</p>
                </td>
                <td>
                    <span class="title">Valor</span>
                    <br/>
                    <p class="content center" id="configbank_valorCotacao" >{{.ConfigBank.ValorCotacao}}</p>
                </td>
                <td width="30%">
                    <span class="title">(=) Valor do Documento</span>
                    <br/>
                    <br/>
                    <p class="content right" id="amount_in_cents" >{{fmtNumber .View.Boleto.Title.AmountInCents}}</p>
                </td>
            </tr>
            <tr>
                <td colspan="6" rowspan="4">
                    <span class="title">Instruções de responsabilidade do BENEFICIÁRIO. Qualquer dúvida sobre este boleto contate o beneficiário.</span>
                    <p class="content" id="instructions">{{.View.Boleto.Title.Instructions }}</p>
                </td>
            </tr>
            <tr>
                <td>
                    <span class="title">(-) Descontos/Abatimento</span>
                    <br/>
                    <p class="content right">&nbsp;</p>
                </td>
            </tr>
            <tr>
                <td>
                    <span class="title">(+) Juros/Multa</span>
                    <br/>
                    <p class="content right">&nbsp;</p>
                </td>
            </tr>
            <tr>
                <td>
                    <span class="title">(=) Valor Pago</span>
                    <br/>
                    <p class="content right">&nbsp;</p>
                </td>
            </tr>
            <tr>
                <td colspan="7">
                    <table border="0" style="border:none">
                        <tr>
                            <td width="60%"><span class="text" id="buyer_name"><b>Nome do Pagador: </b>&nbsp;{{.View.Boleto.Buyer.Name}}</span></td>
                            <td><span class="text" id="buyer_document"><b>CNPJ/CPF: </b>&nbsp;{{fmtDoc .View.Boleto.Buyer.Document}}</span></td>
                        </tr>
                        <tr>
                            <td><span class="text" id="buyer_address"><b>Endereço: </b>&nbsp;{{.View.Boleto.Buyer.Address.Street}}&nbsp;{{.View.Boleto.Buyer.Address.Number}}, {{.View.Boleto.Buyer.Address.District}} - {{.View.Boleto.Buyer.Address.City}}, {{.View.Boleto.Buyer.Address.StateCode}} - {{.View.Boleto.Buyer.Address.ZipCode}}</span></td>
                            <td>&nbsp;</td>
                        </tr>
                        <tr>
                            <td><span class="text"><b>Sacador/Avalista: </b> &nbsp;</span></td>
                            <td><span class="text"><b>CNPJ/CPF: </b> &nbsp;</span></td>
                        </tr>
                    </table>
                </td>
            </tr>            
        </table>
		<br/>
    </div>

	{{end}}
`

const boletoPropostaLayoutFormBradescoShopFacil = `
{{define "boletoForm"}}
<div class="document">
    <table cellspacing="0" cellpadding="0">
        <tbody>
            <tr class="topLine">
                <td class="bankLogo">{{.ConfigBank.Logo}}</td>
                <td class="sideBorders center"><span style="font-weight:700;font-size:.9em">{{.View.BankNumber}}</span></td>
                <td class="boletoNumber center"><img src="data:image/png;base64,{{.DigitableLine}}" line="{{printIfNotProduction .View.DigitableLine}}"  /></td>
            </tr>
        </tbody>
    </table>
    <table cellspacing="0" cellpadding="0" border="1">
        <tbody>
            <tr>
                <td colspan="6">
                    <p class="center" style="font-weight:bold">BOLETO DE PROPOSTA</p>
                    <br>
                    <p class="text" style="font-weight:bold">ESTE BOLETO SE REFERE A UMA PROPOSTA JÁ FEITA A VOCÊ E O
                        SEU PAGAMENTO NÃO É OBRIGATÓRIO.
                        <br>Deixar de pagá-lo não dará causa a protesto, a cobrança judicial ou extrajudicial, nem a
                        inserção de seu nome em cadastro de restrição de credito.
                        <br>Pagar até a data do vencimento significa aceitar a proposta.
                        <br>Informações adicionais sobre a proposta e sobre o respectivo contrato poderão ser
                        solicitadas a qualquer momento ao beneficiário,por meio de seus canais de atendimento.</p>
                </td>
            </tr>
            <tr>
                <td colspan="6"><span class="title">Local de Pagamento</span><br><span class="text">ATÉ O VENCIMENTO EM
                        QUALQUER BANCO OU CORRESPONDENTE NÃO BANCÁRIO</span></td>
            </tr>
            <tr>
                <td colspan="5"><span class="title">Nome do Beneficiário / CNPJ / CPF / Endereço:</span>
                    <br>
                    <span class="text" id="recipient_name">{{.View.Boleto.Recipient.Name}}</span>
                    <span class="text" id="recipient_document" style="margin-left:5em"><b>{{.View.Boleto.Recipient.Document.Type}}</b> {{fmtDoc .View.Boleto.Recipient.Document}}</span>
                    <br>
                    <span class="text" id="recipient_address">{{.View.Boleto.Recipient.Address.Street}}, 
                    {{.View.Boleto.Recipient.Address.Number}} - 
                    {{.View.Boleto.Recipient.Address.District}}, 
                    {{.View.Boleto.Recipient.Address.StateCode}} - 
                    {{.View.Boleto.Recipient.Address.ZipCode}}</span>
                </td>
                <td width="20%"><span class="title">Data de Vencimento</span><br><br>
                    <p class="content right text" style="font-weight:700" id="expire_date">{{.View.Boleto.Title.ExpireDateTime | brdate}}</p>
                </td>
            </tr>
            <tr>
                <td width="20%"><span class="title">Data Processamento</span><br>
                    <p class="content center" id="process_date">{{.View.Boleto.Title.CreateDate | brdate}}</p>
                </td>
                <td width="17%"><span class="title">Num. do Documento</span><br>
                    <p class="content center" id="boleto_document_number">{{.View.Boleto.Title.DocumentNumber}}</p>
                </td>
                <td><span class="title">Nosso Número</span><br><br>
                    <p class="content center" id="ournumber">
                    {{padLeft (toString .View.Boleto.Title.OurNumber) "0" 11}}-{{mod11BradescoShopFacilDv (padLeft (toString .View.Boleto.Title.OurNumber) "0" 11) (padLeft (toString16 .View.Boleto.Agreement.Wallet) "0" 2)}}
                    </p>
                </td>
                <td width="22%"><span class="title">Agência/Código Beneficiário</span><br>
                    <p class="content center" id="agreement_agency_account">
                    {{padLeft .View.Boleto.Agreement.Agency "0" 4}}{{if eq .View.Boleto.Agreement.AgencyDigit ""}}{{else}}-{{.View.Boleto.Agreement.AgencyDigit}}{{end}}
                    /
                    {{padLeft .View.Boleto.Agreement.Account "0" 7}}{{if eq .View.Boleto.Agreement.AccountDigit ""}}{{else}}-{{.View.Boleto.Agreement.AccountDigit}}{{end}}
                    </p>
                </td>
                <td width="8%"><span class="title">Carteira</span><br>
                    <p class="content center" id="wallet">{{.View.Boleto.Agreement.Wallet}}</p>
                </td>
                <td><span class="title">(=) Valor do Documento</span><br><br>
                    <p class="content right" id="amount_in_cents">{{fmtNumber .View.Boleto.Title.AmountInCents}}</p>
                </td>
            </tr>
            <tr>
                <td colspan="5" rowspan="2"><span class="title">Informações de responsabilidade do
                        beneficiário</span><br><br><p class="content" id="instructions">{{.View.Boleto.Title.Instructions }}</p></td>
                <td width="20%"><span class="title">(-) Descontos/Abatimento</span><br>
                    <p class="content right">&nbsp;</p>
                </td>
            </tr>
            <tr>
                <td><span class="title">(=) Valor Pago</span><br>
                    <p class="content right">&nbsp;</p>
                </td>
            </tr>
            <tr>
                <td colspan="7">
                    <table border="0" style="border:none">
                        <tbody>
                            <tr>
                                <td width="60%"><span class="text" id="buyer_name"><b>Nome do Pagador: </b>&nbsp;{{.View.Boleto.Buyer.Name}}</span></td>
                                <td><span class="text" id="buyer_document"><b>CNPJ/CPF: </b>&nbsp;{{fmtDoc .View.Boleto.Buyer.Document}}</span>
                                </td>
                            </tr>
                            <tr>
                                <td><span class="text" id="buyer_address"><b>Endereço: </b>&nbsp;{{.View.Boleto.Buyer.Address.Street}}&nbsp;{{.View.Boleto.Buyer.Address.Number}}, {{.View.Boleto.Buyer.Address.District}} - {{.View.Boleto.Buyer.Address.City}}, {{.View.Boleto.Buyer.Address.StateCode}} - {{.View.Boleto.Buyer.Address.ZipCode}}</span></td>
                                <td>&nbsp;</td>
                            </tr>
                            <tr>
                                <td><span class="text"><b>Sacador/Avalista:</b> &nbsp;</span></td>
                                <td><span class="text"><b>CNPJ/CPF:</b> &nbsp;</span></td>
                            </tr>
                        </tbody>
                    </table>
                </td>
            </tr>
        </tbody>

    </table><br>
</div>
	{{end}}
`

func getTemplateBradescoShopFacil(bt string) (string, string) {
	if bt == "BDP" {
		return templateBoletoBradescoShopFacil, boletoPropostaLayoutFormBradescoShopFacil
	}
	return templateBoletoBradescoShopFacil, boletoFormBradescoShopFacil
}

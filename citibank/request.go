package citibank

const responseCiti = `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Header/>
   <soapenv:Body>
      <RegisterBoletoResponse>
         <actionCode>{{returnCode}}</actionCode>
         <reasonMessage>{{returnMessage}}</reasonMessage>
      </RegisterBoletoResponse>
   </soapenv:Body>
</soapenv:Envelope>
`

const registerBoletoCiti = `

## SOAPAction:RegisterBoleto
## Content-Type:text/xml

<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Header/>
   <soapenv:Body>
      <GrpREMColTit>
         <GrpBenf>
            <CdClrSys>745</CdClrSys>
            <CdIspb>33479023</CdIspb>
            <CdtrId>{{.Agreement.AgreementNumber}}</CdtrId>
            <CdtrNm>{{.Recipient.Name}}</CdtrNm>
            <CdtrTaxId>{{.Recipient.Document.Number}}</CdtrTaxId>
            <CdtrTaxTp>J</CdtrTaxTp>
         </GrpBenf>
         <GrpClPgd>
            <DbtrNm>{{.Buyer.Name}}</DbtrNm>
            {{if (eq .Buyer.Document.Type "CPF")}}				
            	<DbtrTaxTp>F</DbtrTaxTp>
			{{else}}
            	<DbtrTaxTp>J</DbtrTaxTp>
			{{end}}
			<DbtrTaxId>{{.Buyer.Document.Number}}</DbtrTaxId>
            <GrpClPgdAdr>
               <DbtrAdrTp>{{.Buyer.Address.Street}} {{.Buyer.Address.Number}} {{.Buyer.Address.Complement}}</DbtrAdrTp>
               <DbtrCtrySubDvsn>{{.Buyer.Address.StateCode}}</DbtrCtrySubDvsn>
               <DbtrPstCd>{{.Buyer.Address.ZipCode}}</DbtrPstCd>
               <DbtrTwnNm>{{.Buyer.Address.City}}</DbtrTwnNm>
            </GrpClPgdAdr>
         </GrpClPgd>
         <CdOccTp>01</CdOccTp>
         <DbtrGrntNm />
         <DbtrMsg>{{.Title.Instructions}}</DbtrMsg>
         <TitlAmt>{{.Title.AmountInCents}}</TitlAmt>
         <TitlBarCdInd>0</TitlBarCdInd>
         <TitlCcyCd>09</TitlCcyCd>
         <TitlCiaCdId>{{.Title.DocumentNumber}}</TitlCiaCdId>
         <TitlDueDt>{{enDate .Title.ExpireDateTime "-"}}</TitlDueDt>
         <TitlInstrNmDtExec>0</TitlInstrNmDtExec>
         <TitlInstrProtInd> </TitlInstrProtInd>
         <TitlInstrWrtOffInd> </TitlInstrWrtOffInd>
         <TitlIOFAmt>0</TitlIOFAmt>
         <TitlIssDt>{{enDate today "-"}}</TitlIssDt>
         <TitlOurNb>{{padLeft (toString .Title.OurNumber) "0" 12}}</TitlOurNb>
         <TitlPortCd>1</TitlPortCd>
         <TitlRbtAmt>0</TitlRbtAmt>
         <TitlTpCd>03</TitlTpCd>
         <TitlYourNb>{{.Title.DocumentNumber}}</TitlYourNb>
         <GrpDscnt>
            <TitlDscntAmtOrPrct>0</TitlDscntAmtOrPrct>
            <TitlDscntEndDt> </TitlDscntEndDt>
            <TitlDscntTp> </TitlDscntTp>
         </GrpDscnt>
         <GrpItrs>
            <TitlItrsAmtOrPrct>0</TitlItrsAmtOrPrct>
            <TitlItrsStrDt> </TitlItrsStrDt>
            <TitlItrsTp> </TitlItrsTp>
         </GrpItrs>
         <GrpFn>
            <TitlFnAmtOrPrct>0</TitlFnAmtOrPrct>
            <TitlFnStrDt> </TitlFnStrDt>
            <TitlFnTp> </TitlFnTp>
         </GrpFn>
      </GrpREMColTit>
   </soapenv:Body>
</soapenv:Envelope>
`

//GetRegisterBoletoCitiTmpl retorna o padrão de registro de boleto do Citibank
func GetRegisterBoletoCitiTmpl() string {
	return registerBoletoCiti
}

//GetResponseTemplateCiti retorna o template de mensagem do Citibank
func GetResponseTemplateCiti() string {
	return responseCiti
}

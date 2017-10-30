package mock

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerBoletoCiti(c *gin.Context) {
	sData := `
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:s1="http://www.citibank.com.br/comercioeletronico/registerboleto">
    <soap:Body>
        <s1:RegisterBoletoResponse>
            <actionCode>0</actionCode>
            <reasonMessage>Data received                           </reasonMessage>
            <TitlBarCd>74591728800000001033100087772012000000421265</TitlBarCd>
            <TitlDgtLine>74593100048777201200800004212650172880000000103 </TitlDgtLine>
        </s1:RegisterBoletoResponse>
    </soap:Body>
</soap:Envelope>
	`

	sDataErr := `
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:s1="http://www.citibank.com.br/comercioeletronico/registerboleto">
    <soap:Body>
        <s1:RegisterBoletoResponse>
            <actionCode>99</actionCode>
            <reasonMessage>Data not processed                      </reasonMessage>
            <TitlBarCd></TitlBarCd>
            <TitlDgtLine></TitlDgtLine>
        </s1:RegisterBoletoResponse>
    </soap:Body>
</soap:Envelope>
`
	d, _ := ioutil.ReadAll(c.Request.Body)
	xml := string(d)
	if strings.Contains(xml, "<TitlAmt>200</TitlAmt>") {
		c.Data(200, "text/xml", []byte(sData))
	} else {
		c.Data(200, "text/xml", []byte(sDataErr))
	}

}

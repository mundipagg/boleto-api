package mock

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerStone(c *gin.Context) {
	b, _ := ioutil.ReadAll(c.Request.Body)
	const resp = `{
    "id": "2c32b477-2da1-4754-8252-a85398e2ef71",
    "account_id": "6fba8647-4821-4aa5-8e70-3b7387cbeff2",
    "created_by": "user:17f75b93-80c2-4c18-9fad-a2fe74d759da",
    "created_at": "2020-01-02T22:01:09Z",
    "registered_at": nil,
    "settled_at": nil,
    "amount": 5000,
    "barcode": "19797770000000020007115000002186110706666952",
    "writable_line": "19797115040000218611207066669529777000000002000",
    "expiration_date": "2020-05-30",
    "invoice_type": "proposal",
    "issuance_date": "2019-01-01",
    "limit_date": "2020-01-02",
    "status": "CREATED",
    "our_number": "20200427112605941833",
    "beneficiary": {
      "account_code":"498910",
      "branch_code":"1",
      "document": "12345678912451",
      "document_type": "cnpj",
      "legal_name": "Valim da Serra",
      "trade_name": "Empresa Legal"
    },
    "payer": {
      "document": "12345678912451",
      "document_type": "cnpj",
      "legal_name": "Rafaela Almeida",
      "trade_name": "ABC"
    }
  }`
	if strings.Contains(string(b), `"amount": "5000"`) {
		c.Data(200, "text/json", []byte(resp))
	} else {
		c.Data(400, "text/json", []byte(`
			{}
		`))
	}

}

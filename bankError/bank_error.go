package bankError

import "github.com/mundipagg/boleto-api/models"

func ParseError(bankError models.ErrorResponse, bankName string) error {

	mapping := getBankMap(bankName)

	if mapping[bankError.Code] != "" {
		err := models.BadGatewayError{
			Message: bankError.Message,
		}
		return err
	}

	return nil
}

func getBankMap(bankName string) map[string]string {

	switch bankName {
	case "BradescoNetEmpresa":
		return bradescoErrorMap
	default:
		return make(map[string]string)
	}
}

var bradescoErrorMap = map[string]string{
	"800":     "Erro de certificado ou formatação de campos",
	"810":     "Erro de certificado ou formatação de campos",
	"1290030": "Fim anormal do programa",
	"-99":     "Serviço indisponível",
}

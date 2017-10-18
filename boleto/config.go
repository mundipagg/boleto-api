package boleto

import (
	"html/template"

	"github.com/mundipagg/boleto-api/models"
)

//ConfigBank configure html template for each bank
type ConfigBank struct {
	Logo         template.HTML
	EspecieDoc   string
	Aceite       string
	Quantidade   string
	ValorCotacao string
	Moeda        string
}

//GetConfig returns boleto configution for each bank
func GetConfig(number models.BankNumber) ConfigBank {
	switch number {
	case models.BancoDoBrasil:
		return configBB()
	case models.Santander:
		return configSantander()
	case models.Citibank:
		return configCiti()
	case models.Bradesco:
		return configBradesco()
	case models.Caixa:
		return configCaixa()
	default:
		return configBB()
	}
}

func configCiti() ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoCiti), EspecieDoc: "DMI", Aceite: "N", Quantidade: "", ValorCotacao: "", Moeda: "R$"}
}

func configBB() ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoBB), EspecieDoc: "DM", Aceite: "N", Quantidade: "N", ValorCotacao: "", Moeda: "R$"}
}

func configCaixa() ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoCaixa), EspecieDoc: "DM", Aceite: "N", Quantidade: "N", ValorCotacao: "", Moeda: "R$"}
}

func configSantander() ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoSantander), EspecieDoc: "DM", Aceite: "N", Quantidade: "N", ValorCotacao: "", Moeda: "R$"}
}

func configBradesco() ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoBradesco), EspecieDoc: "Outro", Aceite: "N", Quantidade: "", ValorCotacao: "", Moeda: "Real"}
}

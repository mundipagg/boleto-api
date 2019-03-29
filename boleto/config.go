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
func GetConfig(boleto models.BoletoRequest) ConfigBank {
	switch boleto.BankNumber {
	case models.BancoDoBrasil:
		return configBB(boleto)
	case models.Santander:
		return configSantander(boleto)
	case models.Citibank:
		return configCiti(boleto)
	case models.Bradesco:
		return configBradesco(boleto)
	case models.Caixa:
		return configCaixa(boleto)
	case models.Itau:
		return configItau(boleto)
	case models.Pefisa:
		return configPefisa(boleto)
	default:
		return configBB(boleto)
	}
}

func configCiti(boleto models.BoletoRequest) ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoCiti), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "", ValorCotacao: "", Moeda: "R$"}
}

func configBB(boleto models.BoletoRequest) ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoBB), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "N", ValorCotacao: "", Moeda: "R$"}
}

func configCaixa(boleto models.BoletoRequest) ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoCaixa), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "", ValorCotacao: "", Moeda: "R$"}
}

func configSantander(boleto models.BoletoRequest) ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoSantander), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "N", ValorCotacao: "", Moeda: "R$"}
}

func configItau(boleto models.BoletoRequest) ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoItau), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "N", ValorCotacao: "", Moeda: "R$"}
}

func configPefisa(boleto models.BoletoRequest) ConfigBank {
	return ConfigBank{Logo: template.HTML(LogoPefisa), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "", ValorCotacao: "", Moeda: "R$"}
}

func configBradesco(boleto models.BoletoRequest) ConfigBank {
	switch boleto.Agreement.Wallet {
	case 4, 9, 19:
		return ConfigBank{Logo: template.HTML(LogoBradesco), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "", ValorCotacao: "", Moeda: "Real"}
	case 25, 26:
		return ConfigBank{Logo: template.HTML(LogoBradesco), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "", ValorCotacao: "", Moeda: "Real"}
	default:
		return ConfigBank{Logo: template.HTML(LogoBradesco), EspecieDoc: boleto.Title.BoletoType, Aceite: "N", Quantidade: "", ValorCotacao: "", Moeda: "Real"}
	}
}

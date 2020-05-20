package models

import (
	"time"

	"github.com/mundipagg/boleto-api/util"

	"gopkg.in/mgo.v2/bson"

	"github.com/mundipagg/boleto-api/config"

	"github.com/PMoneda/flow"
	"github.com/google/uuid"

	"fmt"

	"encoding/json"
	"strconv"
)

// BoletoRequest entidade de entrada para o boleto
type BoletoRequest struct {
	Authentication Authentication `json:"authentication"`
	Agreement      Agreement      `json:"agreement"`
	Title          Title          `json:"title"`
	Recipient      Recipient      `json:"recipient"`
	Buyer          Buyer          `json:"buyer"`
	BankNumber     BankNumber     `json:"bankNumber"`
	RequestKey     string         `json:"requestKey,omitempty"`
}

// BoletoResponse entidade de saída para o boleto
type BoletoResponse struct {
	StatusCode    int    `json:"-"`
	Errors        Errors `json:"errors,omitempty"`
	ID            string `json:"id,omitempty"`
	DigitableLine string `json:"digitableLine,omitempty"`
	BarCodeNumber string `json:"barCodeNumber,omitempty"`
	OurNumber     string `json:"ourNumber,omitempty"`
	Links         []Link `json:"links,omitempty"`
}

//Link é um tipo padrão no restfull para satisfazer o HATEOAS
type Link struct {
	Href   string `json:"href,omitempty"`
	Rel    string `json:"rel,omitempty"`
	Method string `json:"method,omitempty"`
}

// BoletoView contem as informações que serão preenchidas no boleto
type BoletoView struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	UID           string
	SecretKey     string
	PublicKey     string        `json:"pk,omitempty"`
	Format        string        `json:"format,omitempty"`
	Boleto        BoletoRequest `json:"boleto,omitempty"`
	BankID        BankNumber    `json:"bankId,omitempty"`
	CreateDate    time.Time     `json:"createDate,omitempty"`
	BankNumber    string        `json:"bankNumber,omitempty"`
	DigitableLine string        `json:"digitableLine,omitempty"`
	OurNumber     string        `json:"ourNumber,omitempty"`
	Barcode       string        `json:"barcode,omitempty"`
	Barcode64     string        `json:"barcode64,omitempty"`
	Links         []Link        `json:"links,omitempty"`
}

// NewBoletoView cria um novo objeto view de boleto a partir de um boleto request, codigo de barras e linha digitavel
func NewBoletoView(boleto BoletoRequest, response BoletoResponse, bankName string) BoletoView {
	boleto.Authentication = Authentication{}
	uid, _ := uuid.NewUUID()
	id := bson.NewObjectId()
	view := BoletoView{
		ID:            id,
		UID:           uid.String(),
		SecretKey:     uid.String(),
		BankID:        boleto.BankNumber,
		Boleto:        boleto,
		Barcode:       response.BarCodeNumber,
		DigitableLine: response.DigitableLine,
		OurNumber:     response.OurNumber,
		BankNumber:    boleto.BankNumber.GetBoletoBankNumberAndDigit(),
		CreateDate:    time.Now(),
	}
	view.GeneratePublicKey()
	view.Links = view.CreateLinks()
	if len(response.Links) > 0 && bankName == "BradescoShopFacil" {
		view.Links = append(view.Links, response.Links[0])
	}
	return view
}

//EncodeURL tranforma o boleto view na forma que será escrito na url
func (b *BoletoView) EncodeURL(format string) string {
	idBson, _ := b.ID.MarshalText()
	_url := fmt.Sprintf("%s?fmt=%s&id=%s&pk=%s", config.Get().AppURL, format, string(idBson), b.PublicKey)

	return _url
}

//CreateLinks cria a lista de links com os formatos suportados
func (b *BoletoView) CreateLinks() []Link {
	links := make([]Link, 0, 3)
	for _, f := range []string{"html", "pdf"} {
		links = append(links, Link{Href: b.EncodeURL(f), Rel: f, Method: "GET"})
	}
	return links
}

//ToJSON tranforma o boleto view em json
func (b BoletoView) ToJSON() string {
	json, _ := json.Marshal(b)
	return string(json)
}

//GeneratePublicKey Gera a chave pública criptografada para geração da URL do boleto
func (b *BoletoView) GeneratePublicKey() {
	s := b.SecretKey + b.CreateDate.String() + b.Barcode + b.Boleto.Buyer.Document.Number + strconv.FormatUint(b.Boleto.Title.AmountInCents, 10)
	b.PublicKey = util.Sha256(s, "hex")
}

// BankNumber número de identificação do banco
type BankNumber int

// IsBankNumberValid verifica se o banco enviado existe
func (b BankNumber) IsBankNumberValid() bool {
	switch b {
	case BancoDoBrasil, Itau, Santander, Caixa, Bradesco, Citibank, Stone:
		return true
	default:
		return false
	}
}

//GetBoletoBankNumberAndDigit Retorna o numero da conta do banco do boleto
func (b BankNumber) GetBoletoBankNumberAndDigit() string {
	switch b {
	case BancoDoBrasil:
		return "001-9"
	case Caixa:
		return "104-0"
	case Santander:
		return "033-7"
	case Itau:
		return "341-7"
	case Bradesco:
		return "237-2"
	case Pefisa:
		return "174"
	case Stone:
		return "197"
	default:
		return ""
	}
}

const (
	// BancoDoBrasil constante do Banco do Brasil
	BancoDoBrasil = 1

	// Santander constante do Santander
	Santander = 33

	// Itau constante do Itau
	Itau = 341

	//Bradesco constante do Bradesco
	Bradesco = 237
	// Caixa constante do Caixa
	Caixa = 104

	// Citibank constante do Citi
	Citibank = 745

	//Real constante do REal
	Real = 9

	// Pefisa constante do Pefisa
	Pefisa = 174

	// Stone
	Stone = 197
)

// BoletoErrorConector é um connector flow para criar um objeto de erro
func BoletoErrorConector(e *flow.ExchangeMessage, u flow.URI, params ...interface{}) error {
	b := "Erro interno"
	switch t := e.GetBody().(type) {
	case error:
		b = t.Error()
	case string:
		b = t
	}

	st, err := strconv.Atoi(e.GetHeader("status"))
	if err != nil {
		st = 0
	}
	resp := BoletoResponse{}
	resp.Errors = make(Errors, 0, 0)
	resp.Errors.Append("MP"+e.GetHeader("status"), b)
	resp.StatusCode = st
	e.SetBody(resp)
	return nil
}

//HasErrors verify if Response has any error
func (b *BoletoResponse) HasErrors() bool {
	return b.Errors != nil && len(b.Errors) > 0
}

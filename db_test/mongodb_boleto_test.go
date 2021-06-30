package db_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/caixa"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetBoletoById(t *testing.T) {
	mock.StartMockService("9089")

	bank := caixa.New()
	input := newStubBoletoRequestDb(models.Caixa).Build()
	resp, err := bank.ProcessBoleto(input)
	assert.Nil(t, err)

	boView := models.NewBoletoView(*input, resp, bank.GetBankNameIntegration())
	// Set values in order to force attributes that can be checked
	oldID := boView.ID.Hex()
	oldPk := boView.PublicKey
	newID := "60be988b3193b131b8061835"
	newPk := "4f5316763275cc72e6313017171512a5c0d6f4710436a8a9cb506c36655be7f2"
	boView.ID, _ = primitive.ObjectIDFromHex(newID)
	boView.UID = "b32ace53-c7dc-11eb-86b5-00059a3c7a00"
	boView.SecretKey = "b32ace53-c7dc-11eb-86b5-00059a3c7a00"
	boView.PublicKey = newPk
	boView.Boleto.RequestKey = "b3279a06-c7dc-11eb-86b5-00059a3c7a00"
	boView.Links[0].Href = strings.ReplaceAll(boView.Links[0].Href, oldID, newID)
	boView.Links[0].Href = strings.ReplaceAll(boView.Links[0].Href, oldPk, newPk)
	boView.Links[1].Href = strings.ReplaceAll(boView.Links[1].Href, oldID, newID)
	boView.Links[1].Href = strings.ReplaceAll(boView.Links[1].Href, oldPk, newPk)

	expectedExpireDateTime := time.Now().Add(5 * time.Hour * 24)
	boView.Boleto.Title.ExpireDateTime = expectedExpireDateTime
	boView.Boleto.Title.ExpireDate = expectedExpireDateTime.Format("2006-01-02")

	expectedCreateDate := time.Now()
	boView.CreateDate = expectedCreateDate

	mID := boView.ID.Hex()
	resp.ID = string(mID)
	resp.Links = boView.Links

	err = deleteBoletoById(newID)
	assert.Nil(t, err)

	err = db.SaveBoleto(boView)
	assert.Nil(t, err)

	b, _, err := db.GetBoletoByID(newID, newPk)
	assert.Nil(t, err)
	createDateZone, _ := b.CreateDate.Zone()
	expireDateTimeZone, _ := b.Boleto.Title.ExpireDateTime.Zone()
	titleCreateDateZone, _ := b.Boleto.Title.CreateDate.Zone()

	assert.Equal(t, "60be988b3193b131b8061835", b.ID.Hex(), fmt.Sprintf("RecoveryBoleto ID error: expected [%s] got [%s]", "60be988b3193b131b8061835", b.ID.Hex()))
	assert.Equal(t, "b32ace53-c7dc-11eb-86b5-00059a3c7a00", b.UID, fmt.Sprintf("RecoveryBoleto UID error: expected [%s] got [%s]", "b32ace53-c7dc-11eb-86b5-00059a3c7a00", b.UID))
	assert.Equal(t, "b32ace53-c7dc-11eb-86b5-00059a3c7a00", b.SecretKey, fmt.Sprintf("RecoveryBoleto SecretKey error: expected [%s] got [%s]", "b32ace53-c7dc-11eb-86b5-00059a3c7a00", b.SecretKey))
	assert.Equal(t, "4f5316763275cc72e6313017171512a5c0d6f4710436a8a9cb506c36655be7f2", b.PublicKey, fmt.Sprintf("RecoveryBoleto PublicKey error: expected [%s] got [%s]", "4f5316763275cc72e6313017171512a5c0d6f4710436a8a9cb506c36655be7f2", b.PublicKey))
	assert.Equal(t, uint(123456), b.Boleto.Agreement.AgreementNumber, fmt.Sprintf("RecoveryBoleto Boleto.Agreement.AgreementNumber error: expected [%d] got [%d]", uint(123456), b.Boleto.Agreement.AgreementNumber))
	assert.Equal(t, uint16(0), b.Boleto.Agreement.Wallet, fmt.Sprintf("RecoveryBoleto Boleto.Agreement.Wallet error: expected [%d] got [%d]", uint16(0), b.Boleto.Agreement.Wallet))
	assert.Equal(t, "1234", b.Boleto.Agreement.Agency, fmt.Sprintf("RecoveryBoleto Boleto.Agreement.Agency error: expected [%s] got [%s]", "1234", b.Boleto.Agreement.Agency))
	assert.Equal(t, "", b.Boleto.Agreement.AgencyDigit, fmt.Sprintf("RecoveryBoleto Boleto.Agreement.AgencyDigit error: expected [%s] got [%s]", "", b.Boleto.Agreement.AgencyDigit))
	assert.Equal(t, "", b.Boleto.Agreement.Account, fmt.Sprintf("RecoveryBoleto Boleto.Agreement.Account error: expected [%s] got [%s]", "01234567", b.Boleto.Agreement.Account))
	assert.Equal(t, "", b.Boleto.Agreement.AccountDigit, fmt.Sprintf("RecoveryBoleto Boleto.Agreement.AccountDigit error: expected [%s] got [%s]", "", b.Boleto.Agreement.AccountDigit))
	// Check b.Boleto.Title.CreateDate attributes
	assert.Equal(t, boView.Boleto.Title.CreateDate.Unix(), b.Boleto.Title.CreateDate.Unix(), fmt.Sprintf("RecoveryBoleto Boleto.Title.CreateDate[timestamp] error: expected [%s] got [%s]", boView.Boleto.Title.CreateDate, b.Boleto.Title.CreateDate))
	assert.Equal(t, "-03", titleCreateDateZone, fmt.Sprintf("RecoveryBoleto Boleto.Title.CreateDate error: expected [%s] got [%s]", "-03", titleCreateDateZone))
	assert.Equal(t, "Local", b.Boleto.Title.CreateDate.Location().String(), fmt.Sprintf("RecoveryBoleto Boleto.Title.CreateDate[Location] error: expected [%s] got [%s]", "Local", b.Boleto.Title.CreateDate.Location().String()))
	// Check b.Boleto.Title.ExpireDateTime attributes
	assert.Equal(t, expectedExpireDateTime.Unix(), b.Boleto.Title.ExpireDateTime.Unix(), fmt.Sprintf("RecoveryBoleto Boleto.Title.ExpireDateTime[timestamp] error: expected [%s] got [%s]", expectedExpireDateTime, b.Boleto.Title.ExpireDateTime))
	assert.Equal(t, "-03", expireDateTimeZone, fmt.Sprintf("RecoveryBoleto Boleto.Title.ExpireDateTime[Zone] error: expected [%s] got [%s]", "-03", expireDateTimeZone))
	assert.Equal(t, "Local", b.Boleto.Title.ExpireDateTime.Location().String(), fmt.Sprintf("RecoveryBoleto Boleto.Title.ExpireDateTime[location] error: expected [%s] got [%s]", "Local", b.Boleto.Title.ExpireDateTime.Location().String()))
	assert.Equal(t, expectedExpireDateTime.Format("2006-01-02"), b.Boleto.Title.ExpireDate, fmt.Sprintf("RecoveryBoleto Boleto.Title.ExpireDate error: expected [%s] got [%s]", expectedExpireDateTime.Format("2006-01-02"), b.Boleto.Title.ExpireDate))
	assert.Equal(t, uint64(200), b.Boleto.Title.AmountInCents, fmt.Sprintf("RecoveryBoleto Boleto.Title.AmountInCents error: expected [%d] got [%d]", uint64(200), b.Boleto.Title.AmountInCents))
	assert.Equal(t, uint(14012345678901234), b.Boleto.Title.OurNumber, fmt.Sprintf("RecoveryBoleto Boleto.Title.OurNumber error: expected [%d] got [%d]", uint(14012345678901234), b.Boleto.Title.OurNumber))
	assert.Equal(t, "Campo de instruções -  max 40 caracteres", b.Boleto.Title.Instructions, fmt.Sprintf("RecoveryBoleto Boleto.Title.Instructions error: expected [%s] got [%s]", "Campo de instruções -  max 40 caracteres", b.Boleto.Title.Instructions))
	assert.Equal(t, "1234567890A", b.Boleto.Title.DocumentNumber, fmt.Sprintf("RecoveryBoleto Boleto.Title.DocumentNumber error: expected [%s] got [%s]", "1234567890A", b.Boleto.Title.DocumentNumber))
	assert.Equal(t, "OUT", b.Boleto.Title.BoletoType, fmt.Sprintf("RecoveryBoleto Boleto.Title.BoletoType error: expected [%s] got [%s]", "OUT", b.Boleto.Title.BoletoType))
	assert.Equal(t, "99", b.Boleto.Title.BoletoTypeCode, fmt.Sprintf("RecoveryBoleto Boleto.Title.BoletoTypeCode error: expected [%s] got [%s]", "99", b.Boleto.Title.BoletoTypeCode))
	assert.Equal(t, "", b.Boleto.Recipient.Name, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Name error: expected [%s] got [%s]", "", b.Boleto.Recipient.Name))
	assert.Equal(t, "CNPJ", b.Boleto.Recipient.Document.Type, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Document.Type error: expected [%s] got [%s]", "CNPJ", b.Boleto.Recipient.Document.Type))
	assert.Equal(t, "12123123000112", b.Boleto.Recipient.Document.Number, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Document.Number error: expected [%s] got [%s]", "12123123000112", b.Boleto.Recipient.Document.Number))
	assert.Equal(t, "", b.Boleto.Recipient.Address.Street, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Address.Street error: expected [%s] got [%s]", "", b.Boleto.Recipient.Address.Street))
	assert.Equal(t, "", b.Boleto.Recipient.Address.Number, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Address.Number error: expected [%s] got [%s]", "", b.Boleto.Recipient.Address.Number))
	assert.Equal(t, "", b.Boleto.Recipient.Address.Complement, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Address.Complement error: expected [%s] got [%s]", "", b.Boleto.Recipient.Address.Complement))
	assert.Equal(t, "", b.Boleto.Recipient.Address.ZipCode, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Address.ZipCode error: expected [%s] got [%s]", "", b.Boleto.Recipient.Address.ZipCode))
	assert.Equal(t, "", b.Boleto.Recipient.Address.City, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Address.City error: expected [%s] got [%s]", "", b.Boleto.Recipient.Address.City))
	assert.Equal(t, "", b.Boleto.Recipient.Address.District, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Address.District error: expected [%s] got [%s]", "", b.Boleto.Recipient.Address.District))
	assert.Equal(t, "", b.Boleto.Recipient.Address.StateCode, fmt.Sprintf("RecoveryBoleto Boleto.Recipient.Address.StateCode error: expected [%s] got [%s]", "", b.Boleto.Recipient.Address.StateCode))
	assert.Equal(t, "Willian Jadson Bezerra Menezes Tupinambá", b.Boleto.Buyer.Name, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Name error: expected [%s] got [%s]", "Willian Jadson Bezerra Menezes Tupinambá", b.Boleto.Buyer.Name))
	assert.Equal(t, "CPF", b.Boleto.Buyer.Document.Type, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Document.Type error: expected [%s] got [%s]", "CPF", b.Boleto.Buyer.Document.Type))
	assert.Equal(t, "12312312312", b.Boleto.Buyer.Document.Number, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Document.Number error: expected [%s] got [%s]", "12312312312", b.Boleto.Buyer.Document.Number))
	assert.Equal(t, "Rua da Assunção de Sá", b.Boleto.Buyer.Address.Street, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Address.Street error: expected [%s] got [%s]", "Rua da Assunção de Sá", b.Boleto.Buyer.Address.Street))
	assert.Equal(t, "123", b.Boleto.Buyer.Address.Number, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Address.Number error: expected [%s] got [%s]", "123", b.Boleto.Buyer.Address.Number))
	assert.Equal(t, "Seção A, s 02", b.Boleto.Buyer.Address.Complement, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Address.Complement error: expected [%s] got [%s]", "Seção A, s 02", b.Boleto.Buyer.Address.Complement))
	assert.Equal(t, "20520051", b.Boleto.Buyer.Address.ZipCode, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Address.ZipCode error: expected [%s] got [%s]", "20520051", b.Boleto.Buyer.Address.ZipCode))
	assert.Equal(t, "Belém do Pará", b.Boleto.Buyer.Address.City, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Address.City error: expected [%s] got [%s]", "Belém do Pará", b.Boleto.Buyer.Address.City))
	assert.Equal(t, "Açaí", b.Boleto.Buyer.Address.District, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Address.District error: expected [%s] got [%s]", "Açaí", b.Boleto.Buyer.Address.District))
	assert.Equal(t, "PA", b.Boleto.Buyer.Address.StateCode, fmt.Sprintf("RecoveryBoleto Boleto.Buyer.Address.StateCode error: expected [%s] got [%s]", "PA", b.Boleto.Buyer.Address.StateCode))
	assert.Equal(t, models.BankNumber(104), b.Boleto.BankNumber, fmt.Sprintf("RecoveryBoleto Boleto.BankNumber error: expected [%d] got [%d]", models.BankNumber(104), b.Boleto.BankNumber))
	assert.Equal(t, "b3279a06-c7dc-11eb-86b5-00059a3c7a00", b.Boleto.RequestKey, fmt.Sprintf("RecoveryBoleto Boleto.RequestKey error: expected [%s] got [%s]", "b3279a06-c7dc-11eb-86b5-00059a3c7a00", b.Boleto.RequestKey))
	assert.Equal(t, models.BankNumber(104), b.BankID, fmt.Sprintf("RecoveryBoleto BankID error: expected [%d] got [%d]", models.BankNumber(104), b.BankID))
	// Check CreateDate attributes
	assert.Equal(t, expectedCreateDate.Unix(), b.CreateDate.Unix(), fmt.Sprintf("RecoveryBoleto CreateDate error: expected [%s] got [%s]", expectedCreateDate, b.CreateDate))
	assert.Equal(t, "-03", createDateZone, fmt.Sprintf("RecoveryBoleto CreateDate error: expected [%s] got [%s]", "-03", createDateZone))
	assert.Equal(t, "Local", b.CreateDate.Location().String(), fmt.Sprintf("RecoveryBoleto CreateDate error: expected [%s] got [%s]", "Local", b.CreateDate.Location().String()))
	assert.Equal(t, "104-0", b.BankNumber, fmt.Sprintf("RecoveryBoleto BankNumber error: expected [%s] got [%s]", "104-0", b.BankNumber))
	assert.Equal(t, "10492.00650 61000.100042 09922.269841 3 72670000001000", b.DigitableLine, fmt.Sprintf("RecoveryBoleto DigitableLine error: expected [%s] got [%s]", "10492.00650 61000.100042 09922.269841 3 72670000001000", b.DigitableLine))
	assert.Equal(t, "10493726700000010002006561000100040992226984", b.Barcode, fmt.Sprintf("RecoveryBoleto Barcode error: expected [%s] got [%s]", "10493726700000010002006561000100040992226984", b.Barcode))
	assert.Equal(t, "http://localhost:3000/boleto?fmt=html&id=60be988b3193b131b8061835&pk=4f5316763275cc72e6313017171512a5c0d6f4710436a8a9cb506c36655be7f2", b.Links[0].Href, fmt.Sprintf("RecoveryBoleto Links[0].Href error: expected [%s] got [%s]", "http://localhost:3000/boleto?fmt=html&id=60be988b3193b131b8061835&pk=4f5316763275cc72e6313017171512a5c0d6f4710436a8a9cb506c36655be7f2", b.Links[0].Href))
	assert.Equal(t, "html", b.Links[0].Rel, fmt.Sprintf("RecoveryBoleto Links[0].Rel error: expected [%s] got [%s]", "html", b.Links[0].Rel))
	assert.Equal(t, "GET", b.Links[0].Method, fmt.Sprintf("RecoveryBoleto Links[0].Method error: expected [%s] got [%s]", "GET", b.Links[0].Method))
	assert.Equal(t, "http://localhost:3000/boleto?fmt=pdf&id=60be988b3193b131b8061835&pk=4f5316763275cc72e6313017171512a5c0d6f4710436a8a9cb506c36655be7f2", b.Links[1].Href, fmt.Sprintf("RecoveryBoleto Links[1].Href error: expected [%s] got [%s]", "http://localhost:3000/boleto?fmt=pdf&id=60be988b3193b131b8061835&pk=4f5316763275cc72e6313017171512a5c0d6f4710436a8a9cb506c36655be7f2", b.Links[1].Href))
	assert.Equal(t, "pdf", b.Links[1].Rel, fmt.Sprintf("RecoveryBoleto Links[1].Rel error: expected [%s] got [%s]", "pdf", b.Links[1].Rel))
	assert.Equal(t, "GET", b.Links[1].Method, fmt.Sprintf("RecoveryBoleto Links[1].Method error: expected [%s] got [%s]", "GET", b.Links[1].Method))
}

func TestMongoDb_GetUserCredentials(t *testing.T) {
	mock.StartMockService("9089")

	gandalfID := "60c293944808daa6fdf2f3b1"
	gID, _ := primitive.ObjectIDFromHex(gandalfID)
	cGandalf := models.Credentials{
		ID:       gID,
		Username: "gandalf",
		Password: "grey",
	}
	deleteCredentialById(gandalfID)
	err := saveCredential(cGandalf)
	assert.Nil(t, err)

	sarumanID := "60c293944808daa6fdf2f3b3"
	sID, _ := primitive.ObjectIDFromHex(sarumanID)
	cSaruman := models.Credentials{
		ID:       sID,
		Username: "saruman",
		Password: "white",
	}
	deleteCredentialById(sarumanID)
	err = saveCredential(cSaruman)
	assert.Nil(t, err)

	c, err := db.GetUserCredentials()
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(c), 2)

	g, err := getUserCredentialByID(gandalfID)
	assert.Nil(t, err)

	assert.Equal(t, "60c293944808daa6fdf2f3b1", g.ID.Hex(), fmt.Sprintf("Gandalf ID error: expected [%s] got [%s]", "60c293944808daa6fdf2f3b1", g.ID.Hex()))
	assert.Equal(t, "gandalf", g.Username, fmt.Sprintf("Gandalf Username error: expected [%s] got [%s]", "gandalf", g.Username))
	assert.Equal(t, "grey", g.Password, fmt.Sprintf("Gandalf Password error: expected [%s] got [%s]", "grey", g.Password))

	s, err := getUserCredentialByID(sarumanID)
	assert.Nil(t, err)
	assert.Equal(t, "60c293944808daa6fdf2f3b3", s.ID.Hex(), fmt.Sprintf("Saruman ID error: expected [%s] got [%s]", "60c293944808daa6fdf2f3b3", s.ID.Hex()))
	assert.Equal(t, "saruman", s.Username, fmt.Sprintf("Saruman Username error: expected [%s] got [%s]", "saruman", s.Username))
	assert.Equal(t, "white", s.Password, fmt.Sprintf("Saruman Password error: expected [%s] got [%s]", "white", s.Password))

	deleteCredentialById(gandalfID)
	deleteCredentialById(sarumanID)
}

package api

import (
	"net/http"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gin-gonic/gin"

	"strings"

	"encoding/json"
	"os"

	"io/ioutil"

	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/boleto"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

//Regista um boleto em um determinado banco
func registerBoleto(c *gin.Context) {
	_boleto, _ := c.Get("boleto")
	boleto := _boleto.(models.BoletoRequest)
	bank, err := bank.Get(boleto.BankNumber)
	if checkError(c, err, log.CreateLog()) {
		return
	}
	lg := bank.Log()
	lg.Operation = "RegisterBoleto"
	lg.NossoNumero = boleto.Title.OurNumber
	lg.Recipient = boleto.Recipient.Name
	lg.RequestKey = boleto.RequestKey
	lg.BankName = bank.GetBankNumber().BankName()

	repo, err := db.GetDB()
	if checkError(c, err, lg) {
		return
	}
	resp, errR := bank.ProcessBoleto(&boleto)
	if checkError(c, errR, lg) {
		return
	}
	st := http.StatusOK
	if len(resp.Errors) > 0 {
		if resp.StatusCode > 0 {
			st = resp.StatusCode
		} else {
			st = http.StatusBadRequest
		}
	} else {
		boView := models.NewBoletoView(boleto, resp)
		resp.ID = boView.ID
		resp.Links = boView.Links
		errMongo := repo.SaveBoleto(boView)
		if errMongo != nil {
			saveBoletoJSONFile(boView, lg, errMongo)
		}
	}
	c.JSON(st, resp)
	c.Set("boletoResponse", resp)
}

func saveBoletoJSONFile(boView models.BoletoView, lg *log.Log, err error) {
	lg.Warn(err.Error(), "Boleto cannot be saved at Database")
	fd, errOpen := os.Create(config.Get().BoletoJSONFileStore + "/boleto_" + boView.UID + ".json")
	if errOpen != nil {
		lg.Fatal(boView, "[BOLETO_ONLINE_CONTINGENCIA]"+errOpen.Error())
	}
	data, _ := json.Marshal(boView)
	_, errW := fd.Write(data)
	if errW != nil {
		lg.Fatal(boView, "[BOLETO_ONLINE_CONTINGENCIA]"+errW.Error())
	}
	fd.Close()
}

func getBoleto(c *gin.Context) {
	c.Status(200)

	id := c.Query("id")
	format := c.Query("fmt")
	repo, errCon := db.GetDB()
	if checkError(c, errCon, log.CreateLog()) {
		return
	}
	bleto, err := repo.GetBoletoByID(id)
	if err != nil {
		uid := util.Decrypt(id)
		fd, err := os.Open(config.Get().BoletoJSONFileStore + "/boleto_" + uid + ".json")
		if err != nil {
			checkError(c, models.NewHTTPNotFound("Boleto não encontrado na base de dados", "MP404"), log.CreateLog())
			return
		}
		data, errR := ioutil.ReadAll(fd)
		if errR != nil {
			checkError(c, models.NewHTTPNotFound("Boleto não encontrado na base de dados", "MP404"), log.CreateLog())
			return
		}
		json.Unmarshal(data, &bleto)
		fd.Close()
	}

	s, err := boleto.HTML(bleto, format)
	if checkError(c, err, log.CreateLog()) {
		return
	}
	if format == "html" {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteString(s)
	} else {
		c.Header("Content-Type", "application/pdf")
		buf, _ := toPdf(s)
		c.Writer.Write(buf)
	}

}

func toPdf(page string) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	pdfg.Dpi.Set(600)
	pdfg.NoCollate.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(page)))
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}
	return pdfg.Bytes(), nil
}

func getBoletoByID(c *gin.Context) {
	id := c.Param("id")
	db, errDb := db.GetDB()
	if errDb != nil {
		checkError(c, models.NewInternalServerError("MP500", "Erro interno"), log.CreateLog())
	}
	boleto, err := db.GetBoletoByID(id)
	if err != nil {
		checkError(c, models.NewHTTPNotFound("MP404", "Boleto não encontrado"), nil)
		return
	}
	c.JSON(http.StatusOK, boleto)
}

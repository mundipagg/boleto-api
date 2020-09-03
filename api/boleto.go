package api

import (
	"encoding/json"
	"errors"
	"github.com/mundipagg/boleto-api/queue"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	jm "github.com/tdewolff/minify/json"

	"strings"

	"fmt"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/boleto"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"io/ioutil"
)

//Regista um boleto em um determinado banco
func registerBoleto(c *gin.Context) {

	if _, hasErr := c.Get("error"); hasErr {
		return
	}

	_user, _ := c.Get("serviceuser")
	_boleto, _ := c.Get("boleto")
	_bank, _ := c.Get("bank")
	bol := _boleto.(models.BoletoRequest)
	bank := _bank.(bank.Bank)

	lg := bank.Log()
	lg.Operation = "RegisterBoleto"
	lg.NossoNumero = bol.Title.OurNumber
	lg.Recipient = bol.Recipient.Name
	lg.RequestKey = bol.RequestKey
	lg.BankName = bank.GetBankNameIntegration()
	lg.IPAddress = c.ClientIP()
	lg.ServiceUser = _user.(string)
	resp, err := bank.ProcessBoleto(&bol)
	if checkError(c, err, lg) {
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
		mongo, errMongo := db.CreateMongo(lg)

		boView := models.NewBoletoView(bol, resp, bank.GetBankNameIntegration())
		mID, _ := boView.ID.MarshalText()
		resp.ID = string(mID)
		resp.Links = boView.Links

		redis := db.CreateRedis()

		if errMongo == nil {
			errMongo = mongo.SaveBoleto(boView)
		}

		errMongo = errors.New("teste")

		if errMongo != nil {
			lg.Warn(errMongo.Error(), fmt.Sprintf("Error saving to mongo - %s", errMongo.Error()))
			b := minifyJSON(boView)
			p := queue.NewPublisher(b)

			if !queue.WriteMessage(p) {
				err = redis.SetBoletoJSON(b, resp.ID, boView.PublicKey, lg)
				if checkError(c, err, lg) {
					return
				}
			}
		}

		defer queue.CloseConnection()

		bhtml, _ := boleto.HTML(boView, "html")
		s := minifyString(bhtml, "text/html")
		redis.SetBoletoHTML(s, resp.ID, boView.PublicKey, lg)
	}
	c.JSON(st, resp)
	c.Set("boletoResponse", resp)
}

func getBoleto(c *gin.Context) {
	c.Status(200)

	id := c.Query("id")
	format := c.Query("fmt")
	pk := c.Query("pk")

	log := log.CreateLog()
	log.Operation = "GetBoleto"
	log.IPAddress = c.ClientIP()

	redis := db.CreateRedis()
	b := redis.GetBoletoHTMLByID(id, pk, log)

	if b == "" {
		mongo, errMongo := db.CreateMongo(log)
		if checkError(c, errMongo, log) {
			return
		}

		boView, err := mongo.GetBoletoByID(id, pk)

		if err != nil && err.Error() == db.NotFoundDoc {
			e := fmt.Sprintf("%s - %s", err.Error(), c.Request.RequestURI)
			log.Warn(e, fmt.Sprintf("Boleto notfound - %s", id))
			checkError(c, models.NewHTTPNotFound("MP404", "Not Found"), log)
			return
		} else if err != nil && err.Error() == db.InvalidPK {
			e := fmt.Sprintf("%s - %s", err.Error(), c.Request.RequestURI)
			log.Warn(e, fmt.Sprintf("Boleto with invalid pk - %s", id))
			checkError(c, models.NewHTTPNotFound("MP404", "Not Found"), log)
			return
		} else if err != nil {
			checkError(c, models.NewInternalServerError("MP500", err.Error()), log)
			return
		}

		bhtml, err := boleto.HTML(boView, "html")
		b = minifyString(bhtml, "text/html")
	}

	if format == "html" {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteString(b)
	} else {
		c.Header("Content-Type", "application/pdf")
		bpdf, err := toPdf(b)

		if err != nil {
			c.Header("Content-Type", "application/json")
			checkError(c, models.NewInternalServerError(err.Error(), "Internal error"), log)
			c.Abort()
		} else {
			c.Writer.Write(bpdf)
		}

	}

}

func toPdf(page string) ([]byte, error) {
	url := config.Get().PdfAPIURL
	payload := strings.NewReader(page)
	if req, err := http.NewRequest("POST", url, payload); err != nil {
		return nil, err
	} else if res, err := http.DefaultClient.Do(req); err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
		return ioutil.ReadAll(res.Body)
	}
}

func getBoletoByID(c *gin.Context) {
	id := c.Param("id")
	pk := c.Param("pk")
	log := log.CreateLog()
	log.Operation = "GetBoletoV1"

	mongo, errDb := db.CreateMongo(log)
	if errDb != nil {
		checkError(c, models.NewInternalServerError("MP500", "Internal error"), log)
	}
	boleto, err := mongo.GetBoletoByID(id, pk)
	if err != nil {
		checkError(c, models.NewHTTPNotFound("MP404", "Boleto não encontrado"), nil)
		return
	}
	c.JSON(http.StatusOK, boleto)
}

//minifyJSON converte um model BoletoView para um JSON/STRING
func minifyJSON(m models.BoletoView) string {
	j, _ := json.Marshal(m)

	return minifyString(string(j), "application/json")
}

func minifyString(mString, tp string) string {
	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepDocumentTags:        true,
		KeepEndTags:             true,
		KeepWhitespace:          false,
		KeepConditionalComments: true,
	})
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("application/json", jm.Minify)

	s, err := m.String(tp, mString)

	if err != nil {
		return mString
	}

	return s
}

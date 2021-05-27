package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mundipagg/boleto-api/queue"

	"github.com/gin-gonic/gin"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	jm "github.com/tdewolff/minify/json"

	"strings"

	"fmt"
	"io/ioutil"

	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/boleto"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
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

		bhtml, _ := boleto.HTML(boView, "html")
		s := minifyString(bhtml, "text/html")
		redis.SetBoletoHTML(s, resp.ID, boView.PublicKey, lg)
	}
	c.JSON(st, resp)
	c.Set("boletoResponse", resp)
}

func getBoleto(c *gin.Context) {
	start := time.Now()
	var boletoHtml string

	c.Status(200)
	log := log.CreateLog()
	log.Operation = "GetBoleto"
	log.IPAddress = c.ClientIP()

	var result = models.NewGetBoletoResult(c)

	defer logResult(result, log, start)

	if !result.HasValidKeys() {
		result.SetErrorResponse(c, models.NewErrorResponse("MP404", "Not Found"), http.StatusNotFound)
		result.LogSeverity = "Warning"
		return
	}

	redis := db.CreateRedis()
	boletoHtml, result.CacheElapsedTimeInMilliseconds = redis.GetBoletoHTMLByID(result.Id, result.PrivateKey, log)

	if boletoHtml == "" {
		var err error
		var boView models.BoletoView
		mongo, errMongo := db.CreateMongo(log)

		if errMongo != nil {
			result.SetErrorResponse(c, models.NewErrorResponse("MP500", errMongo.Error()), http.StatusInternalServerError)
			result.LogSeverity = "Error"
			return
		}

		boView, result.DatabaseElapsedTimeInMilliseconds, err = mongo.GetBoletoByID(result.Id, result.PrivateKey)

		if err != nil && (err.Error() == db.NotFoundDoc || err.Error() == db.InvalidPK) {
			result.SetErrorResponse(c, models.NewErrorResponse("MP404", "Not Found"), http.StatusNotFound)
			result.LogSeverity = "Warning"
			return
		} else if err != nil {
			result.SetErrorResponse(c, models.NewErrorResponse("MP500", err.Error()), http.StatusInternalServerError)
			result.LogSeverity = "Error"
			return
		}
		result.BoletoSource = "mongo"
		html, err := boleto.HTML(boView, "html")
		boletoHtml = minifyString(html, "text/html")
	} else {
		result.BoletoSource = "redis"
	}

	if result.Format == "html" {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteString(boletoHtml)
	} else {
		c.Header("Content-Type", "application/pdf")
		if boletoPdf, err := toPdf(boletoHtml); err == nil {
			c.Writer.Write(boletoPdf)
		} else {
			c.Header("Content-Type", "application/json")
			result.SetErrorResponse(c, models.NewErrorResponse("MP500", err.Error()), http.StatusInternalServerError)
			result.LogSeverity = "Error"
			return
		}
	}

	result.LogSeverity = "Information"
}

func logResult(result *models.GetBoletoResult, log *log.Log, start time.Time) {
	result.TotalElapsedTimeInMilliseconds = time.Since(start).Milliseconds()
	log.GetBoleto(result, result.LogSeverity)
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
	boleto, _, err := mongo.GetBoletoByID(id, pk)
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

package api

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"

	"github.com/mundipagg/boleto-api/models"
)

//InstallRestAPI "instala" e sobe o servico de rest
func InstallRestAPI() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(executionController())
	if config.Get().DevMode && !config.Get().MockMode {
		router.Use(gin.Logger())
	}
	InstallV1(router)
	router.GET("/boleto", getBoleto)
	router.GET("/boleto/confirmation", confirmation)
	router.POST("/boleto/confirmation", confirmation)
	router.Run(config.Get().APIPort)
}

func confirmation(c *gin.Context) {
	if dump, err := httputil.DumpRequest(c.Request, true); err == nil {
		l := log.CreateLog()
		l.BankName = "Bradesco"
		l.Operation = "BoletoConfirmation"
		l.Request(string(dump), c.Request.URL.String(), nil)
	}
	c.String(200, "OK")
}

func checkError(c *gin.Context, err error, l *log.Log) bool {

	if err != nil {
		errResp := models.BoletoResponse{
			Errors: models.NewErrors(),
		}
		switch v := err.(type) {
		case models.IErrorResponse:
			errResp.Errors.Append(v.ErrorCode(), v.Error())
			c.JSON(http.StatusBadRequest, errResp)
		case models.IHttpNotFound:
			errResp.Errors.Append("MP404", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusNotFound, errResp)
		case models.IFormatError:
			errResp.Errors.Append("MP400", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusBadRequest, errResp)
		case models.IServerError:
			errResp.Errors.Append("MP500", "Internal Error")
			l.Fatal(v.Error(), v.Message())
			c.JSON(http.StatusInternalServerError, errResp)
		default:
			l.Fatal(err.Error(), "")
			errResp.Errors.Append("MP500", "Internal Error")
			c.JSON(http.StatusInternalServerError, errResp)
		}
		return true
	}
	return false
}

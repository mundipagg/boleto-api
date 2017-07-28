package api

import (
	"fmt"
	"net/http"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"

	"os"

	"github.com/mundipagg/boleto-api/models"
	gin "gopkg.in/gin-gonic/gin.v1"
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
	if config.Get().HTTPOnly || config.Get().DevMode {
		router.Run(config.Get().APIPort)
	} else {
		err := router.RunTLS(config.Get().APIPort, config.Get().TLSCertPath, config.Get().TLSKeyPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
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
			c.JSON(http.StatusNotFound, errResp)
		case models.IFormatError:
			errResp.Errors.Append("MP400", v.Error())
			c.JSON(http.StatusBadRequest, errResp)
		case models.IServerError:
			errResp.Errors.Append("MP500", "Erro interno")
			l.Fatal(v.Error(), v.Message())
			c.JSON(http.StatusInternalServerError, errResp)
		default:
			l.Fatal(err.Error(), "")
			errResp.Errors.Append("MP500", "Erro interno")
			c.JSON(http.StatusInternalServerError, errResp)
		}
		return true
	}
	return false
}

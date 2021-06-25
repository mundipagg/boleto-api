package api

import (
	"context"
	stdlog "log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/mundipagg/boleto-api/metrics"

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
	useNewRelic(router)

	if config.Get().DevMode && !config.Get().MockMode {
		router.Use(gin.Logger())
	}
	InstallV1(router)
	router.StaticFile("/favicon.ico", "./boleto/favicon.ico")
	router.GET("/boleto/memory-check/:unit", memory)
	router.GET("/boleto/memory-check/", memory)
	router.GET("/boleto", getBoleto)
	router.GET("/boleto/confirmation", confirmation)
	router.POST("/boleto/confirmation", confirmation)
	router.GET("/test-shutdown", func(c *gin.Context) {
		stdlog.Println("Test Gracefully Shutdown")
		time.Sleep(5*time.Second)
		c.Status(200)
	})

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    config.Get().APIPort,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			interrupt <- syscall.SIGTERM
			stdlog.Println("err: ", err)
		}
	}()

	<-interrupt
	stdlog.Println("shutdown server")
	err := server.Shutdown(context.Background())
	if err != nil {
		stdlog.Print(err)
	}

	stdlog.Println("shutdown completed")
}

func useNewRelic(router *gin.Engine) {
	if !config.Get().TelemetryEnabled {
		return
	}

	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName(config.Get().NewRelicAppName),
		newrelic.ConfigLicense(config.Get().NewRelicLicence),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	router.Use(nrgin.Middleware(app))
}

func memory(c *gin.Context) {
	unit := c.Param("unit")
	c.JSON(200, metrics.GetMemoryReport(unit))
}

func confirmation(c *gin.Context) {
	if dump, err := httputil.DumpRequest(c.Request, true); err == nil {
		l := log.CreateLog()
		l.BankName = "BradescoShopFacil"
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

		case models.ErrorResponse:
			errResp.Errors.Append(v.ErrorCode(), v.Error())
			c.JSON(http.StatusBadRequest, errResp)

		case models.HttpNotFound:
			errResp.Errors.Append("MP404", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusNotFound, errResp)

		case models.InternalServerError:
			errResp.Errors.Append("MP500", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusInternalServerError, errResp)

		case models.BadGatewayError:
			errResp.Errors.Append("MP502", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusBadGateway, errResp)

		case models.FormatError:
			errResp.Errors.Append("MP400", v.Error())
			l.Warn(errResp, v.Error())
			c.JSON(http.StatusBadRequest, errResp)

		default:
			errResp.Errors.Append("MP500", "Internal Error")
			l.Fatal(errResp, v.Error())
			c.JSON(http.StatusInternalServerError, errResp)
		}

		c.Set("boletoResponse", errResp)
		return true
	}
	return false
}

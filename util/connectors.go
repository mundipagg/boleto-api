package util

import (
	"fmt"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"

	"errors"
	"net/http"
	"strings"

	"github.com/PMoneda/flow"
)

// SeqLogConector é um connector flow para logar no Seq
func SeqLogConector(e *flow.ExchangeMessage, u flow.URI, params ...interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var b string
	switch t := e.GetBody().(type) {
	case string:
		if t == "" {
			b = "Nenhum retorno do serviço"
		} else {
			b = t
		}
	case error:
		b = t.Error()
	default:
		b = fmt.Sprintln(t)
	}
	if len(params) > 0 {
		l := params[0].(*log.Log)
		if u.GetOption("type") == "request" {
			l.Request(b, u.GetOption("url"), e.GetHeaderMap())
		}
		if u.GetOption("type") == "response" {
			l.Response(b, u.GetOption("url"))
		}
	}
	return nil
}

//TlsConector is a connector to send https request client certificate Params[0] *http.Transport (http.Transport configuration with certificate files config)
func TlsConector(e *flow.ExchangeMessage, u flow.URI, params ...interface{}) error {
	var b string
	switch t := e.GetBody().(type) {
	case string:
		if t == "" {
			b = "Nenhum retorno do serviço"
		} else {
			b = t
		}
	case error:
		b = t.Error()
	default:
		b = fmt.Sprintln(t)
	}
	if len(params) > 0 {
		switch t := params[0].(type) {
		case *http.Transport:
			var url string
			var response string
			var status int
			var err error
			if config.Get().MockMode {
				url = strings.Replace(u.GetRaw(), "tls", "http", 1)
				response, status, err = Post(url, b, e.GetHeaderMap())
			} else {
				url = strings.Replace(u.GetRaw(), "tls", "https", 1)
				response, status, err = PostTLS(url, b, e.GetHeaderMap(), t)
			}
			if err != nil {
				e.SetHeader("error", err.Error())
				e.SetBody(err)
				return err
			}
			e.SetHeader("status", fmt.Sprintf("%d", status))
			e.SetBody(response)
		default:
			return errors.New("invalid data type")
		}
	} else {
		return errors.New("Http Transport is required")
	}
	return nil
}

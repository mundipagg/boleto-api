package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/env"
)

func mockInstallApi() *gin.Engine {
	env.Config(true, true, true)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(executionController())
	Base(r)
	V1(r)
	V2(r)
	return r
}

package middleware

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/env"
)

func mockServerEngine() *gin.Engine {
	env.Config(true, true, true)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

func arrangeMiddlewareRoute(route string, handlers ...gin.HandlerFunc) (*gin.Engine, *httptest.ResponseRecorder) {
	router := mockServerEngine()
	router.POST(route, handlers...)
	w := httptest.NewRecorder()
	return router, w
}

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mundipagg/boleto-api/metrics"
	"github.com/stretchr/testify/assert"
)

func Test_Memory_WhenCalled_ReturnMemoryReportSuccessful(t *testing.T) {
	var resp = metrics.MemoryReport{}

	router := mockInstallApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/boleto/memory-check/", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Nil(t, json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&resp))
	assert.NotEmpty(t, resp)
	assert.Equal(t, "MB", resp.MemoryUnit)
}

func Test_Memory_WhenCalledWithKilobytesParameter_ReturnMemoryReportSuccessful(t *testing.T) {
	var resp = metrics.MemoryReport{}

	router := mockInstallApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/boleto/memory-check/KB", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Nil(t, json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&resp))
	assert.NotEmpty(t, resp)
	assert.Equal(t, "KB", resp.MemoryUnit)
}

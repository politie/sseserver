package main

import (
	"flag"
	"testing"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestReadData(t *testing.T) {
	flag.StringVar(&inputFile, "input-file", "./examples/test.json", "test default")
	flag.Parse()
	content := readData()
	expectedContent := "{  \"hello\": \"world\" }"

	assert.Equal(t, expectedContent, content)
}

func TestEndpointRouteWithPlainHttp(t *testing.T) {
	inputFileContent = "foobar"

	router := gin.New()
	endpointRoute(router.Group("/"))

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, inputFileContent, resp.Body.String())
}
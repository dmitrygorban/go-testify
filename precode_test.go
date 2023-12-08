package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func urlWithQuery(city string, count int) string {
	return fmt.Sprintf("/cafe?city=%s&count=%d", city, count)
}

func TestMainHandlerWithValidParams(t *testing.T) {
	totalCount := 4
	url := urlWithQuery("moscow", totalCount)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "handler returned wrong status code, expected 200")
	assert.NotEmpty(t, responseRecorder.Body.String(), "response body is empty")
}

func TestMainHandlerWithInvalidCityName(t *testing.T) {
	totalCount := 4
	url := urlWithQuery("saint-petersburg", totalCount)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "handler returned wrong status code, expected 400")

	expectedErrorBody := "wrong city value"

	assert.Equal(t, expectedErrorBody, responseRecorder.Body.String(), "unexpected error body")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 6
	url := urlWithQuery("moscow", totalCount)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)

	responseBody := responseRecorder.Body.String()
	response := strings.Split(responseBody, ",")

	expectedCafesLen := len(CafeList["moscow"])

	assert.Len(t, response, expectedCafesLen, "unexpected number of cafes in response")
}

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": {"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func urlWithQuery(city string, count int) string {
	return fmt.Sprintf("/cafe?city=%s&count=%d", city, count)
}

func TestMainHandlerWithValidParams(t *testing.T) {
	totalCount := 4
	url := urlWithQuery("moscow", totalCount)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
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
	handler := http.HandlerFunc(mainHandle)
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
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	responseBody := responseRecorder.Body.String()
	response := strings.Split(responseBody, ",")

	expectedCafesLen := len(cafeList["moscow"])

	assert.Len(t, response, expectedCafesLen, "unexpected number of cafes in response")
}

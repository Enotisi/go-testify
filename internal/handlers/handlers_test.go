package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeListTest = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

const (
	validCity   = "moscow"
	invalidCity = "izhevsk"
	count       = 2
	totalCount  = 4
)

// Проверка, что API отвечает 200 при валидном запросе
func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe", nil)

	q := req.URL.Query()
	q.Set("count", strconv.Itoa(count))
	q.Set("city", validCity)

	req.URL.RawQuery = q.Encode()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body.String())
	assert.Len(t, strings.Split(responseRecorder.Body.String(), ","), count)
}

// Проверка, что API отвечает 400 при невалидном городе
func TestMainHandlerWithoutCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe", nil)

	q := req.URL.Query()
	q.Set("count", strconv.Itoa(count))
	q.Set("city", invalidCity)

	req.URL.RawQuery = q.Encode()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body.String())
	assert.Equal(t, "wrong city value", responseRecorder.Body.String())
}

// Проверка, что API возвращает все кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe", nil)

	q := req.URL.Query()
	q.Set("count", strconv.Itoa(totalCount))
	q.Set("city", validCity)
	req.URL.RawQuery = q.Encode()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body.String())
	assert.Len(t, strings.Split(responseRecorder.Body.String(), ","), totalCount)
	assert.Equal(t, responseRecorder.Body.String(), strings.Join(cafeListTest[validCity], ","))
}

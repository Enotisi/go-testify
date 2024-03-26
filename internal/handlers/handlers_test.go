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
)

type handlerFuncType func(w http.ResponseWriter, req *http.Request)

func createGetReq(url string, query map[string]string) *http.Request {

	req := httptest.NewRequest(http.MethodGet, url, nil)

	if len(query) > 0 {

		q := req.URL.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return req
}

func createRecorder(handlerFunc handlerFuncType, req *http.Request) *httptest.ResponseRecorder {

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(responseRecorder, req)

	return responseRecorder
}

// Проверка, что API отвечает 200 при валидном запросе
func TestMainHandlerWhenOk(t *testing.T) {

	query := map[string]string{
		"count": strconv.Itoa(count),
		"city":  validCity,
	}
	req := createGetReq("/cafe", query)
	responseRecorder := createRecorder(MainHandle, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body.String())
	assert.Len(t, strings.Split(responseRecorder.Body.String(), ","), count)
}

// Проверка, что API отвечает 400 при невалидном городе
func TestMainHandlerWithoutCity(t *testing.T) {

	query := map[string]string{
		"count": strconv.Itoa(count),
		"city":  invalidCity,
	}
	req := createGetReq("/cafe", query)
	responseRecorder := createRecorder(MainHandle, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body.String())
	assert.Equal(t, "wrong city value", responseRecorder.Body.String())
}

// Проверка, что API возвращает все кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {

	totalCount := len(cafeListTest[validCity])
	query := map[string]string{
		"count": strconv.Itoa(totalCount + 5),
		"city":  validCity,
	}
	req := createGetReq("/cafe", query)
	responseRecorder := createRecorder(MainHandle, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body.String())
	assert.Len(t, strings.Split(responseRecorder.Body.String(), ","), totalCount)
	assert.Equal(t, responseRecorder.Body.String(), strings.Join(cafeListTest[validCity], ","))
}

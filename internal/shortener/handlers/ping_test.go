package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestAPIHandler_PingDB(t *testing.T) {

	testCases := []struct {
		name     string
		fields   fields
		expected expected
		PingOK   bool
	}{

		{
			name: "Тест 1 - Успешный",
			fields: fields{
				method: http.MethodGet, URL: "http://localhost/",
			},
			expected: expected{
				code: http.StatusOK,
			},
			PingOK: true,
		},
		{
			name: "Тест 2 - Ошибка в сервисном слою",
			fields: fields{
				method: http.MethodGet, URL: "http://localhost",
			},
			expected: expected{
				code: http.StatusInternalServerError,
				body: "\n",
			},
			PingOK: false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var (
				ctrl        = gomock.NewController(t)
				mockService = NewMockShortenerServicer(ctrl)

				h, _ = NewAPIHandler(mockService)

				bodyRequest  = strings.NewReader("")
				response     = httptest.NewRecorder()
				mockReq, err = http.NewRequest(tt.fields.method, tt.fields.URL, bodyRequest)
			)
			require.NoError(t, err)

			mockService.EXPECT().
				PingStorage().
				Return(tt.PingOK)

			h.PingDB(response, mockReq)

			assert.Equal(t, tt.expected.code, response.Code, "Код ответа не совпадает с ожидаемым.")
			assert.Equal(t, tt.expected.body, response.Body.String(), "Тело ответа не совпадает с ожидаемым.")
		})
	}
}

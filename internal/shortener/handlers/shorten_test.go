package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
)

func TestAPIHandler_JSONGenerateShortkeyHandler(t *testing.T) {

	testCases := []struct {
		name     string
		fields   fields
		expected expected
	}{

		{
			name: "Тест 1 - Успешный",
			fields: fields{
				method: http.MethodPost,
				URL:    "http://localhost",
				key:    "http://e.com",
				userID: uuid.New(),
			},
			expected: expected{
				code: http.StatusCreated, body: "http://localhost/testKey",
			},
		},
		{
			name: "Тест 2 - Пустой URL",
			fields: fields{
				method: http.MethodPost,
				URL:    "http://localhost",
				key:    "",
				userID: uuid.New(),
			},
			expected: expected{
				code: http.StatusBadRequest, body: "Bad Request\n",
			},
		},

		{
			name: "Тест 3 - В контексте нет UserID",
			fields: fields{
				method: http.MethodPost, URL: "http://localhost", key: "http://e.com", userID: uuid.Nil,
			},
			expected: expected{
				code: http.StatusInternalServerError, body: "\n",
			},
		},

		{
			name: "Тест 4 - GenerateShortURL возвращает ошибку InternalServicesError",
			fields: fields{
				method: http.MethodPost, URL: "http://localhost", key: "http://e.com", userID: uuid.New(),
			},
			expected: expected{
				code:       lErrors.InternalServicesError.HTTPCode,
				body:       lErrors.InternalServicesError.Error() + "\n",
				serviceErr: lErrors.InternalServicesError,
			},
		},

		{
			name: "Тест 5 - GenerateShortURL возвращает ошибку ErrorDuplicateFullURLServicesError",
			fields: fields{
				method: http.MethodPost, URL: "http://localhost", key: "http://e.com", userID: uuid.New(),
			},
			expected: expected{
				code:       lErrors.ErrorDuplicateFullURLServicesError.HTTPCode,
				body:       "http://localhost/testKey",
				serviceErr: lErrors.ErrorDuplicateFullURLServicesError,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			var (
				ctrl        = gomock.NewController(t)
				mockService = NewMockShortenerServicer(ctrl)

				h, _ = NewAPIHandler(mockService)

				bodyRequest  = strings.NewReader(fmt.Sprintf(`{"url" : "%s"}`, tt.fields.key))
				response     = httptest.NewRecorder()
				mockReq, err = http.NewRequest(tt.fields.method, tt.fields.URL, bodyRequest)
			)

			require.NoError(t, err)

			if tt.fields.userID != uuid.Nil {
				ctx := auth.WithUserID(mockReq.Context(), tt.fields.userID)
				mockReq = mockReq.WithContext(ctx)
			}

			if tt.fields.key != "" && tt.fields.userID != uuid.Nil {
				mockService.EXPECT().
					GenerateShortURL(mockReq.Context(), tt.fields.userID, tt.fields.key, mockReq.Host).
					Return(tt.expected.body, tt.expected.serviceErr)
			}

			h.JSONGenerateShortkeyHandler(response, mockReq)

			assert.Equal(t, tt.expected.code, response.Code, "Код ответа не совпадает с ожидаемым.")
			if tt.expected.code == http.StatusCreated || tt.expected.code == http.StatusConflict {
				assert.JSONEq(t, fmt.Sprintf(`{"result" : "%s"}`, tt.expected.body), response.Body.String(), "Тело ответа не совпадает с ожидаемым.")
			} else {
				assert.Equal(t, tt.expected.body, response.Body.String(), "Тело ответа не совпадает с ожидаемым.")

			}
		})
	}

}

func TestAPIHandler_JSONMultyGenerateShortkeyHandler(t *testing.T) {

	testCases := []struct {
		name     string
		fields   fields
		expected expected
		data     []models.JSONShortURL
	}{

		{
			name: "Тест 1 - Успешный",
			fields: fields{
				method: http.MethodPost,
				URL:    "http://localhost",
				key:    `[{"correlation_id":"123","original_url":"http://e.com","short_url":"http://short.ru"}]`,
				userID: uuid.New(),
			},
			expected: expected{
				code: http.StatusCreated, body: "http://localhost/testKey",
			},
			data: []models.JSONShortURL{{CorrelationID: "123", OriginalURL: "http://e.com", ShortURL: "http://short.ru"}},
		},
		{
			name: "Тест 2 - не верное тело запроса",
			fields: fields{
				method: http.MethodPost,
				URL:    "http://localhost",
				key:    `[{"correlation_id":123,"original_url":"","short_url":"http://short.ru"}]`,
				userID: uuid.New(),
			},
			expected: expected{
				code: http.StatusBadRequest, body: "Bad Request\n",
			},
		},

		{
			name: "Тест 3 - В контексте нет UserID",
			fields: fields{
				method: http.MethodPost,
				URL:    "http://localhost",
				key:    `[{"correlation_id":"123","original_url":"http://e.com","short_url":"http://short.ru"}]`,
			},
			expected: expected{
				code: http.StatusInternalServerError, body: "\n",
			},
		},

		{
			name: "Тест 4 - GenerateShortURL возвращает ошибку InternalServicesError",
			fields: fields{
				method: http.MethodPost,
				URL:    "http://localhost",
				key:    `[{"correlation_id":"123","original_url":"http://e.com","short_url":"http://short.ru"}]`,
				userID: uuid.New(),
			},
			expected: expected{
				code:       http.StatusInternalServerError,
				body:       "\n",
				serviceErr: lErrors.InternalServicesError,
			},
			data: []models.JSONShortURL{{CorrelationID: "123", OriginalURL: "http://e.com", ShortURL: "http://short.ru"}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			var (
				ctrl        = gomock.NewController(t)
				mockService = NewMockShortenerServicer(ctrl)

				h, _ = NewAPIHandler(mockService)

				bodyRequest  = strings.NewReader(tt.fields.key)
				response     = httptest.NewRecorder()
				mockReq, err = http.NewRequest(tt.fields.method, tt.fields.URL, bodyRequest)
			)

			require.NoError(t, err)

			if tt.fields.userID != uuid.Nil {
				ctx := auth.WithUserID(mockReq.Context(), tt.fields.userID)
				mockReq = mockReq.WithContext(ctx)
			}

			if tt.fields.key != "" && tt.fields.userID != uuid.Nil && tt.data != nil {
				mockService.EXPECT().
					GenerateMultyShortURL(mockReq.Context(), tt.fields.userID, tt.data, mockReq.Host).
					Return(tt.expected.serviceErr)

			}

			h.JSONMultyGenerateShortkeyHandler(response, mockReq)

			assert.Equal(t, tt.expected.code, response.Code, "Код ответа не совпадает с ожидаемым.")
			if tt.expected.code == http.StatusCreated || tt.expected.code == http.StatusConflict {
				assert.JSONEq(t, tt.fields.key, response.Body.String(), "Тело ответа не совпадает с ожидаемым.")
			} else {
				assert.Equal(t, tt.expected.body, response.Body.String(), "Тело ответа не совпадает с ожидаемым.")

			}
		})
	}

}

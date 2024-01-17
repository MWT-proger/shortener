package handlers

import (
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

func TestAPIHandler_GetListUserURLsHandler(t *testing.T) {

	testCases := []struct {
		name     string
		fields   fields
		expected expected
		listUrls []*models.JSONShortURL
	}{

		{
			name: "Тест 1 - Успешный",
			fields: fields{
				method: http.MethodGet, URL: "http://localhost", userID: uuid.New(),
			},
			expected: expected{
				code: http.StatusOK, body: `
											[
												{
													"original_url": "https://github.com/MaximMNsk/go-url-shortener2",
													"short_url": "http://localhost:7000/hEwcj"
												}
											]
											`,
			},
			listUrls: []*models.JSONShortURL{
				{
					OriginalURL: "https://github.com/MaximMNsk/go-url-shortener2",
					ShortURL:    "http://localhost:7000/hEwcj",
				},
			},
		},
		{
			name: "Тест 2 - Список пользователя пуст",
			fields: fields{
				method: http.MethodGet, URL: "http://localhost", userID: uuid.New(),
			},
			expected: expected{
				code:       lErrors.NoContentUserServicesError.HTTPCode,
				body:       lErrors.NoContentUserServicesError.Error() + "\n",
				serviceErr: lErrors.NoContentUserServicesError,
			},
		},

		{
			name: "Тест 3 - не указан UserID в context",
			fields: fields{
				method: http.MethodGet, URL: "http://localhost",
			},
			expected: expected{
				code: http.StatusInternalServerError,
				body: "\n",
			},
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

			if tt.fields.userID != uuid.Nil {
				ctx := auth.WithUserID(mockReq.Context(), tt.fields.userID)
				mockReq = mockReq.WithContext(ctx)

				mockService.EXPECT().
					GetListUserURLs(mockReq.Context(), tt.fields.userID, mockReq.Host).
					Return(tt.listUrls, tt.expected.serviceErr)
			}

			h.GetListUserURLsHandler(response, mockReq)

			assert.Equal(t, tt.expected.code, response.Code, "Код ответа не совпадает с ожидаемым.")

			if tt.listUrls != nil {
				assert.JSONEq(t, tt.expected.body, response.Body.String(), "Тело ответа не совпадает с ожидаемым.")
			} else {
				assert.Equal(t, tt.expected.body, response.Body.String(), "Тело ответа не совпадает с ожидаемым.")
			}
		})
	}

}

func TestAPIHandler_DeleteListUserURLsHandler(t *testing.T) {

	testCases := []struct {
		name     string
		fields   fields
		expected expected
		listKeys []string
	}{

		{
			name: "Тест 1 - Успешный",
			fields: fields{
				method: http.MethodDelete,
				URL:    "http://localhost",
				key:    `["LnJVM", "LnJVM2"]`,
				userID: uuid.New(),
			},
			expected: expected{
				code: http.StatusAccepted,
			},
			listKeys: []string{"LnJVM", "LnJVM2"},
		},
		{
			name: "Тест 2 - не указан UserID в context",
			fields: fields{
				method: http.MethodDelete,
				URL:    "http://localhost",
				key:    `["LnJVM", "LnJVM2"]`,
			},
			expected: expected{
				code: http.StatusInternalServerError,
				body: "\n",
			},
		},
		{
			name: "Тест 3 - не верное тело запроса",
			fields: fields{
				method: http.MethodDelete,
				URL:    "http://localhost",
				key:    "",
				userID: uuid.New(),
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: "Bad Request\n",
			},
			listKeys: []string{},
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

			if tt.fields.key != "" && tt.fields.userID != uuid.Nil {
				mockService.EXPECT().
					DeleteListUserURLs(mockReq.Context(), tt.fields.userID, tt.listKeys)
			}

			h.DeleteListUserURLsHandler(response, mockReq)

			assert.Equal(t, tt.expected.code, response.Code, "Код ответа не совпадает с ожидаемым.")
			assert.Equal(t, tt.expected.body, response.Body.String(), "Тело ответа не совпадает с ожидаемым.")

		})
	}

}

package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
)

type fields struct {
	method string
	URL    string
	key    string
	userID uuid.UUID
}
type expected struct {
	code       int
	body       string
	serviceErr error
}

func TestAPIHandler_GenerateShortkeyHandler(t *testing.T) {

	testCases := []struct {
		name     string
		fields   fields
		expected expected
	}{

		{
			name: "Тест 1 - Успешный",
			fields: fields{
				method: http.MethodPost, URL: "http://localhost", key: "http://e.com", userID: uuid.New(),
			},
			expected: expected{
				code: http.StatusCreated, body: "http://localhost/testKey",
			},
		},
		{
			name: "Тест 2 - В теле запроса нет значения",
			fields: fields{
				method: http.MethodPost, URL: "http://localhost", key: "", userID: uuid.New(),
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
			name: "GenerateShortURL возвращает ошибку InternalServicesError",
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
			name: "GenerateShortURL возвращает ошибку ErrorDuplicateFullURLServicesError",
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
					GenerateShortURL(mockReq.Context(), tt.fields.userID, tt.fields.key, mockReq.Host).
					Return(tt.expected.body, tt.expected.serviceErr)
			}

			h.GenerateShortkeyHandler(response, mockReq)

			assert.Equal(t, tt.expected.code, response.Code, "Код ответа не совпадает с ожидаемым.")
			assert.Equal(t, tt.expected.body, response.Body.String(), "Тело ответа не совпадает с ожидаемым.")

		})
	}

}

func TestAPIHandler_GetURLByKeyHandler(t *testing.T) {

	testCases := []struct {
		name     string
		fields   fields
		expected expected
	}{

		{
			name: "Тест 1 - Успешный",
			fields: fields{
				method: http.MethodGet, URL: "http://localhost/{shortKey}", key: "qwert",
			},
			expected: expected{
				code: http.StatusTemporaryRedirect, body: "https://e.ru",
			},
		},
		{
			name: "Тест 2 - Ошибка в сервисном слою",
			fields: fields{
				method: http.MethodGet, URL: "http://localhost/{shortKey}", key: "qwert",
			},
			expected: expected{
				code:       lErrors.GetFullURLServicesError.HTTPCode,
				body:       lErrors.GetFullURLServicesError.Error() + "\n",
				serviceErr: lErrors.GetFullURLServicesError,
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
				rctx         = chi.NewRouteContext()
			)
			require.NoError(t, err)

			rctx.URLParams.Add("shortKey", tt.fields.key)
			mockReq = mockReq.WithContext(context.WithValue(mockReq.Context(), chi.RouteCtxKey, rctx))

			mockService.EXPECT().
				GetFullURLByShortKey(mockReq.Context(), tt.fields.key).
				Return(tt.expected.body, tt.expected.serviceErr)

			h.GetURLByKeyHandler(response, mockReq)

			assert.Equal(t, tt.expected.code, response.Code, "Код ответа не совпадает с ожидаемым.")
			if tt.expected.serviceErr == nil {
				assert.Equal(t, tt.expected.body, response.Header().Get("Location"), "Тело ответа не совпадает с ожидаемым.")
			}
		})
	}
}

func BenchmarkAPIHandler_GetURLByKeyHandler(b *testing.B) {

	type testCase struct {
		name     string
		fields   fields
		expected expected
	}
	tt := testCase{

		name: "Замер получения URL по ключу",
		fields: fields{
			method: http.MethodGet, URL: "http://localhost/{shortKey}", key: "qwert",
		},
		expected: expected{
			code: http.StatusTemporaryRedirect, body: "https://e.ru",
		},
	}

	b.Run(tt.name, func(b *testing.B) {

		var (
			ctrl        = gomock.NewController(b)
			mockService = NewMockShortenerServicer(ctrl)

			h, _ = NewAPIHandler(mockService)

			bodyRequest  = strings.NewReader("")
			response     = httptest.NewRecorder()
			mockReq, err = http.NewRequest(tt.fields.method, tt.fields.URL, bodyRequest)
			rctx         = chi.NewRouteContext()
		)
		require.NoError(b, err)

		rctx.URLParams.Add("shortKey", tt.fields.key)
		mockReq = mockReq.WithContext(context.WithValue(mockReq.Context(), chi.RouteCtxKey, rctx))

		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			b.StopTimer()
			mockService.EXPECT().
				GetFullURLByShortKey(mockReq.Context(), tt.fields.key).
				Return(tt.expected.body, nil)
			b.StartTimer()

			h.GetURLByKeyHandler(response, mockReq)
		}

	})
}

func BenchmarkAPIHandler_GenerateShortkeyHandler(b *testing.B) {

	type testCase struct {
		name     string
		fields   fields
		expected expected
	}
	tt := testCase{

		name: "Замер генерации ShortKey",
		fields: fields{
			method: http.MethodPost, URL: "http://localhost", key: "http://e.com", userID: uuid.New(),
		},
		expected: expected{
			code: http.StatusCreated, body: "http://localhost/testKey",
		},
	}

	b.Run(tt.name, func(b *testing.B) {

		var (
			ctrl        = gomock.NewController(b)
			mockService = NewMockShortenerServicer(ctrl)

			h, _ = NewAPIHandler(mockService)

			bodyRequest  = strings.NewReader(tt.fields.key)
			response     = httptest.NewRecorder()
			mockReq, err = http.NewRequest(tt.fields.method, tt.fields.URL, bodyRequest)
		)

		require.NoError(b, err)

		if tt.fields.userID != uuid.Nil {
			ctx := auth.WithUserID(mockReq.Context(), tt.fields.userID)
			mockReq = mockReq.WithContext(ctx)
		}
		mockService.EXPECT().
			GenerateShortURL(mockReq.Context(), tt.fields.userID, tt.fields.key, mockReq.Host).
			Return(tt.expected.body, tt.expected.serviceErr)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.GenerateShortkeyHandler(response, mockReq)
		}

	})
}

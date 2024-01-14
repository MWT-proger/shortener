package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MWT-proger/shortener/configs"
	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestShortenerService_GenerateShortURL(t *testing.T) {

	type args struct {
		userID  uuid.UUID
		fullURL string
		rHost   string
	}
	type wantStorage struct {
		key string
		err error
	}
	testsCases := []struct {
		name        string
		args        args
		want        string
		wantStorage wantStorage
	}{
		{
			name:        "Тест 1 - Успешный",
			args:        args{fullURL: "http://full.ru", userID: uuid.UUID{}, rHost: "localhost"},
			wantStorage: wantStorage{key: "qwert", err: nil},
			want:        "http://localhost/qwert",
		},
		{
			name:        "Тест 2 - Ошибка в storage: lErrors.ErrorDuplicateFullURLServicesError",
			args:        args{fullURL: "http://full.ru", userID: uuid.UUID{}, rHost: "localhost"},
			wantStorage: wantStorage{key: "qwert", err: lErrors.ErrorDuplicateFullURLServicesError},
			want:        "http://localhost/qwert",
		},
		{
			name:        "Тест 3 - Ошибка в storage: lErrors.InternalServicesError",
			args:        args{fullURL: "http://full.ru", userID: uuid.UUID{}, rHost: "localhost"},
			wantStorage: wantStorage{key: "qwert", err: lErrors.InternalServicesError},
			want:        "",
		},
	}
	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {

			var (
				ctx         = context.TODO()
				conf        = configs.Config{}
				ctrl        = gomock.NewController(t)
				mockStorage = NewMockStorager(ctrl)
				s           = NewShortenerService(ctx, conf, mockStorage)
			)
			mockStorage.EXPECT().
				Set(ctx, models.ShortURL{FullURL: tt.args.fullURL, UserID: tt.args.userID}).
				Return(tt.wantStorage.key, tt.wantStorage.err)

			got, err := s.GenerateShortURL(ctx, tt.args.userID, tt.args.fullURL, tt.args.rHost)

			assert.ErrorIs(t, err, tt.wantStorage.err, "ошибка не совпадает с запланированным")
			assert.Equal(t, tt.want, got, "ShortURL не совпадает с запланированным")

		})
	}
}

func TestShortenerService_GenerateMultyShortURL(t *testing.T) {

	type args struct {
		userID uuid.UUID
		data   []models.JSONShortURL
		rHost  string
	}
	type wantStorage struct {
		err error
	}
	testsCases := []struct {
		name        string
		args        args
		want        error
		wantStorage wantStorage
	}{
		{
			name:        "Тест 1 - Успешный",
			args:        args{data: []models.JSONShortURL{}, userID: uuid.UUID{}, rHost: "localhost"},
			wantStorage: wantStorage{err: nil},
			want:        nil,
		},
		{
			name:        "Тест 2 - Ошибка в хранилище",
			args:        args{data: []models.JSONShortURL{}, userID: uuid.UUID{}, rHost: "localhost"},
			wantStorage: wantStorage{err: errors.New("")},
			want:        lErrors.InternalServicesError,
		},
	}
	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {

			var (
				ctx         = context.TODO()
				conf        = configs.Config{}
				ctrl        = gomock.NewController(t)
				mockStorage = NewMockStorager(ctrl)
				s           = NewShortenerService(ctx, conf, mockStorage)
			)

			mockStorage.EXPECT().
				SetMany(ctx, tt.args.data, "http://"+tt.args.rHost+"/", tt.args.userID).
				Return(tt.wantStorage.err)

			err := s.GenerateMultyShortURL(ctx, tt.args.userID, tt.args.data, tt.args.rHost)

			assert.ErrorIs(t, err, tt.want, "ошибка не совпадает с запланированным")

		})
	}
}

func TestShortenerService_GetFullURLByShortKey(t *testing.T) {

	type args struct {
		shortKey string
	}
	type wantStorage struct {
		data models.ShortURL
		err  error
	}
	type want struct {
		data string
		err  error
	}
	testsCases := []struct {
		name        string
		args        args
		want        want
		wantStorage wantStorage
	}{
		{
			name:        "Тест 1 - Успешный",
			args:        args{shortKey: "qwert"},
			wantStorage: wantStorage{data: models.ShortURL{FullURL: "http://f.ru"}, err: nil},
			want:        want{data: "http://f.ru", err: nil},
		},
		{
			name:        "Тест 2 - Storage возвращает ошибку",
			args:        args{shortKey: "qwert"},
			wantStorage: wantStorage{err: errors.New("")},
			want:        want{data: "", err: lErrors.GetFullURLServicesError},
		},
		{
			name:        "Тест 3 - Storage возвращает data c FullURL == ''",
			args:        args{shortKey: "qwert"},
			wantStorage: wantStorage{data: models.ShortURL{FullURL: ""}, err: nil},
			want:        want{data: "", err: lErrors.NotFoundServicesError},
		},
		{
			name:        "Тест 4 - Storage возвращает data c DeletedFlag == true",
			args:        args{shortKey: "qwert"},
			wantStorage: wantStorage{data: models.ShortURL{FullURL: "http://f.ru", DeletedFlag: true}, err: nil},
			want:        want{data: "", err: lErrors.GoneServicesError},
		},
	}

	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {

			var (
				ctx         = context.TODO()
				conf        = configs.Config{}
				ctrl        = gomock.NewController(t)
				mockStorage = NewMockStorager(ctrl)
				s           = NewShortenerService(ctx, conf, mockStorage)
			)

			mockStorage.EXPECT().
				Get(ctx, tt.args.shortKey).
				Return(tt.wantStorage.data, tt.wantStorage.err)

			got, err := s.GetFullURLByShortKey(ctx, tt.args.shortKey)

			assert.ErrorIs(t, err, tt.want.err, "ошибка не совпадает с запланированным")
			assert.Equal(t, tt.want.data, got, "FullURL не совпадает с запланированным")

		})
	}
}

func TestShortenerService_GetListUserURLs(t *testing.T) {

	type args struct {
		userID uuid.UUID
		rHost  string
	}
	type wantStorage struct {
		data []*models.JSONShortURL
		err  error
	}
	type want struct {
		data []*models.JSONShortURL
		err  error
	}
	testsCases := []struct {
		name        string
		args        args
		want        want
		wantStorage wantStorage
	}{
		{
			name:        "Тест 1 - Успешный",
			args:        args{userID: uuid.UUID{}, rHost: "localhost"},
			wantStorage: wantStorage{data: []*models.JSONShortURL{{ShortURL: "qwerty"}}, err: nil},
			want:        want{data: []*models.JSONShortURL{{ShortURL: "http://localhost/qwerty"}}, err: nil},
		},
		{
			name:        "Тест 2 - Storage возвращает ошибку",
			args:        args{userID: uuid.UUID{}, rHost: "localhost"},
			wantStorage: wantStorage{data: []*models.JSONShortURL{{ShortURL: "qwerty"}}, err: errors.New("")},
			want:        want{err: lErrors.GetFullURLServicesError},
		},
		{
			name:        "Тест 3 - Storage возвращает nil",
			args:        args{userID: uuid.UUID{}, rHost: "localhost"},
			wantStorage: wantStorage{},
			want:        want{err: lErrors.NoContentUserServicesError},
		},
	}
	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {

			var (
				ctx         = context.TODO()
				conf        = configs.Config{}
				ctrl        = gomock.NewController(t)
				mockStorage = NewMockStorager(ctrl)
				s           = NewShortenerService(ctx, conf, mockStorage)
			)
			mockStorage.EXPECT().
				GetList(ctx, tt.args.userID).
				Return(tt.wantStorage.data, tt.wantStorage.err)

			got, err := s.GetListUserURLs(ctx, tt.args.userID, tt.args.rHost)

			assert.ErrorIs(t, err, tt.want.err, "ошибка не совпадает с запланированным")
			assert.Equal(t, tt.want.data, got, "listData не совпадает с запланированным")

		})
	}
}

func TestShortenerService_DeleteListUserURLs(t *testing.T) {

	type args struct {
		userID uuid.UUID
		data   []string
	}

	type want struct {
		data []models.DeletedShortURL
	}
	testsCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Тест 1 - Успешный",
			args: args{userID: uuid.Nil, data: []string{"qwert", "asdf"}},
			want: want{data: []models.DeletedShortURL{{Payload: "qwert"}, {Payload: "asdf"}}},
		},
	}
	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {

			var (
				data []models.DeletedShortURL

				ctx  = context.TODO()
				conf = configs.Config{}

				ctrl        = gomock.NewController(t)
				mockStorage = NewMockStorager(ctrl)
				s           = NewShortenerService(ctx, conf, mockStorage)
			)

			s.DeleteListUserURLs(ctx, tt.args.userID, tt.args.data)
			ticker := time.NewTicker(6 * time.Second)

		L:
			for {
				select {
				case d := <-s.deletedChan:
					data = append(data, d)

				case <-ticker.C:
					break L
				}
				if len(data) == len(tt.want.data) {
					break L
				}

			}
			assert.Equal(t, tt.want.data, data, "listData не совпадает с запланированным")

		})
	}
}

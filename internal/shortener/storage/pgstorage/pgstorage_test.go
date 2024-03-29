package pgstorage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
)

func TestPgStorageGet(t *testing.T) {
	testCases := []struct {
		name     string
		shortKey string
		result   models.ShortURL
		errors   error
	}{
		{
			name:     "Тест 1 - Проверяем на успех",
			shortKey: "testkey",
			result:   models.ShortURL{FullURL: "http://example.ru"},
			errors:   nil,
		},

		{
			name:     "Тест 2 - Проверяем на успех",
			shortKey: "testkey",
			result:   models.ShortURL{FullURL: ""},
			errors:   sql.ErrNoRows,
		},
	}
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	s := &pgStorage{
		db: sqlxDB,
	}
	querySQL := "SELECT full_url, is_deleted FROM content.shorturl WHERE short_key = $1 LIMIT 1;"

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			if tt.errors != nil {
				mock.ExpectQuery(querySQL).
					WithArgs(tt.shortKey).
					WillReturnError(tt.errors)
			} else {
				rows := sqlmock.NewRows([]string{"full_url", "is_deleted"}).AddRow(tt.result.FullURL, tt.result.DeletedFlag)

				mock.ExpectQuery(querySQL).
					WithArgs(tt.shortKey).
					WillReturnRows(rows)
			}
			got, _ := s.Get(context.Background(), tt.shortKey)

			assert.Equal(t, got.FullURL, tt.result.FullURL, "Результат не совпадает с ожиданием")
		})
	}
}

func TestPgStorageDoSet(t *testing.T) {
	testCases := []struct {
		name        string
		model       models.ShortURL
		errorsDB    error
		errorString string
	}{

		{
			name:        "Тест 1 - Проверяем на дубликат short_key",
			model:       models.ShortURL{ShortKey: "testKey", FullURL: "http://example.ru", UserID: uuid.New()},
			errorsDB:    &pgconn.PgError{Code: "23505", ConstraintName: "shorturl_short_key_key"},
			errorString: lErrors.ErrorDuplicateShortKey.Error(),
		},
		{
			name:        "Тест 2 - Проверяем на успех",
			model:       models.ShortURL{ShortKey: "testKey", FullURL: "http://example.ru", UserID: uuid.New()},
			errorsDB:    nil,
			errorString: "",
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	s := &pgStorage{
		db: sqlxDB,
	}
	querySQL := "INSERT INTO content.shorturl (short_key, full_url, user_id) VALUES($1,$2,$3)"

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			if tt.errorsDB != nil {
				mock.ExpectBegin()
				mock.ExpectPrepare(querySQL).ExpectExec().
					WithArgs(tt.model.ShortKey, tt.model.FullURL, tt.model.UserID).
					WillReturnError(tt.errorsDB)
			} else {
				mock.ExpectBegin()
				mock.ExpectPrepare(querySQL).ExpectExec().
					WithArgs(tt.model.ShortKey, tt.model.FullURL, tt.model.UserID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}
			ctx := context.TODO()
			ctx = auth.WithUserID(ctx, tt.model.UserID)
			err := s.doSet(ctx, &tt.model)

			if tt.errorString != "" {
				assert.EqualError(t, err, tt.errorString, "Ошибка не совпадает")
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

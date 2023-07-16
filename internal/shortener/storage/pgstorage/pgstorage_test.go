package pgstorage

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestPgStorageGet(t *testing.T) {
	testCases := []struct {
		name     string
		shortKey string
		result   string
		errors   error
	}{
		{
			name:     "Тест 1 - Проверяем на успех",
			shortKey: "testkey",
			result:   "http://example.ru",
			errors:   nil,
		},

		{
			name:     "Тест 2 - Проверяем на успех",
			shortKey: "testkey",
			result:   "",
			errors:   sql.ErrNoRows,
		},
	}
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	s := &PgStorage{
		db: db,
	}
	querySQL := "SELECT full_url FROM content.shorturl WHERE short_key = $1 LIMIT 1;"

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			if tt.errors != nil {
				mock.ExpectQuery(querySQL).
					WithArgs(tt.shortKey).
					WillReturnError(tt.errors)
			} else {
				rows := sqlmock.NewRows([]string{"full_url"}).AddRow(tt.result)

				mock.ExpectQuery(querySQL).
					WithArgs(tt.shortKey).
					WillReturnRows(rows)
			}
			got, _ := s.Get(tt.shortKey)

			assert.Equal(t, got, tt.result, "Результат не совпадает с ожиданием")
		})
	}
}

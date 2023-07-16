package pgstorage

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type PgStorage struct {
	storage.Storage
	db *sql.DB
}

func (s *PgStorage) Init(ctx context.Context) error {
	conf := configs.GetConfig()

	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		return err
	}
	s.db = db

	if err := s.Migration(); err != nil {
		return err
	}

	return nil

}

// Migration() проверяет нувые миграции и при неообходимости добавляет в БД
func (s *PgStorage) Migration() error {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(s.db, "migrations"); err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

// Добавляет в хранилище полную ссылку и присваевает ей ключ
func (s *PgStorage) Set(fullURL string) (string, error) {

	return "", nil

}

// Достаёт из хранилища и возвращает полную ссылку по ключу
func (s *PgStorage) Get(shortURL string) (string, error) {
	var FullURL string

	row := s.db.QueryRowContext(context.Background(),
		"SELECT full_url "+
			"FROM content.shorturl WHERE short_key = $1 LIMIT 1;", shortURL)

	err := row.Scan(&FullURL)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		logger.Log.Error(err.Error())
		return "", err
	}

	return FullURL, nil
}

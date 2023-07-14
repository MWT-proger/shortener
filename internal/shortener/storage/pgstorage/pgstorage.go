package pgstorage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

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
	return "", nil
}

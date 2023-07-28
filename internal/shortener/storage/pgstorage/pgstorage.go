package pgstorage

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/MWT-proger/shortener/configs"
	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
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
func (s *PgStorage) Set(newModel models.ShortURL) (string, error) {

	ctx := context.Background()

	for {
		newModel.ShortKey = utils.StringWithCharset(5)

		if err := s.doSet(ctx, &newModel); err != nil {

			if errors.Is(err, &lErrors.ErrorDuplicateShortKey{}) {
				continue
			}
			if errors.Is(err, &lErrors.ErrorDuplicateFullURL{}) {
				newModel.ShortKey, _ = s.getShortKey(newModel.FullURL)
				return newModel.ShortKey, &lErrors.ErrorDuplicateFullURL{}
			}
			return "", err

		}
		break

	}
	return newModel.ShortKey, nil

}

// Добавляет в хранилище полную ссылку и присваевает ей ключ
func (s *PgStorage) SetMany(data []models.JSONShortURL, baseShortURL string, userID uuid.UUID) error {

	ctx := context.Background()

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO content.shorturl (short_key, full_url, user_id) VALUES($1,$2,$3) ON CONFLICT (short_key) DO NOTHING RETURNING short_key")

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	defer stmt.Close()

	for i, v := range data {

		shortKey := ""
		for {
			row := stmt.QueryRowContext(ctx, utils.StringWithCharset(5), v.OriginalURL, userID)

			err := row.Scan(&shortKey)

			if err != nil {

				if errors.Is(err, sql.ErrNoRows) {
					logger.Log.Debug("Возвращается пустой ключ из за дублирования в БД")
					continue
				}

				logger.Log.Error(err.Error())
				return err
			}
			break
		}

		data[i].ShortURL = baseShortURL + shortKey
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

// doSet() Добавляет в БД или возвращает ошибку
func (s *PgStorage) doSet(ctx context.Context, model *models.ShortURL) error {

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO content.shorturl (short_key, full_url, user_id) VALUES($1,$2,$3)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, model.ShortKey, model.FullURL, model.UserID)

	if err != nil {
		logger.Log.Error(err.Error())
		var pgError *pgconn.PgError
		if errors.As(err, &pgError); errors.Is(err, pgError) {

			if pgError.Code == "23505" && pgError.ConstraintName == "shorturl_short_key_key" {
				return &lErrors.ErrorDuplicateShortKey{}
			}
			if pgError.Code == "23505" && pgError.ConstraintName == "shorturl_full_url_key" {
				return &lErrors.ErrorDuplicateFullURL{}
			}
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

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

func (s *PgStorage) getShortKey(FullURL string) (string, error) {
	var shortURL string

	row := s.db.QueryRowContext(context.Background(),
		"SELECT short_key "+
			"FROM content.shorturl WHERE full_url = $1 LIMIT 1;", FullURL)

	err := row.Scan(&shortURL)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		logger.Log.Error(err.Error())
		return "", err
	}

	return shortURL, nil
}

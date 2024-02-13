package pgstorage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	"github.com/MWT-proger/shortener/configs"
	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// pgStorage хранит данные в Postgres.
type pgStorage struct {
	db *sqlx.DB
}

// NewPgStorage создаёт и возвращает новый экземпляр pgStorage.
func NewPgStorage(ctx context.Context, conf configs.Config) (*pgStorage, error) {

	var (
		s       = &pgStorage{}
		db, err = sqlx.Open("pgx", conf.DatabaseDSN)
	)

	if err != nil {
		return nil, err
	}

	s.db = db

	if err := s.migration(); err != nil {
		return nil, err
	}

	return s, nil
}

// Close Закрывает соединение.
func (s *pgStorage) Close() error {

	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

// Set Добавляет в хранилище полную ссылку и присваевает ей ключ.
func (s *pgStorage) Set(ctx context.Context, newModel models.ShortURL) (string, error) {

	for {
		newModel.ShortKey = utils.StringWithCharset(5)

		if err := s.doSet(ctx, &newModel); err != nil {

			if errors.Is(err, lErrors.ErrorDuplicateShortKey) {
				continue
			}
			if errors.Is(err, lErrors.ErrorDuplicateFullURLServicesError) {
				newModel.ShortKey, _ = s.getShortKey(ctx, newModel.FullURL)
				return newModel.ShortKey, lErrors.ErrorDuplicateFullURLServicesError
			}
			return "", err

		}
		break

	}
	return newModel.ShortKey, nil

}

// SetMany Добавляет в хранилище несколько полных ссылок и присваевает им ключи.
func (s *pgStorage) SetMany(ctx context.Context, data []models.JSONShortURL, baseShortURL string, userID uuid.UUID) error {
	var (
		pgError *pgconn.PgError
		tx, err = s.db.BeginTx(ctx, nil)
	)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO content.shorturl (short_key, full_url, user_id) VALUES($1,$2,$3) ON CONFLICT (short_key) DO NOTHING RETURNING short_key")

	if err != nil {
		logger.Log.Error("Ошибка создания предварительного запроса", logger.ErrorField(err))
		return err
	}
	defer stmt.Close()

	for i, v := range data {
		nameSavePoint := fmt.Sprintf("multy_savepoint_%s", strconv.Itoa(i))
		tx.ExecContext(ctx, "SAVEPOINT "+nameSavePoint)
		FullURLIsExist := false
		shortKey := ""

		for {

			row := stmt.QueryRowContext(ctx, utils.StringWithCharset(5), v.OriginalURL, userID)

			err := row.Scan(&shortKey)

			if err != nil {

				if errors.Is(err, sql.ErrNoRows) {
					logger.Log.Debug("Возвращается пустой ключ из за дублирования в БД")
					continue
				}

				if errors.As(err, &pgError); errors.Is(err, pgError) {
					if pgError.Code == "23505" && pgError.ConstraintName == "shorturl_full_url_key" {
						tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT "+nameSavePoint)
						logger.Log.Debug("FullURL is exist", logger.ErrorField(err))
						FullURLIsExist = true
						break
					}
				}

				logger.Log.Error("Ошибка добавления новой строки", logger.ErrorField(err))
				return err
			}

			break
		}
		if FullURLIsExist {

			row := tx.QueryRowContext(ctx,
				"SELECT short_key "+
					"FROM content.shorturl WHERE full_url = $1 LIMIT 1;", v.OriginalURL)

			err := row.Scan(&shortKey)

			if err != nil {
				return err
			}
		}
		data[i].ShortURL = baseShortURL + shortKey
		tx.ExecContext(ctx, "RELEASE SAVEPOINT "+nameSavePoint)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

// Get Достаёт из хранилища и возвращает полную ссылку по ключу.
func (s *pgStorage) Get(ctx context.Context, shortURL string) (models.ShortURL, error) {
	var model models.ShortURL

	row := s.db.QueryRowContext(ctx,
		"SELECT full_url, is_deleted "+
			"FROM content.shorturl WHERE short_key = $1 LIMIT 1;", shortURL)

	err := row.Scan(&model.FullURL, &model.DeletedFlag)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return model, nil
		}

		logger.Log.Error("ошибка", logger.ErrorField(err))
		return model, err
	}

	return model, nil
}

// GetList достает список url users.
func (s *pgStorage) GetList(ctx context.Context, userID uuid.UUID) ([]*models.JSONShortURL, error) {
	var (
		list      = []*models.JSONShortURL{}
		args      = map[string]interface{}{"user_id": userID}
		stmt, err = s.db.PrepareNamedContext(ctx, "SELECT * "+
			"FROM content.shorturl WHERE user_id = :user_id;")
	)

	if err != nil {
		logger.Log.Error("ошибка", logger.ErrorField(err))
		return nil, err
	}
	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &list, args); err != nil {
		logger.Log.Error("ошибка", logger.ErrorField(err))
		return nil, err
	}

	return list, nil
}

// getShortKey получает короткий url по полному.
func (s *pgStorage) getShortKey(ctx context.Context, FullURL string) (string, error) {
	var shortURL string

	row := s.db.QueryRowContext(ctx,
		"SELECT short_key "+
			"FROM content.shorturl WHERE full_url = $1 LIMIT 1;", FullURL)

	err := row.Scan(&shortURL)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		logger.Log.Error("ошибка", logger.ErrorField(err))
		return "", err
	}

	return shortURL, nil
}

// DeleteList удаляет список коротких ссылок пользователя.
func (s *pgStorage) DeleteList(ctx context.Context, data ...models.DeletedShortURL) error {
	var values []string
	var args = []any{true}

	for i, v := range data {
		base := i * 2
		params := fmt.Sprintf("(user_id = $%d AND short_key= $%d )", base+2, base+3)

		values = append(values, params)
		args = append(args, v.UserID, v.Payload)
	}

	query := `UPDATE content.shorturl SET is_deleted= $1 WHERE ` +
		strings.Join(values, " OR ") + `;`
	_, err := s.db.ExecContext(ctx, query, args...)

	if err != nil {
		logger.Log.Error("ошибка", logger.ErrorField(err))
		return err
	}

	return nil
}

// Ping Прверяет соединение.
func (s *pgStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

// migration() проверяет нувые миграции и при неообходимости добавляет в БД.
func (s *pgStorage) migration() error {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(s.db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}

// doSet() Добавляет в БД или возвращает ошибку.
func (s *pgStorage) doSet(ctx context.Context, model *models.ShortURL) error {

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
		logger.Log.Error("ошибка", logger.ErrorField(err))
		var pgError *pgconn.PgError
		if errors.As(err, &pgError); errors.Is(err, pgError) {

			if pgError.Code == "23505" && pgError.ConstraintName == "shorturl_short_key_key" {
				return lErrors.ErrorDuplicateShortKey
			}
			if pgError.Code == "23505" && pgError.ConstraintName == "shorturl_full_url_key" {
				return lErrors.ErrorDuplicateFullURLServicesError
			}
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

// CountUsersAndUrls возвращает количество пользователей и сокращенных URL в сервисе.
func (s *pgStorage) CountUsersAndUrls(ctx context.Context) (urls int, users int, err error) {

	row := s.db.QueryRowContext(ctx,
		"SELECT count(distinct user_id) as users, count(short_key) as urls FROM content.shorturl;")

	err = row.Scan(&users, &urls)

	if err != nil {

		logger.Log.Error("ошибка", logger.ErrorField(err))
		return 0, 0, err
	}

	return urls, users, err
}

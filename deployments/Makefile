db_up:
	docker compose --env-file .env  up
db_migration:
	@if [ -z "${name}" ]; then echo "Использование: make db_migration name=name_of_migration_file"; exit 1; fi
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir ../internal/shortener/storage/pgstorage/migrations create ${name} sql

# mockgen:
# 	go run go.uber.org/mock/mockgen@latest -destination=internal/shortener/storage/mocks/mock_pgstorage.go -package=mocks github.com/MWT-proger/shortener/internal/shortener/storage/pgstorage PgStorager

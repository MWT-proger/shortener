-- +goose Up
-- +goose StatementBegin
ALTER TABLE "content"."shorturl" ADD COLUMN "is_deleted" bool NOT NULL DEFAULT 'FALSE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "content"."shorturl" DROP COLUMN "is_deleted";
-- +goose StatementEnd

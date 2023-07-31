-- +goose Up
-- +goose StatementBegin
ALTER TABLE "content"."shorturl" ADD COLUMN "is_deleted" bool;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "content"."shorturl" DROP COLUMN "is_deleted";
-- +goose StatementEnd

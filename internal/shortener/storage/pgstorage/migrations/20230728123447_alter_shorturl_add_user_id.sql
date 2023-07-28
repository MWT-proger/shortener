-- +goose Up
-- +goose StatementBegin
ALTER TABLE "content"."shorturl" ADD COLUMN "user_id" uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "content"."shorturl" DROP COLUMN "user_id";
-- +goose StatementEnd

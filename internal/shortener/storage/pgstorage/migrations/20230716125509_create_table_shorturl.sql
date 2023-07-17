-- +goose Up
-- +goose StatementBegin
CREATE TABLE "content"."shorturl" (
    short_key VARCHAR(10) PRIMARY KEY,
    full_url TEXT NOT NULL
);


ALTER TABLE ONLY "content"."shorturl"
    ADD CONSTRAINT "shorturl_short_key_key" UNIQUE (short_key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."shorturl";
-- +goose StatementEnd

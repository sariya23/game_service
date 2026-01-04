-- +goose Up
-- +goose StatementBegin
alter table game rename column image_url to image_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table game rename column image_key to image_url;
-- +goose StatementEnd

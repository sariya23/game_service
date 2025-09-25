-- +goose Up
-- +goose StatementBegin
alter table game
add column published boolean not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table drop column published;
-- +goose StatementEnd

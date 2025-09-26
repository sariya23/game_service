-- +goose Up
-- +goose StatementBegin
alter table game
drop column published;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table game
add column published boolean not null default false;
-- +goose StatementEnd

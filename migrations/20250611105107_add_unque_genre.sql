-- +goose Up
-- +goose StatementBegin
alter table genre add constraint unique_genre unique (genre_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop constraint unique_genre;
-- +goose StatementEnd

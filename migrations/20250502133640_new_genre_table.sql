-- +goose Up
-- +goose StatementBegin
create table if not exists genre(
    genre_id smallint generated always as identity primary key,
    genre_name varchar(70) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists genre;
-- +goose StatementEnd

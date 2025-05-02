-- +goose Up
-- +goose StatementBegin
create table if not exists game(
    game_id bigint generated always as identity primary key,
    title varchar(160) not null check(char_length(title) > 1),
    description text default 'Описание не добавили :(',
    release_date date not null,
    image_url text 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete table if exists game;
-- +goose StatementEnd

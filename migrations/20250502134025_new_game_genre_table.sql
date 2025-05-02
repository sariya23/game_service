-- +goose Up
-- +goose StatementBegin
create table if not exists game_genre(
    game_id bigint not null references game(game_id),
    genre_id smallint not null references genre(genre_id),
    primary key (game_id, genre_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists game_genre;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE game_genre
ADD CONSTRAINT game_genre__genre_id_fkey
FOREIGN KEY (genre_id) REFERENCES genre(genre_id)
ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table game_genre drop constraint game_genre__genre_id_fkey;
-- +goose StatementEnd

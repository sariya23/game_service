-- +goose Up
-- +goose StatementBegin
ALTER TABLE game_tag
ADD CONSTRAINT game_tag__game_id_fkey
FOREIGN KEY (game_id) REFERENCES game(game_id)
ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table game_genre drop constraint game_tag__game_id_fkey;
-- +goose StatementEnd

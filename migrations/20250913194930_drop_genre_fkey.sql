-- +goose Up
-- +goose StatementBegin
alter table game_genre drop constraint game_genre_genre_id_fkey;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

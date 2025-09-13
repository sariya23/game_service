-- +goose Up
-- +goose StatementBegin
alter table game_tag drop constraint game_tag_game_id_fkey;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

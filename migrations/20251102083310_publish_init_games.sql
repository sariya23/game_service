-- +goose Up
-- +goose StatementBegin
update game set game_status_id = 2;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update game set game_status_id = 0;
-- +goose StatementEnd

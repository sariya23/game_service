-- +goose Up
-- +goose StatementBegin
insert into game_status values
(0, 'draft'),
(1, 'pending'),
(2, 'publish');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate game_status;
-- +goose StatementEnd

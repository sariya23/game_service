-- +goose Up
-- +goose StatementBegin
insert into game_tag (game_id, tag_id) values 
(1, 2), 
(1, 3), 
(1, 5), 
(1, 7), 
(1, 10), 
(2, 1),
(2, 9), 
(3, 2),
(3, 3), 
(3, 4),
(3, 5),
(4, 2), 
(4, 4), 
(4, 5), 
(5, 2), 
(5, 6),
(5, 9);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate game_tag;
-- +goose StatementEnd

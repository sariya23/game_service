-- +goose Up
-- +goose StatementBegin
insert into game_genre (game_id, genre_id) values 
(1, 2),
(1, 3), 
(2, 1),
(2, 9), 
(3, 3), 
(3, 4),
(4, 2),
(4, 10), 
(5, 1),
(5, 8);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate game_genre;
-- +goose StatementEnd

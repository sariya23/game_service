-- +goose Up
-- +goose StatementBegin
alter table game 
add column game_status_id smallint not null default 0 references game_status(game_status_id);

update game 
set game_status_id = case when published=false then 0
else 2
end;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table game
drop column game_status_id;
-- +goose StatementEnd

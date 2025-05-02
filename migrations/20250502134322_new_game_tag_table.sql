-- +goose Up
-- +goose StatementBegin
create table if not exists game_tag(
    game_id bigint not null references game(game_id),
    tag_id int not null references tag(tag_id),
    primary key (game_id, tag_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists game_tag;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
create table if not exists game_status (
    game_status_id smallint not null primary key,
    name varchar(10) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists game_status;
-- +goose StatementEnd

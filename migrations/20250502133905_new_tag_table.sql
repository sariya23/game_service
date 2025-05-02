-- +goose Up
-- +goose StatementBegin
create table if not exists tag(
    tag_id int generated always as identity primary key,
    tag_name varchar(70) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tag;
-- +goose StatementEnd

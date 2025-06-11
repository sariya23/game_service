-- +goose Up
-- +goose StatementBegin
alter table tag add constraint unique_tag unique (tag_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop constraint unique_tag;
-- +goose StatementEnd

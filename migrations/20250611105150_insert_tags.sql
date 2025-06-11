-- +goose Up
-- +goose StatementBegin
insert into tag (tag_name) values 
('Мультиплеер'),
('Одиночная игра'),
('Открытый мир'),
('Пиксельная графика'),
('Сюжетная'),
('Выживание'),
('Фэнтези'),
('Научная фантастика'),
('От первого лица'),
('От третьего лица');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate tag;
-- +goose StatementEnd

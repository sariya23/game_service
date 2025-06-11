-- +goose Up
-- +goose StatementBegin
insert into genre (genre_name) values 
('Экшен'),
('Приключения'),
('Ролевая игра'),
('Симулятор'),
('Стратегия'),
('Спорт'),
('Головоломка'),
('Хоррор'),
('Шутер'),
('Платформер');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate genre;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
insert into game (title, description, release_date) values 
(
    'The Witcher 3: Wild Hunt',
    'Эпическая ролевая игра с открытым миром, в которой вы играете за ведьмака Геральта в поисках приемной дочери.',
    '2015-05-19'::date
),
(
    'Counter-Strike 2',
    'Тактический шутер от первого лица, продолжение культовой серии игр от Valve.',
    '2023-09-27'::date
),
(
    'Stardew Valley',
    'Фермерский симулятор с элементами RPG, где игрок восстанавливает старую ферму и развивает отношения с жителями.',
    '2016-02-26'::date
),
(
    'Hollow Knight',
    'Атмосферный платформер с элементами метроидвании, исследующий подземный мир Халлоунест.',
    '2017-02-24'::date
),
(
    'Resident Evil Village',
    'Хоррор от первого лица, продолжающий историю Итана Уинтерса в таинственной деревне.',
    '2021-05-07'::date
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate game;
-- +goose StatementEnd

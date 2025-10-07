
ENV ?= local
ENV_FILE = ./config/$(ENV).env

include ${ENV_FILE}


# ЛОКАЛЬНЫЙ ЗАПУСК
# usage: Нужно указать префикс env файла. То есть
# если используем local.env, то пишем make migrate ENV=local.
.PHONY: migrate
migrate:
	goose -dir migrations postgres \
	"postgresql://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)\
	@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)\
	?sslmode=$(SSL_MODE)" up

.PHONY: run
run:
	go run cmd/main.go --config config/local.env

.PHONY: test
test:
	go test ./...

.PHONY: mock
mock:
	find . -name '*_mock.go' -delete
	mockgen -source internal/service/game/game.go \
	-destination=internal/service/game/mocks/game.go -package=mock_gameservice


# ДЛЯ ТЕСТОВ В ДОКЕРЕ
.PHONY: test_compose_up
test_compose_up:
	docker-compose -p test_game_service -f deployments/docker/test/docker-compose.yaml  \
	--env-file ./config/test.env up -d

.PHONY: test_migrate
test_migrate:
	goose -dir migrations postgres \
	"postgresql://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)\
	@$(POSTGRES_HOST_INNER_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)\
	?sslmode=$(SSL_MODE)" up

.PHONY: test_compose_down
test_compose_down:
	docker-compose -p test_game_service -f deployments/docker/test/docker-compose.yaml \
	--env-file ./config/test.env rm -fvs
	docker rmi test_game_service-app || true
	docker rmi test_game_service-migration || true
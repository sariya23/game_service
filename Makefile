ENV ?= local
ENV_FILE = ./config/$(ENV).env

# usage: Нужно указать префикс env файла. То есть
# если используем local.env, то пишем make migrate ENV=local.

include ${ENV_FILE}


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


# ДЛЯ ЛОКАЛЬНОГО ЗАПУСКА
.PHONY: service_compose_up
service_compose_up:
	docker-compose -p game_service -f deployments/docker/local/docker-compose.yaml  \
	--env-file ./config/local.env up -d

.PHONY: service_migrate_inner
service_migrate_inner:
	goose -dir migrations postgres \
	"postgresql://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)\
	@$(POSTGRES_HOST_INNER_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)\
	?sslmode=$(SSL_MODE)" up

.PHONY: service_migrate_outer
service_migrate_outer:
	goose -dir migrations postgres \
	"postgresql://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)\
	@$(POSTGRES_HOST_OUTER_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)\
	?sslmode=$(SSL_MODE)" up

.PHONY: service_compose_down
service_compose_down:
	docker-compose -p game_service -f deployments/docker/local/docker-compose.yaml \
	--env-file ./config/local.env rm -fvs
	docker rmi game_service-app || true
	docker rmi game_service-migration || true


# DEBUG

.PHONY: run
run:
	go run cmd/main.go --config config/debug.env

.PHONY: test
test:
	go test ./...


.PHONY: infra
infra:
	docker-compose -p game_infra -f deployments/docker/debug/docker-compose.yaml  \
	--env-file ./config/debug.env up -d

.PHONY: migrate
migrate:
	goose -dir migrations postgres \
	"postgresql://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)\
	@$(POSTGRES_HOST_OUTER_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)\
	?sslmode=$(SSL_MODE)" up
.PHONY: run migrate

MOCKGEN_BIN= ""
ENV ?= local
ENV_FILE = ./config/$(ENV).env

include ${ENV_FILE}

# usage: Нужно указать префикс env файла. То есть
# если используем local.env, то пишем make migrate ENV=local.
migrate:
	goose -dir migrations postgres \
	"postgresql://$(POSTGRES_USERNAME):$(POSTRGRES_PASSWORD)\
	@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)\
	?sslmode=disable" up

.PHONY: run
run:
	go run cmd/main.go --config config/local.env

.PHONY: test
test:
	go test ./...

.PHONY: .mock
mock:
	find . -name '*_mock.go' -delete
	mockgen -source internal/service/game/game.go -destination=internal/service/game/mocks/game.go -package=mock_gameservice

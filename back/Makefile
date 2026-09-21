COMPOSE := docker compose

.DEFAULT_GOAL := help

.PHONY: help
help: 
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

## API と Swagger UI を起動
.PHONY: up
up:
	$(COMPOSE) up --build -d

## コンテナを停止する
.PHONY: down
down: 
	$(COMPOSE) down

## コンテナとDBのボリュームごと削除する
.PHONY: clean
clean: 
	$(COMPOSE) down -v

## APIのログを追う
.PHONY: logs
logs: 
	$(COMPOSE) logs -f api

 ## 開発コンテナのシェルに入る
.PHONY: sh
sh:
	$(COMPOSE) exec api sh

 ## 本番相当（ビルド済みバイナリ）で起動する
.PHONY: prod
prod:
	$(COMPOSE) -f compose.prod.yaml up --build

# 開発タスク（コンテナ内で実行）

## api/openapi.yaml から Go のコードを再生成する
.PHONY: generate
generate: 
	$(COMPOSE) run --rm api go tool oapi-codegen --config api/oapi-codegen.yaml api/openapi.yaml

## テストを実行する
.PHONY: test
test: 
	$(COMPOSE) run --rm api go test ./... -count=1

 ## gofmt と golangci-lint をかける
.PHONY: lint
lint:
	$(COMPOSE) run --rm api sh -c 'test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)'
	$(COMPOSE) run --rm api golangci-lint run ./...

 ## go.mod / go.sum を整理する
.PHONY: tidy
tidy:
	$(COMPOSE) run --rm api go mod tidy

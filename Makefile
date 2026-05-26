.PHONY: run build tidy docker-up docker-down docker-rebuild gen-protos swagger

run:
	@go run cmd/main.go

build:
	@CGO_ENABLED=0 go build -o bin/server cmd/main.go

tidy:
	@go mod tidy

swagger:
	@swag init -g cmd/main.go --parseInternal --parseDependency --parseDepth 2 -o docs

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

docker-rebuild:
	@docker compose up -d --build

PROTO_DIR     := $(shell pwd)/protos
PROTO_OUT_DIR := $(shell pwd)/genprotos

gen-protos:
	@if [ -z "$(folder)" ]; then \
		echo "Error: pass folder=<name> (e.g. make gen-protos folder=user)"; \
		exit 1; \
	fi
	@rm -rf "$(PROTO_OUT_DIR)/$(folder)" && mkdir -p "$(PROTO_OUT_DIR)/$(folder)"
	@protoc \
		-I=$(PROTO_DIR) \
		--go_out="$(PROTO_OUT_DIR)" --go_opt=paths=source_relative \
		--go-grpc_out="$(PROTO_OUT_DIR)" --go-grpc_opt=paths=source_relative \
		$$(find "$(PROTO_DIR)/$(folder)" -type f -name "*.proto")

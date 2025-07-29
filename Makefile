SERVICE ?= auth
PROTO_FILE := $(SERVICE).proto
OUT_DIR := services/$(SERVICE)/pkg/pb

up:
	docker compose -f docker-compose-local.yaml build --no-cache && docker compose -f docker-compose-local.yaml up -d

down:
	docker compose down -v

gen:
	@mkdir -p $(OUT_DIR)
	@protoc \
		--proto_path=protobuf \
		--go_out=paths=source_relative:$(OUT_DIR) \
		--go-grpc_out=paths=source_relative:$(OUT_DIR) \
		$(PROTO_FILE)

run:
	@bash -c 'export $$(grep -vE "^\s*#|^\s*$$" .env.local | xargs) && go run services/$(SERVICE)/cmd/main.go'

rest:
	@bash -c 'export $$(grep -vE "^\s*#|^\s*$$" .env.local | xargs) && go run gateway/rest/cmd/main.go'

run_worker:
	@bash -c 'export $$(grep -vE "^\s*#|^\s*$$" .env.local | xargs) && go run services/$(SERVICE)/cmd/worker/main.go'
SERVICE ?= auth
PROTO_FILE := $(SERVICE).proto
OUT_DIR := services/$(SERVICE)/pkg/pb

gen:
	@mkdir -p $(OUT_DIR)
	@protoc \
		--proto_path=protobuf \
		--go_out=paths=source_relative:$(OUT_DIR) \
		--go-grpc_out=paths=source_relative:$(OUT_DIR) \
		$(PROTO_FILE)

run:
	@bash -c 'export $$(grep -vE "^\s*#|^\s*$$" .env | xargs) && go run services/$(SERVICE)/cmd/main.go'

rest:
	@bash -c 'export $$(grep -vE "^\s*#|^\s*$$" .env | xargs) && go run gateway/rest/cmd/main.go'

run_worker:
	@bash -c 'export $$(grep -vE "^\s*#|^\s*$$" .env | xargs) && go run services/$(SERVICE)/cmd/worker/main.go'
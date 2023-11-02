.PHONY: gen build
SERVICE_NAME ?= profiles

gen: gen_proto gen_modules

gen_proto:
	protoc \
	-I api/proto \
	-I api/proto/${SERVICE_NAME} \
	--go_out=api/proto/profiles \
	--go_opt=paths=source_relative \
	--go-grpc_out=require_unimplemented_servers=false:api/proto/profiles \
	--go-grpc_opt=paths=source_relative \
	--validate_opt=paths=source_relative \
	--validate_out=lang=go:api/proto/profiles \
	${SERVICE_NAME}.proto

gen_modules:
	go mod tidy
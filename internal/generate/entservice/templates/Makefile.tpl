.PHONY: grpc
# generate grpc code
grpc:
	protoc -I ../proto -I . --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./proto/example/*.proto

.PHONY: http
# generate kratos-http code
http:
	protoc -I ../proto -I . --go_out=paths=source_relative:. --go-http_out=paths=source_relative:. ./proto/example/*.proto

.PHONY: gateway
# generate grpc-gateway code
gateway:
	protoc -I ../proto -I . --go_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. ./proto/example/*.proto

.PHONY: openapi
# generate openapi documents
openapi:
	protoc -I ../proto -I . --openapi_out=. ./proto/example/*.proto

.PHONY: proto
# generate all proto code
proto:
	make grpc http gateway openapi

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
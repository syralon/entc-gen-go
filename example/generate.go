package example

//go:generate protoc -I ../proto -I . --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --go-http_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. --openapi_out=. ./proto/example/*.proto

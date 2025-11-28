package example

//go:generate protoc -I . -I ../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --go-http_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. ./proto/example/*.proto

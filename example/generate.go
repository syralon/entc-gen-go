package example

//go:generate protoc -I . --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --go-http_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. proto/syralon/entc-gen-go/*.proto

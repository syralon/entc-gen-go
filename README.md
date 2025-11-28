# entc-gen-go

generate go service base on [ent](https://github.com/ent/ent)

# Usage

1. Add [syralon/entproto/*.proto](./proto) into your `[PROTOC]/bin` directory (you can run `whre protoc` or `whereis protoc` to ensure the `%PROTOC%/bin` directory).
2. Run `go install github.com/syralon/entc-gen-go/cmd/entc-gen@latest` to install `entc-gen`.
3. Run `go mod init github.com/example/example` to create a new mod file.
4. Run `entc-gen new Example` to add a new schema, (This command is same as `ent new`), and then modify the generated ent file to add some fields.
5. Run `entc-gen proto` to generate proto files.
6. Run `entc-gen service` to generate service codes.
7. Run `go generate ./...` to generate other codes.
8. Edit the config file named `config.yaml` which is auto generated in current directory to add your own resources.
9. Run `go run ./cmd/example`
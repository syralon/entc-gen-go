package entproto

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/syralon/entc-gen-go/internal/entcgen"
)

func NewEntityBuilder(path, pkg, goPkg string) ProtoFileBuilder {
	return NewFile(
		WithFilename(path),
		WithPackage(pkg),
		WithGoPackage(goPkg),
		WithMessageBuilder(
			NewMessage(
				WithTypeMapping(EntityTypeMapping),
			),
		),
	)
}

func NewGRPCBuilder(filename, pkg, goPkg string) ProtoFileBuilder {
	return NewFile(
		WithFilename(filename),
		WithPackage(pkg),
		WithGoPackage(goPkg),
		WithEnumBuilder(&orderEnumBuilder{}),
		WithMessageBuilder(
			OptionMessages(),
			UpdateMessages(),
			ListOrderMessage(),
			MethodGetMessages(),
			MethodListMessages(),
			MethodCreateMessages(),
			MethodUpdateMessages(),
			MethodDeleteMessages(),
			MethodSetMessages(),
			MethodListEdgesMessage(),
		),
		WithServiceBuilder(
			GRPCServiceBuilder(),
		),
	)
}

func New(module, output string) entcgen.Generator {
	protoPath := path.Join("proto", path.Base(module))
	protoPkg := strings.ReplaceAll(protoPath, "/", ".")
	gopkg := path.Join(module, "proto", path.Base(module))
	generate(module, output)
	return NewProto(
		WithProtoOutput(output),
		WithBuilder(
			NewEntityBuilder(path.Join(protoPath, "ent.proto"), protoPkg, gopkg),
			NewGRPCBuilder(path.Join(protoPath, "ent_service.proto"), protoPkg, gopkg),
		),
	)
}

func generate(module, output string) {
	content := []byte(fmt.Sprintf(
		"package %s\n\n"+
			"//go:generate protoc -I . "+
			"--go_out=paths=source_relative:. "+
			"--go-grpc_out=paths=source_relative:. "+
			"--go-http_out=paths=source_relative:. "+
			"--grpc-gateway_out=paths=source_relative:. "+
			"./proto/%s/*.proto\n",
		path.Base(module),
		path.Base(module),
	))
	_ = os.WriteFile(path.Join(output, "generate.go"), content, 644)
}

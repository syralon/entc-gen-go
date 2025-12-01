package entproto

import (
	"fmt"
	"os"
	"path"

	"github.com/syralon/entc-gen-go/internal/entcgen"
	"github.com/syralon/entc-gen-go/internal/tools/text"
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
	// module = github.com/example/hello
	// output = api
	// --------------------------------------
	// protoModule = proto/example/hello
	// protoPackage = example.hello
	// goPackage = github.com/example/hello/api/proto/example/hello;hello

	protoModule := text.ProtoModule(module)
	protoPackage := text.ProtoPackage(module)
	goPackage := path.Join(module, output, protoModule) + ";" + path.Base(module)
	generate(module, protoModule, output)
	return NewProto(
		WithProtoOutput(output),
		WithBuilder(
			NewEntityBuilder(path.Join(protoModule, "ent.proto"), protoPackage, goPackage),
			NewGRPCBuilder(path.Join(protoModule, "ent_service.proto"), protoPackage, goPackage),
		),
	)
}

func generate(module, protoModule, output string) {
	content := []byte(fmt.Sprintf(
		"package %s\n\n"+
			"//go:generate protoc -I . "+
			"--go_out=paths=source_relative:. "+
			"--go-grpc_out=paths=source_relative:. "+
			"--go-http_out=paths=source_relative:. "+
			"--grpc-gateway_out=paths=source_relative:. "+
			"%s\n",
		path.Base(path.Join(module, output)),
		path.Join(protoModule, "*.proto"),
	))
	_ = os.WriteFile(path.Join(output, "generate.go"), content, 644)
}

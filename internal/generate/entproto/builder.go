package entproto

import (
	"path"

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
	protoPath := path.Join(output, "proto", path.Base(module))
	protoPkg := path.Join()
	gopkg := path.Join(module, output, "proto", path.Base(module))
	return NewProto(
		WithProtoOutput(output),
		WithBuilder(
			NewEntityBuilder(path.Join(protoPath, "ent.proto"), protoPkg, gopkg),
			NewGRPCBuilder(path.Join(protoPath, "ent_service.proto"), protoPkg, gopkg),
		),
	)
}

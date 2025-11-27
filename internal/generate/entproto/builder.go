package entproto

import (
	"os"
	"path"
	"strings"

	"github.com/syralon/entc-gen-go/internal/entgen"
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

func findGoMod() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return ""
}

func lastPath(dir string) string {
	_, d := path.Split(dir)
	if d == "" {
		return dir
	}
	return d
}

func New(output, protoPath, pkg string) entgen.Generator {
	dir := strings.ReplaceAll(path.Join(output, protoPath), "\\", "/")
	_ = os.MkdirAll(dir, 0744)
	if pkg == "" {
		pkg = strings.ReplaceAll(dir, "/", ".")
	}
	var gopkg string
	if module := findGoMod(); module != "" {
		gopkg = path.Join(module, dir) + ";" + lastPath(protoPath)
	}
	if gopkg == "" {
		gopkg = "./" + path.Clean(protoPath)
	}
	return NewProto(
		WithProtoOutput(output),
		WithBuilder(
			NewEntityBuilder(path.Join(protoPath, "ent.proto"), pkg, gopkg),
			NewGRPCBuilder(path.Join(protoPath, "ent_service.proto"), pkg, gopkg),
		),
	)
}

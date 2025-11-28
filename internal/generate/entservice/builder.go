package entservice

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"entgo.io/ent/entc/gen"
	"github.com/dave/jennifer/jen"
	"github.com/syralon/entc-gen-go/internal/entcgen"
)

type Option func(b *builder)

func WithModule(module string) Option {
	return func(b *builder) {
		b.module = module
	}
}

func WithOutput(output string) Option {
	return func(b *builder) {
		b.output = output
	}
}

type builder struct {
	output string
	module string

	templates *template.Template
}

func NewBuilder(opts ...Option) entcgen.Generator {
	b := &builder{
		module: "github.com/syralon/example",
	}
	for _, opt := range opts {
		opt(b)
	}

	var err error
	b.templates, err = template.ParseFS(fs, "templates/*.tpl")
	if err != nil {
		panic(err)
	}

	return b
}

func (b *builder) render(data any, name, filename string) error {
	filename = path.Join(b.output, filename)
	_ = os.MkdirAll(path.Dir(filename), 0700)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	return b.templates.ExecuteTemplate(file, name, data)
}

func (b *builder) renders(data any, m map[string]string) error {
	for name, filename := range m {
		if err := b.render(data, name, filename); err != nil {
			return err
		}
	}
	return nil
}

func (b *builder) write(file *jen.File, filename string) error {
	filename = path.Join(b.output, filename)
	_ = os.MkdirAll(path.Dir(filename), 0700)
	return file.Save(filename)
}

func (b *builder) Generate(ctx context.Context, graph *gen.Graph) error {
	var services = make([]string, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		services = append(services, fmt.Sprintf("%sService", node.Name))
	}
	data := map[string]any{
		"module":               b.module,
		"config_proto_package": path.Base(b.module) + ".config",
		"services":             services,
	}
	if err := b.renders(data, map[string]string{
		"config.go.tpl":           "internal/conf/config.go",
		"config.proto.tpl":        "internal/conf/config.proto",
		"data.go.tpl":             "internal/data/data.go",
		"data_provider.go.tpl":    "internal/data/provider.go",
		"helper.go.tpl":           "internal/service/helper.go",
		"main.go.tpl":             path.Join("cmd", path.Base(b.module), "main.go"),
		"Makefile.tpl":            "Makefile",
		"server_grpc.go.tpl":      "internal/server/grpc.go",
		"server_http.go.tpl":      "internal/server/http.go",
		"server_provider.go.tpl":  "internal/server/provider.go",
		"service_provider.go.tpl": "internal/service/provider.go",
		"wire.go.tpl":             path.Join("cmd", path.Base(b.module), "wire.go"),
	}); err != nil {
		return err
	}
	sb := &serviceBuilder{
		entPackage:   path.Join(b.module, "ent"),
		protoPackage: path.Join(b.module, "proto", path.Base(b.module)),
	}
	for _, node := range graph.Nodes {
		file, err := sb.Build(ctx, node)
		if err != nil {
			return err
		}
		if err = b.write(file, fmt.Sprintf("internal/service/%sservice.go", strings.ToLower(node.Name))); err != nil {
			return err
		}
	}
	if err := b.service(ctx, graph); err != nil {
		return err
	}
	return nil
}

func (b *builder) service(_ context.Context, graph *gen.Graph) error {
	entPkg := path.Join(b.module, "ent")
	var fields []jen.Code
	var blocks []jen.Code
	for _, node := range graph.Nodes {
		fields = append(
			fields,
			jen.Id(fmt.Sprintf("%sService", node.Name)).Op("*").Id(fmt.Sprintf("%sService", node.Name)),
		)
		blocks = append(
			blocks,
			jen.Id(fmt.Sprintf("New%sService", node.Name)).Call(jen.Id("client")).Op(","),
		)
	}

	file := jen.NewFile("service")
	file.Type().Id("Services").Struct(fields...).Line()
	file.Func().Id("NewServices").Params(jen.Id("client").Op("*").Qual(entPkg, "Client")).Op("*").Id("Services").Block(
		jen.Return(jen.Op("&").Id("Services").Block(blocks...)),
	)
	return b.write(file, "internal/service/services.go")
}

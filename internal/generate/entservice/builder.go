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
	"github.com/syralon/entc-gen-go/internal/tools/text"
)

type Option func(b *builder)

func WithModule(module string) Option {
	return func(b *builder) {
		b.rootModule = module
	}
}

func WithOutput(output string) Option {
	return func(b *builder) {
		b.output = output
	}
}
func WithOverwrite(overwrite bool) Option {
	return func(b *builder) {
		b.overwrite = overwrite
	}
}

type builder struct {
	output     string
	overwrite  bool
	rootModule string

	templates *template.Template
}

func NewBuilder(opts ...Option) entcgen.Generator {
	b := &builder{
		rootModule: "",
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

func (b *builder) render(data any, name string, o *output) error {
	filename := path.Join(b.output, o.filename)
	if !o.overwrite {
		if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
			return nil
		}
	}
	_ = os.MkdirAll(path.Dir(filename), 0700)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	return b.templates.ExecuteTemplate(file, name, data)
}

func (b *builder) renders(data any, m map[string]*output) error {
	for name, o := range m {
		if err := b.render(data, name, o); err != nil {
			return err
		}
	}
	return nil
}

func (b *builder) generate(ctx context.Context, node *gen.Type, g interface {
	Build(_ context.Context, node *gen.Type) (*jen.File, error)
}, filename string) error {
	file, err := g.Build(ctx, node)
	if err != nil {
		return err
	}
	return b.write(file, filename)
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
	module := path.Join(b.rootModule, b.output)
	protoModule := path.Join(module, text.ProtoModule(b.rootModule))
	data := map[string]any{
		"module":               module,
		"config_proto_package": path.Base(module) + ".config",
		"proto_package":        protoModule,
		"proto_path":           text.ProtoModule(b.rootModule),
		"services":             services,
	}
	if err := b.renders(data, map[string]*output{
		"config.go.tpl":          out("internal/conf/config.go", b.overwrite),
		"config.proto.tpl":       out("internal/conf/config.proto", b.overwrite),
		"config.yaml.tpl":        out("config.yaml", b.overwrite),
		"data.go.tpl":            out("internal/data/data.go", b.overwrite),
		"data_provider.go.tpl":   out("internal/data/provider.go", b.overwrite),
		"helper.go.tpl":          out("internal/service/helper.go", b.overwrite),
		"main.go.tpl":            out(path.Join("cmd", path.Base(module), "main.go"), b.overwrite),
		"Makefile.tpl":           out("Makefile", b.overwrite),
		"server_grpc.go.tpl":     out("internal/server/grpc.go", b.overwrite),
		"server_http.go.tpl":     out("internal/server/http.go", b.overwrite),
		"server_provider.go.tpl": out("internal/server/provider.go", b.overwrite),
		"wire.go.tpl":            out(path.Join("cmd", path.Base(module), "wire.go"), b.overwrite),
		"wire_gen.go.tpl":        out(path.Join("cmd", path.Base(module), "wire_gen.go"), b.overwrite),
	}); err != nil {
		return err
	}
	sb := &serviceBuilder{
		entPackage:   path.Join(module, "ent"),
		protoPackage: protoModule,
	}
	cb := &controllerBuilder{
		pkg:        module,
		entPackage: path.Join(module, "ent"),
	}
	for _, node := range graph.Nodes {
		if err := b.generate(ctx, node, sb, fmt.Sprintf("internal/service/%sservice.go", strings.ToLower(node.Name))); err != nil {
			return err
		}
		if err := b.generate(ctx, node, cb, fmt.Sprintf("internal/controller/%sservice/service.go", strings.ToLower(node.Name))); err != nil {
			return err
		}
	}
	file, err := controllerProvider(module, path.Join(module, "ent"), protoModule, graph)
	if err != nil {
		return err
	}
	if err = b.write(file, "internal/controller/provider.go"); err != nil {
		return err
	}
	if err = b.service(ctx, module, graph); err != nil {
		return err
	}
	return nil
}

func (b *builder) service(_ context.Context, module string, graph *gen.Graph) error {
	entPkg := path.Join(module, "ent")
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

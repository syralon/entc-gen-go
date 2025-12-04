package entservice

import (
	"bytes"
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

type Builder interface {
	Build(_ context.Context, node *gen.Type) (*jen.File, error)
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

func (b *builder) render(data any, name, filename string) error {
	filename = path.Join(b.output, filename)
	if !b.overwrite {
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

func (b *builder) renders(data any, m map[string]string) error {
	for name, filename := range m {
		if err := b.render(data, name, filename); err != nil {
			return err
		}
	}
	return nil
}

func (b *builder) generate(ctx context.Context, node *gen.Type, g Builder, filename string, overwrite bool) error {
	file, err := g.Build(ctx, node)
	if err != nil {
		return err
	}
	return b.write(file, filename, overwrite)
}

func (b *builder) write(file *jen.File, filename string, overwrite bool) error {
	filename = path.Join(b.output, filename)
	if !overwrite {
		if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
			return nil
		}
	}
	_ = os.MkdirAll(path.Dir(filename), 0700)
	return file.Save(filename)
}

func (b *builder) rewrite(file *jen.File, filename string, overwrite bool) error {
	filename = path.Join(b.output, filename)
	if !overwrite {
		if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
			return nil
		}
	}
	_ = os.MkdirAll(path.Dir(filename), 0700)
	buf := &bytes.Buffer{}
	if err := file.Render(buf); err != nil {
		return err
	}
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	var afterPkg bool
	for _, line := range bytes.Split(buf.Bytes(), []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			_, _ = out.Write([]byte("\n"))
			continue
		}
		if afterPkg {
			line = append([]byte("// "), line...)
		} else {
			afterPkg = bytes.HasPrefix(line, []byte("package"))
		}
		_, _ = out.Write(append(line, '\n'))
	}
	return nil
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
	if err := b.renders(data, map[string]string{
		"config.go.tpl":          "internal/conf/config.go",
		"config.proto.tpl":       "internal/conf/config.proto",
		"config.yaml.tpl":        "config.yaml",
		"data.go.tpl":            "internal/data/data.go",
		"data_provider.go.tpl":   "internal/data/provider.go",
		"helper.go.tpl":          "internal/service/helper.go",
		"main.go.tpl":            path.Join("cmd", path.Base(module), "main.go"),
		"Makefile.tpl":           "Makefile",
		"server_grpc.go.tpl":     "internal/server/grpc.go",
		"server_http.go.tpl":     "internal/server/http.go",
		"server_provider.go.tpl": "internal/server/provider.go",
		"wire.go.tpl":            path.Join("cmd", path.Base(module), "wire.go"),
		"wire_gen.go.tpl":        path.Join("cmd", path.Base(module), "wire_gen.go"),
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
		if err := b.generate(ctx, node, sb, fmt.Sprintf("internal/service/%sservice.go", strings.ToLower(node.Name)), true); err != nil {
			return err
		}
		if err := b.generate(ctx, node, cb, fmt.Sprintf("internal/controller/%sservice/service.go", strings.ToLower(node.Name)), b.overwrite); err != nil {
			return err
		}
		methods, err := customControllerMethod(protoModule, node)
		if err != nil {
			return err
		}
		for filename, m := range methods {
			if err = b.rewrite(m, filename, b.overwrite); err != nil {
				return err
			}
		}
	}
	file, err := controllerProvider(module, path.Join(module, "ent"), protoModule, graph)
	if err != nil {
		return err
	}
	if err = b.write(file, "internal/controller/provider.go", b.overwrite); err != nil {
		return err
	}
	return nil
}

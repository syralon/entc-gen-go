package entproto

import (
	"context"
	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/jhump/protoreflect/v2/protoprint"
	"github.com/syralon/entc-gen-go/internal/tools/text"
	"google.golang.org/protobuf/reflect/protoreflect"
	"path"
)

type ProtoBuilder interface {
	Build(ctx *Context, graph *gen.Graph) ([]*protobuilder.FileBuilder, error)
}

type GenerateOption interface {
	applyGenerator(*Generator)
}

type GenerateOptionFunc func(*Generator)

func (fn GenerateOptionFunc) applyGenerator(g *Generator) {
	fn(g)
}

func WithOutput(output string) GenerateOptionFunc {
	return func(g *Generator) {
		g.output = output
	}
}

func WithPrinter(printer *protoprint.Printer) GenerateOptionFunc {
	return func(g *Generator) {
		g.printer = printer
	}
}

type Generator struct {
	options

	output  string
	path    string
	printer *protoprint.Printer

	builders []ProtoBuilder
}

func NewGenerator(options ...GenerateOption) (*Generator, error) {
	g := &Generator{}
	for _, option := range options {
		option.applyGenerator(g)
	}

	if g.printer == nil {
		g.printer = new(protoprint.Printer)
	}
	if g.output == "" {
		g.output = "."
	}
	module, err := text.Module(g.output)
	if err != nil {
		return nil, err
	}
	if g.path == "" {
		g.path = text.ProtoModule(module)
	}
	if g.protoPackage == "" {
		g.protoPackage = text.ProtoPackage(module)
	}
	if g.goPackage == "" {
		g.goPackage = path.Join(module, g.path) + ";" + path.Base(g.path)
	}
	g.builders = []ProtoBuilder{
		NewEntBuilder(
			WithProtoPackage(g.protoPackage),
			WithGoPackage(g.goPackage),
			WithPath(g.path),
		),
		NewServiceBuilder(
			WithProtoPackage(g.protoPackage),
			WithGoPackage(g.goPackage),
			WithPath(g.path),
		),
	}
	return g, nil
}

func (g *Generator) Generate(c context.Context, graph *gen.Graph) error {
	ctx := NewContext(c)
	var files []*protobuilder.FileBuilder
	for _, bu := range g.builders {
		f, err := bu.Build(ctx, graph)
		if err != nil {
			return err
		}
		files = append(files, f...)
	}

	var descriptors = make([]protoreflect.FileDescriptor, 0, len(files))
	for _, file := range files {
		descriptor, err := file.Build()
		if err != nil {
			return err
		}
		descriptors = append(descriptors, descriptor)
	}
	return g.printer.PrintProtosToFileSystem(descriptors, g.output)
}

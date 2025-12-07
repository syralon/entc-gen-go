package entproto

import (
	"context"

	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/jhump/protoreflect/v2/protoprint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoBuilder interface {
	Build(ctx *Context, graph *gen.Graph) ([]*protobuilder.FileBuilder, error)
}

type Generator struct {
	output       string
	protoPackage string
	goPackage    string
	printer      *protoprint.Printer

	builders []ProtoBuilder
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

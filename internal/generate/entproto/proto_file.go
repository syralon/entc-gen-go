package entproto

import (
	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type builder struct {
	mbs []ProtoMessageBuilder
	ebs []ProtoEnumBuilder
	sbs []ProtoServiceBuilder

	filename string
	pkg      string
	goPkg    string
}

type BuildOption func(*builder)

func WithMessageBuilder(mb ...ProtoMessageBuilder) BuildOption {
	return func(builder *builder) {
		builder.mbs = append(builder.mbs, mb...)
	}
}

func WithServiceBuilder(sb ...ProtoServiceBuilder) BuildOption {
	return func(builder *builder) {
		builder.sbs = append(builder.sbs, sb...)
	}
}

func WithEnumBuilder(eb ...ProtoEnumBuilder) BuildOption {
	return func(builder *builder) {
		builder.ebs = append(builder.ebs, eb...)
	}
}

func WithFilename(filename string) BuildOption {
	return func(builder *builder) {
		if filename != "" {
			builder.filename = filename
		}
	}
}

func WithPackage(pkg string) BuildOption {
	return func(builder *builder) {
		if pkg != "" {
			builder.pkg = pkg
		}
	}
}

func WithGoPackage(goPkg string) BuildOption {
	return func(builder *builder) {
		if goPkg != "" {
			builder.goPkg = goPkg
		}
	}
}

func (b *builder) Build(ctx Context, graph *gen.Graph) (*protobuilder.FileBuilder, error) {
	fb := protobuilder.NewFile(b.filename)
	fb.Syntax = protoreflect.Proto3
	fb.SetPackageName(protoreflect.FullName(b.pkg))
	fb.SetOptions(&descriptorpb.FileOptions{GoPackage: &b.goPkg})

	fc := NewFileContext(ctx, fb)

	for _, eb := range b.ebs {
		for _, node := range graph.Nodes {
			ebs, err := eb.Build(fc, node)
			if err != nil {
				return nil, err
			}
			for _, v := range ebs {
				fb.AddEnum(v)
			}
		}
	}

	var edges = Edges{}
	for _, mb := range b.mbs {
		for _, node := range graph.Nodes {
			ms, eg, err := mb.Build(fc, node)
			if err != nil {
				return nil, err
			}
			for _, m := range ms {
				fb.AddMessage(m)
			}
			edges = append(edges, eg)
		}
	}

	for _, node := range graph.Nodes {
		for _, sb := range b.sbs {
			ss, err := sb.Build(fc, node)
			if err != nil {
				return nil, err
			}
			for _, s := range ss {
				fb.AddService(s)
			}
		}
	}
	if err := edges.SetEdge(fc); err != nil {
		return nil, err
	}
	return fb, nil
}

func NewFile(opts ...BuildOption) ProtoFileBuilder {
	b := &builder{pkg: "proto", goPkg: "./proto"}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

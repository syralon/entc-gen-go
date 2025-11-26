package entproto

import (
	"context"

	"github.com/jhump/protoreflect/v2/protobuilder"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Dependency interface {
	GetMessage(name protoreflect.Name) *protobuilder.MessageBuilder
	GetEnum(name protoreflect.Name) *protobuilder.EnumBuilder
}

type Dependencies []Dependency

func (d Dependencies) GetMessage(name protoreflect.Name) *protobuilder.MessageBuilder {
	for _, dep := range d {
		if m := dep.GetMessage(name); m != nil {
			return m
		}
	}
	return nil
}
func (d Dependencies) GetEnum(name protoreflect.Name) *protobuilder.EnumBuilder {
	for _, dep := range d {
		if m := dep.GetEnum(name); m != nil {
			return m
		}
	}
	return nil
}

type Context interface {
	context.Context

	Dependency

	WithDependency(dep ...*protobuilder.FileBuilder) Context
}

type contextImpl struct {
	context.Context
	deps Dependencies
}

func (c *contextImpl) GetMessage(name protoreflect.Name) *protobuilder.MessageBuilder {
	return c.deps.GetMessage(name)
}
func (c *contextImpl) GetEnum(name protoreflect.Name) *protobuilder.EnumBuilder {
	return c.deps.GetEnum(name)
}

func (c *contextImpl) WithDependency(dep ...*protobuilder.FileBuilder) Context {
	var deps = make(Dependencies, len(dep))
	for i := range dep {
		deps[i] = dep[i]
	}
	return &contextImpl{
		Context: c,
		deps:    deps,
	}
}

func NewContext(ctx context.Context, deps ...Dependency) Context {
	return &contextImpl{Context: ctx, deps: deps}
}

type FileContext struct {
	Context

	file *protobuilder.FileBuilder
}

func (fc *FileContext) GetMessage(name protoreflect.Name) *protobuilder.MessageBuilder {
	m := fc.Context.GetMessage(name)
	if m == nil {
		return fc.file.GetMessage(name)
	}
	if m.ParentFile() != fc.file {
		fc.file.AddDependency(m.ParentFile())
	}
	return m
}

func (fc *FileContext) GetEnum(name protoreflect.Name) *protobuilder.EnumBuilder {
	m := fc.Context.GetEnum(name)
	if m == nil {
		return fc.file.GetEnum(name)
	}
	if m.ParentFile() != fc.file {
		fc.file.AddDependency(m.ParentFile())
	}
	return m
}

func NewFileContext(ctx Context, file *protobuilder.FileBuilder) *FileContext {
	return &FileContext{Context: ctx, file: file}
}

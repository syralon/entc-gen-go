package entproto

import (
	"context"
	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Context struct {
	context.Context

	files []*protobuilder.FileBuilder
}

func (c *Context) With(f ...*protobuilder.FileBuilder) *Context {
	c.files = append(c.files, f...)
	return c
}

func (c *Context) GetMessage(name protoreflect.Name) (*protobuilder.MessageBuilder, error) {
	for _, file := range c.files {
		if m := file.GetMessage(name); m != nil {
			return m, nil
		}
	}
	return nil, ErrorMessageNotFound(name)
}

func (c *Context) GetEnum(name protoreflect.Name) (*protobuilder.EnumBuilder, error) {
	for _, file := range c.files {
		if m := file.GetEnum(name); m != nil {
			return m, nil
		}
	}
	return nil, ErrorEnumNotFound(name)
}

type Delay func(*Context) error

func DelayList(delay ...Delay) Delay {
	return func(c *Context) error {
		for _, fn := range delay {
			if err := fn(c); err != nil {
				return err
			}
		}
		return nil
	}
}

type File struct {
	*protobuilder.FileBuilder
	delay []Delay
}

func (f *File) CallDelay(ctx *Context) error {
	for _, fn := range f.delay {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func WithFileBuilder(f *protobuilder.FileBuilder, delay ...Delay) *File {
	return &File{FileBuilder: f, delay: delay}
}

type ProtoBuilder interface {
	Build(ctx *Context, node *gen.Type) (*File, error)
}

type getter interface {
	GetMessage(name protoreflect.Name) *protobuilder.MessageBuilder
}

func GetMessage(g getter, name protoreflect.Name) (*protobuilder.MessageBuilder, error) {
	if val := g.GetMessage(name); val != nil {
		return val, nil
	}
	return nil, ErrorMessageNotFound(name)
}

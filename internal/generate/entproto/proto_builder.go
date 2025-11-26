package entproto

import (
	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/v2/protobuilder"
	googleapi "google.golang.org/genproto/googleapis/api/annotations"
)

type (
	ProtoFileBuilder interface {
		Build(ctx Context, graph *gen.Graph) (*protobuilder.FileBuilder, error)
	}

	ProtoMessageBuilder interface {
		Build(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error)
	}
	ProtoMessageBuildFunc func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error)

	ProtoServiceBuilder interface {
		Build(ctx *FileContext, node *gen.Type) ([]*protobuilder.ServiceBuilder, error)
	}

	ProtoServiceBuildFunc func(ctx *FileContext, node *gen.Type) ([]*protobuilder.ServiceBuilder, error)

	ProtoMethodBuilder interface {
		Build(ctx *FileContext, node *gen.Type) ([]*protobuilder.MethodBuilder, error)
	}
	ProtoMethodBuildFunc func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MethodBuilder, error)

	ProtoEnumBuilder interface {
		Build(_ *FileContext, node *gen.Type) ([]*protobuilder.EnumBuilder, error)
	}

	Pattern interface {
		Name() string
		Rule(prefix string) (*googleapi.HttpRule, error)
	}
)

func (fn ProtoMessageBuildFunc) Build(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
	return fn(ctx, node)
}

func (fn ProtoServiceBuildFunc) Build(ctx *FileContext, node *gen.Type) ([]*protobuilder.ServiceBuilder, error) {
	return fn(ctx, node)
}

func (fn ProtoMethodBuildFunc) Build(ctx *FileContext, node *gen.Type) ([]*protobuilder.MethodBuilder, error) {
	return fn(ctx, node)
}

type MessageInterceptor func(next ProtoMessageBuilder) ProtoMessageBuilder

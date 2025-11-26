package entproto

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GRPCServiceBuilder() ProtoServiceBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.ServiceBuilder, error) {
		opt, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil {
			return nil, err
		}
		var methods = []ProtoMethodBuilder{
			MethodSet(),
			MethodListEdges(),
		}
		for _, m := range opt.Method.Methods() {
			methods = append(methods, NewMethod(m, WithAPI(opt.Pattern)))
		}

		b := NewService(
			protoreflect.Name(fmt.Sprintf("%sService", node.Name)),
			methods...,
		)
		return b.Build(ctx, node)
	}
}

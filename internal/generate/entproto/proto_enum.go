package entproto

import (
	"fmt"

	"entgo.io/ent/entc/gen"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type orderEnumBuilder struct{}

func (b *orderEnumBuilder) Build(_ *FileContext, node *gen.Type) ([]*protobuilder.EnumBuilder, error) {
	eb := protobuilder.NewEnum(protoreflect.Name(fmt.Sprintf("%sOrder", node.Name)))
	eb = eb.AddValue(protobuilder.NewEnumValue(protoreflect.Name(strcase.ToScreamingSnake(node.Name) + "_ORDER_BY_ID")))
	for _, v := range node.Fields {
		opts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return nil, err
		}
		if !opts.Orderable {
			continue
		}
		name := fmt.Sprintf("%s_ORDER_BY_%s", strcase.ToScreamingSnake(node.Name), strcase.ToScreamingSnake(v.Name))
		eb = eb.AddValue(protobuilder.NewEnumValue(protoreflect.Name(name)))
	}
	return []*protobuilder.EnumBuilder{eb}, nil
}

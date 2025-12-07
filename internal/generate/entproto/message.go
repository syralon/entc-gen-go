package entproto

import (
	"fmt"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	openapiv3 "github.com/google/gnostic/openapiv3"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	"github.com/syralon/entc-gen-go/pkg/annotations/openapi"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MessageOption func(*MessageBuildHelper)

func WithTypeMapping(mapping TypeMapping) MessageOption {
	return func(builder *MessageBuildHelper) {
		builder.mapping = mapping
	}
}

func WithForceOptional(optional bool) MessageOption {
	return func(builder *MessageBuildHelper) {
		builder.optional = optional
	}
}

func WithSkipImmutable(skip bool) MessageOption {
	return func(builder *MessageBuildHelper) {
		builder.skipImmutable = skip
	}
}

func WithSkipFunc(fn func(opt entproto.FieldOptions) bool) MessageOption {
	return func(builder *MessageBuildHelper) {
		builder.skip = fn
	}
}

func WithEdgeName(fn func(*gen.Type) protoreflect.Name) MessageOption {
	return func(m *MessageBuildHelper) {
		m.edgeNameFunc = fn
	}
}

func WithSingleEdge(b bool) MessageOption {
	return func(builder *MessageBuildHelper) {
		builder.singleEdge = b
	}
}
func WithSkipID(b bool) MessageOption {
	return func(builder *MessageBuildHelper) {
		builder.skipID = b
	}
}

type MessageBuildHelper struct {
	edgeNameFunc  func(node *gen.Type) protoreflect.Name
	mapping       TypeMapping
	optional      bool
	skipImmutable bool
	skip          func(opt entproto.FieldOptions) bool
	singleEdge    bool
	skipID        bool
}

func NewMessageBuildHelper(opts ...MessageOption) *MessageBuildHelper {
	mb := &MessageBuildHelper{
		mapping:       EntityTypeMapping,
		optional:      false,
		skipImmutable: false,
		skip:          func(opt entproto.FieldOptions) bool { return false },
	}
	for _, opt := range opts {
		opt(mb)
	}
	if mb.edgeNameFunc == nil {
		mb.edgeNameFunc = func(node *gen.Type) protoreflect.Name { return protoreflect.Name(node.Name) }
	}
	return mb
}

func (b *MessageBuildHelper) Build(ctx *Context, mb *protobuilder.MessageBuilder, node *gen.Type) error {
	if !b.skipID {
		mb.AddField(protobuilder.NewField("id", b.mapping.Mapping(node.IDType.Type)))
	}
	if err := b.fields(mb, node); err != nil {
		return err
	}
	if doc, err := openapi.GetSchema(node.Annotations); err != nil {
		return fmt.Errorf("invalid openapi annotation on entity %s: %w", node.Name, err)
	} else if doc != nil {
		messageOption := &descriptorpb.MessageOptions{}
		proto.SetExtension(messageOption, openapiv3.E_Schema, doc)
		mb.SetOptions(messageOption)
	}
	for _, edge := range node.Edges {
		m, err := ctx.GetMessage(b.edgeNameFunc(edge.Type))
		if err != nil {
			return err
		}
		fb := protobuilder.NewField(protoreflect.Name(edge.Name), protobuilder.FieldTypeMessage(m))
		if !b.singleEdge && !edge.Unique {
			fb.SetRepeated()
		}
		mb.AddField(fb)
	}
	return nil
}

func (b *MessageBuildHelper) fields(mb *protobuilder.MessageBuilder, node *gen.Type) error {
	for _, v := range node.Fields {
		if v.Immutable && b.skipImmutable {
			continue
		}
		if v.Type.Type == field.TypeEnum {
			// TODO
		}
		fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return err
		}
		if b.skip(fieldOpts) {
			continue
		}
		var entType = v.Type.Type
		if fieldOpts.Type > 0 {
			entType = fieldOpts.Type
		}
		fieldType := b.mapping.Mapping(entType)
		if fieldType == nil {
			return fmt.Errorf("unsupported entity type: %s", v.Type.Type)
		}
		fb := protobuilder.NewField(protoreflect.Name(v.Name), fieldType)
		fb.SetComments(protobuilder.Comments{LeadingComment: v.Comment()})
		if fieldOpts.TypeRepeated {
			fb.SetRepeated()
		} else {
			fb.SetProto3Optional(v.Optional || b.optional)
		}
		if doc, err := openapi.GetSchema(v.Annotations); err != nil {
			return fmt.Errorf("invalid openapi annotation on field %s.%s: %w", node.Name, v.Name, err)
		} else if doc != nil {
			propertiesOption := &descriptorpb.FieldOptions{}
			if doc.Description == "" {
				doc.Description = v.Comment()
			}
			proto.SetExtension(propertiesOption, openapiv3.E_Property, doc)
			fb.SetOptions(propertiesOption)
		}
		mb.AddField(fb)
	}
	return nil
}

//type MessageBuilders []*MessageBuildHelper
//
//func (mb MessageBuilders) Build(node *gen.Type) ([]*protobuilder.MessageBuilder, func(*Context) error, error) {
//	var messages []*protobuilder.MessageBuilder
//	var fns []Delay
//	for _, b := range mb {
//		message, delay, err := b.Build(node)
//		if err != nil {
//			return nil, nil, err
//		}
//		messages = append(messages, message)
//		fns = append(fns, delay)
//	}
//	return messages, func(ctx *Context) error {
//		for _, fn := range fns {
//			if err := fn(ctx); err != nil {
//				return err
//			}
//		}
//		return nil
//	}, nil
//}

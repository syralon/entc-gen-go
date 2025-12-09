package entproto

import (
	"fmt"

	"entgo.io/ent/entc/gen"
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

func WithSkipFunc(fn func(f *gen.Field, opt entproto.FieldOptions) bool) MessageOption {
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
func WithSkipEdge(b bool) MessageOption {
	return func(builder *MessageBuildHelper) {
		builder.skipEdge = b
	}
}

type MessageBuildHelper struct {
	edgeNameFunc  func(node *gen.Type) protoreflect.Name
	mapping       TypeMapping
	optional      bool
	skipImmutable bool
	skip          func(f *gen.Field, opt entproto.FieldOptions) bool
	singleEdge    bool
	skipEdge      bool
	skipID        bool
}

func NewMessageBuildHelper(opts ...MessageOption) *MessageBuildHelper {
	mb := &MessageBuildHelper{
		mapping:       EntityTypeMapping,
		optional:      false,
		skipImmutable: false,
		skip:          func(f *gen.Field, opt entproto.FieldOptions) bool { return false },
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
		idField, err := NewField("id", node.ID, b.mapping)
		if err != nil {
			return err
		}
		mb.AddField(idField)
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
	if b.skipEdge {
		return nil
	}
	for _, edge := range node.Edges {
		m := ctx.NewMessage(string(b.edgeNameFunc(edge.Type)))
		fb := protobuilder.NewField(protoreflect.Name(edge.Name), protobuilder.FieldTypeMessage(m))
		if !b.singleEdge && !edge.Unique {
			fb.SetRepeated()
		}
		fb.SetProto3Optional((b.optional || edge.Optional) && !fb.IsRepeated())
		mb.AddField(fb)
	}
	return nil
}

func (b *MessageBuildHelper) fields(mb *protobuilder.MessageBuilder, node *gen.Type) error {
	for _, v := range node.Fields {
		opt, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return err
		}
		if v.Immutable && b.skipImmutable {
			continue
		}
		if b.skip(v, opt) {
			continue
		}
		fb, err := NewField(v.Name, v, b.mapping)
		if err != nil {
			return err
		}
		fb.SetProto3Optional((b.optional || v.Optional) && !fb.IsRepeated())
		fb.SetComments(protobuilder.Comments{LeadingComment: v.Comment()})
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

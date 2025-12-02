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

type MessageBuildOption func(*messages)

func WithNameFunc(fn func(node *gen.Type) protoreflect.Name) MessageBuildOption {
	return func(builder *messages) {
		builder.nameFunc = fn
	}
}

func WithOriginName() MessageBuildOption {
	return WithNameFunc(func(node *gen.Type) protoreflect.Name {
		return protoreflect.Name(node.Name)
	})
}

func WithFormatName(format string) MessageBuildOption {
	return WithNameFunc(func(node *gen.Type) protoreflect.Name {
		return protoreflect.Name(fmt.Sprintf(format, node.Name))
	})
}

func WithStringName(name string) MessageBuildOption {
	return WithNameFunc(func(*gen.Type) protoreflect.Name { return protoreflect.Name(name) })
}

func WithTypeMapping(mapping TypeMapping) MessageBuildOption {
	return func(builder *messages) {
		builder.mapping = mapping
	}
}

func WithForceOptional(optional bool) MessageBuildOption {
	return func(builder *messages) {
		builder.optional = optional
	}
}

func WithSkipImmutable(skip bool) MessageBuildOption {
	return func(builder *messages) {
		builder.skipImmutable = skip
	}
}

func WithSkipFunc(fn func(opt entproto.FieldOptions) bool) MessageBuildOption {
	return func(builder *messages) {
		builder.skip = fn
	}
}

func WithEdgeName(fn func(*gen.Type) protoreflect.Name) MessageBuildOption {
	return func(m *messages) {
		m.edgeNameFunc = fn
	}
}

func WithSingleEdge(b bool) MessageBuildOption {
	return func(builder *messages) {
		builder.singleEdge = b
	}
}
func WithSkipID(b bool) MessageBuildOption {
	return func(builder *messages) {
		builder.skipID = b
	}
}

type messages struct {
	nameFunc      func(node *gen.Type) protoreflect.Name
	edgeNameFunc  func(node *gen.Type) protoreflect.Name
	mapping       TypeMapping
	optional      bool
	skipImmutable bool
	skip          func(opt entproto.FieldOptions) bool
	singleEdge    bool
	skipID        bool
}

func (m *messages) Build(_ *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
	name := m.nameFunc(node)
	mb := protobuilder.NewMessage(name)
	if !m.skipID {
		mb.AddField(protobuilder.NewField("id", m.mapping.Mapping(node.IDType.Type)))
	}
	if err := m.fields(mb, node); err != nil {
		return nil, nil, err
	}
	if doc, err := openapi.GetSchema(node.Annotations); err != nil {
		return nil, nil, fmt.Errorf("invalid openapi annotation on entity %s: %w", node.Name, err)
	} else if doc != nil {
		messageOption := &descriptorpb.MessageOptions{}
		proto.SetExtension(messageOption, openapiv3.E_Schema, doc)
		mb.SetOptions(messageOption)
	}
	edges := NewEdges(mb, node, m.singleEdge, m.edgeNameFunc)
	return []*protobuilder.MessageBuilder{mb}, edges, nil
}

func (m *messages) fields(mb *protobuilder.MessageBuilder, node *gen.Type) error {
	for _, v := range node.Fields {
		if v.Immutable && m.skipImmutable {
			continue
		}
		if v.Type.Type == field.TypeEnum {
			// TODO
		}
		fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return err
		}
		if m.skip(fieldOpts) {
			continue
		}
		var entType = v.Type.Type
		if fieldOpts.Type > 0 {
			entType = fieldOpts.Type
		}
		fieldType := m.mapping.Mapping(entType)
		if fieldType == nil {
			return fmt.Errorf("unsupported entity type: %s", v.Type.Type)
		}
		fb := protobuilder.NewField(protoreflect.Name(v.Name), fieldType)
		fb.SetComments(LeadingComment(v.Comment()))
		if fieldOpts.TypeRepeated {
			fb.SetRepeated()
		} else {
			fb.SetProto3Optional(v.Optional || m.optional)
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

func NewEdges(mb *protobuilder.MessageBuilder, node *gen.Type, single bool, name func(p *gen.Type) protoreflect.Name) Edges {
	var edges Edges
	for _, e := range node.Edges {
		edges = append(edges, &edge{
			name:        protoreflect.Name(e.Name),
			messageName: name(e.Type),
			mb:          mb,
			repeated:    !(single || e.Unique),
		})
	}
	return edges
}

type edge struct {
	name        protoreflect.Name
	messageName protoreflect.Name
	mb          *protobuilder.MessageBuilder
	repeated    bool
}

func (e *edge) SetEdge(ctx Context) error {
	m := ctx.GetMessage(e.messageName)
	if m == nil {
		return fmt.Errorf("message %s not found", e.messageName)
	}
	ff := protobuilder.NewField(e.name, protobuilder.FieldTypeMessage(m))
	if e.repeated {
		ff.SetRepeated()
	}
	e.mb.AddField(ff)
	return nil
}

func NewMessage(options ...MessageBuildOption) ProtoMessageBuilder {
	mb := &messages{
		nameFunc:      func(node *gen.Type) protoreflect.Name { return protoreflect.Name(node.Name) },
		mapping:       EntityTypeMapping,
		optional:      false,
		skipImmutable: false,
		skip:          func(opt entproto.FieldOptions) bool { return false },
	}
	for _, opt := range options {
		opt(mb)
	}
	if mb.edgeNameFunc == nil {
		mb.edgeNameFunc = mb.nameFunc
	}
	return mb
}

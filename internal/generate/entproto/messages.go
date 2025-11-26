package entproto

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	entpb "github.com/syralon/entc-gen-go/proto/syralon/entproto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func OptionMessages() ProtoMessageBuilder {
	return NewMessage(
		WithFormatName("%sOptions"),
		WithTypeMapping(OperationTypeMapping),
		WithForceOptional(true),
		WithSingleEdge(true),
		WithSkipFunc(func(opt entproto.FieldOptions) bool { return !opt.Filterable }),
	)
}

func UpdateOptionMessages() ProtoMessageBuilder {
	return NewMessage(
		WithFormatName("%sUpdateOptions"),
		WithTypeMapping(OperationTypeMapping),
		WithForceOptional(true),
		WithSingleEdge(true),
		WithSkipFunc(func(opt entproto.FieldOptions) bool { return opt.Immutable }),
	)
}

func MethodGetMessages() ProtoMessageBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
		opt, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil || opt.Method&entproto.APIGet == 0 {
			return nil, nil, err
		}
		data := ctx.GetMessage(protoreflect.Name(node.Name))
		if data == nil {
			return nil, nil, fmt.Errorf("message %s not found", node.Name)
		}
		request := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Get%sRequest", node.Name))).
			AddField(protobuilder.NewField("id", EntityTypeMapping.Mapping(node.IDType.Type)))
		response := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Get%sResponse", node.Name))).
			AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)))

		return []*protobuilder.MessageBuilder{request, response}, nil, nil
	}
}

func MethodListMessages() ProtoMessageBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
		data := ctx.GetMessage(protoreflect.Name(node.Name))
		if data == nil {
			return nil, nil, fmt.Errorf("message %s not found", node.Name)
		}
		options := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%sOptions", node.Name)))
		if options == nil {
			return nil, nil, fmt.Errorf("message %sOptions not found", node.Name)
		}
		paginator := protobuilder.FieldTypeImportedMessage((&entpb.Paginator{}).ProtoReflect().Descriptor())
		request := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("List%sRequest", node.Name))).
			AddField(protobuilder.NewField("options", protobuilder.FieldTypeMessage(options))).
			AddField(protobuilder.NewField("paginator", paginator))
		response := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("List%sResponse", node.Name))).
			AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)).SetRepeated())
		return []*protobuilder.MessageBuilder{request, response}, nil, nil
	}
}

func MethodCreateMessages() ProtoMessageBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
		opt, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil || opt.Method&entproto.APICreate == 0 {
			return nil, nil, err
		}
		request, edges, err := NewMessage(
			WithFormatName("Create%sRequest"),
			WithEdgeName(func(node *gen.Type) protoreflect.Name { return protoreflect.Name(node.Name) }),
		).Build(ctx, node)
		if err != nil {
			return nil, nil, err
		}
		response := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Create%sResponse", node.Name))).
			AddField(protobuilder.NewField("id", EntityTypeMapping.Mapping(node.IDType.Type)))
		return append(request, response), edges, nil
	}
}

func MethodUpdateMessages() ProtoMessageBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
		opt, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil || opt.Method&entproto.APIUpdate == 0 {
			return nil, nil, err
		}
		options := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%sUpdateOptions", node.Name)))
		if options == nil {
			return nil, nil, fmt.Errorf("message %sUpdateOptions not found", node.Name)
		}
		request := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Update%sRequest", node.Name))).
			AddField(protobuilder.NewField("id", EntityTypeMapping.Mapping(node.IDType.Type))).
			AddField(protobuilder.NewField("options", protobuilder.FieldTypeMessage(options)))
		response := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Update%sResponse", node.Name)))
		return []*protobuilder.MessageBuilder{request, response}, nil, nil
	}
}

func MethodDeleteMessages() ProtoMessageBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
		opt, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil || opt.Method&entproto.APIDelete == 0 {
			return nil, nil, err
		}
		data := ctx.GetMessage(protoreflect.Name(node.Name))
		if data == nil {
			return nil, nil, fmt.Errorf("message %s not found", node.Name)
		}
		request := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Delete%sRequest", node.Name))).
			AddField(protobuilder.NewField("id", EntityTypeMapping.Mapping(node.IDType.Type)))
		response := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Delete%sResponse", node.Name)))
		return []*protobuilder.MessageBuilder{request, response}, nil, nil
	}
}

func MethodSetMessages() ProtoMessageBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
		options := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%sOptions", node.Name)))
		if options == nil {
			return nil, nil, fmt.Errorf("message %sOptions not found", node.Name)
		}
		var ms = make([]*protobuilder.MessageBuilder, 0, len(node.Fields)*2)
		for _, field := range node.Fields {
			opt, err := entproto.GetFieldOptions(field.Annotations)
			if err != nil {
				return nil, nil, err
			}
			if opt.Immutable || !opt.Settable {
				continue
			}
			request := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Set%s%sRequest", node.Name, strcase.ToCamel(field.Name)))).
				AddField(protobuilder.NewField(protoreflect.Name(strcase.ToSnake(node.Name)+"_id"), EntityTypeMapping.Mapping(node.IDType.Type))).
				AddField(protobuilder.NewField("options", protobuilder.FieldTypeMessage(options))).
				AddField(protobuilder.NewField("set", EntityTypeMapping.Mapping(field.Type.Type)))
			response := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("Set%s%sResponse", node.Name, strcase.ToCamel(field.Name)))).
				AddField(protobuilder.NewField("rows", protobuilder.FieldTypeInt32()))
			ms = append(ms, request, response)
		}
		return ms, nil, nil
	}
}

func MethodListEdgesMessage() ProtoMessageBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MessageBuilder, Edge, error) {
		var ms []*protobuilder.MessageBuilder
		for _, ed := range node.Edges {
			if ed.Unique {
				continue
			}
			edgeName := strcase.ToCamel(ed.Name)
			options := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%sOptions", ed.Type.Name)))
			if options == nil {
				return nil, nil, fmt.Errorf("message %sOptions not found", ed.Type.Name)
			}
			paginator := protobuilder.FieldTypeImportedMessage((&entpb.Paginator{}).ProtoReflect().Descriptor())

			request := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("List%s%sRequest", node.Name, edgeName))).
				AddField(protobuilder.NewField(protoreflect.Name(strcase.ToSnake(node.Name)+"_id"), EntityTypeMapping.Mapping(node.IDType.Type))).
				AddField(protobuilder.NewField("options", protobuilder.FieldTypeMessage(options))).
				AddField(protobuilder.NewField("paginator", paginator))
			ms = append(ms, request)
		}
		return ms, nil, nil
	}
}

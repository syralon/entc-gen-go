package entproto

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/entc-gen-go/internal/tools/text"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

type ServiceBuilder struct {
	builder
}

func (b *ServiceBuilder) Build(_ *Context, node *gen.Type) (*File, error) {
	filename := fmt.Sprintf("ent_%s_service.proto", strings.ToLower(node.Name))
	file := b.NewFile(filename)
	messagesDelay, err := b.messages(file, node)
	if err != nil {
		return nil, err
	}
	if err = b.orderMessage(file, node); err != nil {
		return nil, err
	}

	service := b.NewService(node.Name)
	if err = b.methodGet(file, service, node); err != nil {
		return nil, err
	}
	if err = b.methodList(file, service, node); err != nil {
		return nil, err
	}
	if err = b.methodCreate(file, service, node); err != nil {
		return nil, err
	}
	if err = b.methodUpdate(file, service, node); err != nil {
		return nil, err
	}
	if err = b.methodDelete(file, service, node); err != nil {
		return nil, err
	}
	listDelay, err := b.methodListEdge(file, service, node)
	if err != nil {
		return nil, err
	}
	getDelay, err := b.methodGetEdge(file, service, node)
	if err != nil {
		return nil, err
	}
	if err = b.methodSet(file, service, node); err != nil {
		return nil, err
	}
	return WithFileBuilder(file, DelayList(messagesDelay, listDelay, getDelay)), nil
}

func (b *ServiceBuilder) messages(file *protobuilder.FileBuilder, node *gen.Type) (Delay, error) {
	mbs := MessageBuilders{
		NewMessageBuildHelper(
			WithFormatName("%sOptions"),
			WithTypeMapping(OperationTypeMapping),
			WithForceOptional(true),
			WithSingleEdge(true),
			WithSkipFunc(func(opt entproto.FieldOptions) bool { return !opt.Filterable })),
		NewMessageBuildHelper(
			WithFormatName("%sUpdate"),
			WithTypeMapping(EntityTypeMapping),
			WithForceOptional(true),
			WithSingleEdge(true),
			WithSkipFunc(func(opt entproto.FieldOptions) bool { return opt.Immutable }),
			WithSkipID(true)),
		NewMessageBuildHelper(
			WithFormatName("%sCreate"),
			WithEdgeName(func(node *gen.Type) protoreflect.Name { return protoreflect.Name(node.Name) }),
		),
	}
	messages, delay, err := mbs.Build(node)
	if err != nil {
		return nil, err
	}
	for _, message := range messages {
		if err = file.TryAddMessage(message); err != nil {
			return nil, err
		}
	}
	return delay, nil
}

func (b *ServiceBuilder) orderMessage(file *protobuilder.FileBuilder, node *gen.Type) error {
	eb := protobuilder.NewEnum(protoreflect.Name(fmt.Sprintf("%sOrder", node.Name)))
	eb = eb.AddValue(protobuilder.NewEnumValue(protoreflect.Name(strcase.ToScreamingSnake(node.Name) + "_ORDER_BY_ID")))
	for _, v := range node.Fields {
		opts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return err
		}
		if !opts.Orderable {
			continue
		}
		name := fmt.Sprintf("%s_ORDER_BY_%s", strcase.ToScreamingSnake(node.Name), strcase.ToScreamingSnake(v.Name))
		eb = eb.AddValue(protobuilder.NewEnumValue(protoreflect.Name(name)))
	}
	file.AddEnum(eb)
	orderMessage := protobuilder.NewMessage(protoreflect.Name(fmt.Sprintf("List%sOrder", node.Name))).
		AddField(protobuilder.NewField("by", protobuilder.FieldTypeEnum(eb))).
		AddField(protobuilder.NewField("desc", protobuilder.FieldTypeBool()))
	return file.TryAddMessage(orderMessage)
}

func (b *ServiceBuilder) methodGet(file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	request, response := b.NewMethod(file, service, "Get", "Get"+node.Name)
	data, err := GetMessage(file, protoreflect.Name(node.Name))
	if err != nil {
		return err
	}
	request.AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type]))
	response.AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)))
	return nil
}

func (b *ServiceBuilder) methodList(file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	request, response := b.NewMethod(file, service, "List", "List"+node.Name)
	data, err := GetMessage(file, protoreflect.Name(node.Name))
	if err != nil {
		return err
	}
	options, err := GetMessage(file, protoreflect.Name(fmt.Sprintf("%sOptions", node.Name)))
	if err != nil {
		return err
	}
	order, err := GetMessage(file, protoreflect.Name(fmt.Sprintf("List%sOrder", node.Name)))
	if err != nil {
		return err
	}
	request.
		AddField(protobuilder.NewField("options", protobuilder.FieldTypeMessage(options))).
		AddField(protobuilder.NewField("orders", protobuilder.FieldTypeMessage(order)).SetRepeated()).
		AddField(protobuilder.NewField("paginator", TypePaginator))
	response.
		AddField(protobuilder.NewField("paginator", TypePaginator)).
		AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)).SetRepeated())
	return nil
}

func (b *ServiceBuilder) methodCreate(file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	request, response := b.NewMethod(file, service, "Create", "Create"+node.Name)
	create, err := GetMessage(file, protoreflect.Name(fmt.Sprintf("%sCreate", node.Name)))
	if err != nil {
		return err
	}
	data, err := GetMessage(file, protoreflect.Name(node.Name))
	if err != nil {
		return err
	}
	request.AddField(protobuilder.NewField("create", protobuilder.FieldTypeMessage(create)))
	response.AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)))
	return nil
}

func (b *ServiceBuilder) methodUpdate(file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	request, _ := b.NewMethod(file, service, "Update", "Update"+node.Name)
	update, err := GetMessage(file, protoreflect.Name(fmt.Sprintf("%sUpdate", node.Name)))
	if err != nil {
		return err
	}
	request.AddField(protobuilder.NewField("update", protobuilder.FieldTypeMessage(update)))
	return nil
}

func (b *ServiceBuilder) methodDelete(file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	request, _ := b.NewMethod(file, service, "Delete", "Delete"+node.Name)
	request.AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type]))
	return nil
}

func (b *ServiceBuilder) methodListEdge(file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) (Delay, error) {
	var fns []Delay
	for _, edge := range node.Edges {
		if edge.Unique {
			continue
		}
		request, response := b.NewMethod(file, service, "List"+text.ProtoPackage(edge.Name), "List"+node.Name+text.ProtoPackage(edge.Name))
		request.
			AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type])).
			AddField(protobuilder.NewField("paginator", TypePaginator))
		response.AddField(protobuilder.NewField("paginator", TypePaginator))
		fns = append(fns, func(ctx *Context) error {
			data, err := GetMessage(file, protoreflect.Name(edge.Type.Name))
			if err != nil {
				return err
			}
			options, err := GetMessage(file, protoreflect.Name(fmt.Sprintf("%sOptions", edge.Type.Name)))
			if err != nil {
				return err
			}
			order, err := GetMessage(file, protoreflect.Name(fmt.Sprintf("List%sOrder", edge.Type.Name)))
			if err != nil {
				return err
			}
			request.
				AddField(protobuilder.NewField("options", protobuilder.FieldTypeMessage(options))).
				AddField(protobuilder.NewField("orders", protobuilder.FieldTypeMessage(order)).SetRepeated())
			response.AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)).SetRepeated())
			return nil
		})
	}
	return DelayList(fns...), nil
}

func (b *ServiceBuilder) methodGetEdge(file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) (Delay, error) {
	var fns []Delay
	for _, edge := range node.Edges {
		if edge.Unique {
			continue
		}
		request, response := b.NewMethod(file, service, "Get"+text.ProtoPackage(edge.Name), "Get"+node.Name+text.ProtoPackage(edge.Name))
		request.
			AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type])).
			AddField(protobuilder.NewField(protoreflect.Name(strcase.ToSnake(edge.Type.Name)+"_id"), EntityTypeMapping[edge.Type.IDType.Type]))
		fns = append(fns, func(ctx *Context) error {
			data, err := GetMessage(file, protoreflect.Name(edge.Type.Name))
			if err != nil {
				return err
			}
			response.AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)))
			return nil
		})
	}
	return DelayList(fns...), nil
}

func (b *ServiceBuilder) methodSet(file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	for _, v := range node.Fields {
		opt, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return err
		}
		if opt.Immutable || !opt.Settable {
			continue
		}
		request, _ := b.NewMethod(file, service, "Set"+text.ProtoPackage(v.Name), "Set"+node.Name+text.ProtoPackage(v.Name))
		request.
			AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type])).
			AddField(protobuilder.NewField("set", EntityTypeMapping[v.Type.Type]))
	}
	return nil
}

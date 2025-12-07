package entproto

import (
	"fmt"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/entc-gen-go/internal/tools/text"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	googleapi "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ServiceBuildOption interface {
	applyService(*ServiceBuilder)
}

type ServiceBuildOptionFunc func(*ServiceBuilder)

func (f ServiceBuildOptionFunc) applyService(builder *ServiceBuilder) {
	f(builder)
}

type ServiceBuilder struct {
	builder
}

func NewServiceBuilder(options ...ServiceBuildOption) *ServiceBuilder {
	s := &ServiceBuilder{}
	for _, option := range options {
		option.applyService(s)
	}
	return s
}

func (b *ServiceBuilder) Build(ctx *Context, graph *gen.Graph) ([]*protobuilder.FileBuilder, error) {
	var files = make([]*protobuilder.FileBuilder, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		file, err := b.build(ctx, node)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (b *ServiceBuilder) build(ctx *Context, node *gen.Type) (*protobuilder.FileBuilder, error) {
	filename := fmt.Sprintf("ent_%s_service.proto", strings.ToLower(node.Name))
	file := ctx.NewFile(filename, b.protoPackage, b.goPackage)
	if err := b.messages(ctx, file, node); err != nil {
		return nil, err
	}
	service := ctx.NewService(fmt.Sprintf("%sService", node.Name))
	file.AddService(service)
	if err := b.methods(ctx, file, service, node); err != nil {
		return nil, err
	}
	if err := b.methodEdges(ctx, file, service, node); err != nil {
		return nil, err
	}
	if err := b.methodSet(ctx, file, service, node); err != nil {
		return nil, err
	}
	return file, nil
}

func (b *ServiceBuilder) messages(ctx *Context, file *protobuilder.FileBuilder, node *gen.Type) (err error) {
	options := ctx.NewMessage(fmt.Sprintf("%sOptions", node.Name))
	update := ctx.NewMessage(fmt.Sprintf("%sUpdate", node.Name))
	create := ctx.NewMessage(fmt.Sprintf("%sCreate", node.Name))

	orderEnum := ctx.NewEnum(fmt.Sprintf("%sOrder", node.Name))
	orderEnum = orderEnum.AddValue(protobuilder.NewEnumValue(protoreflect.Name(strcase.ToScreamingSnake(node.Name) + "_ORDER_BY_ID")))
	for _, v := range node.Fields {
		opts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return err
		}
		if !opts.Orderable {
			continue
		}
		name := fmt.Sprintf("%s_ORDER_BY_%s", strcase.ToScreamingSnake(node.Name), strcase.ToScreamingSnake(v.Name))
		orderEnum = orderEnum.AddValue(protobuilder.NewEnumValue(protoreflect.Name(name)))
	}
	orderMessage := ctx.NewMessage(fmt.Sprintf("List%sOrder", node.Name)).
		AddField(protobuilder.NewField("by", protobuilder.FieldTypeEnum(orderEnum))).
		AddField(protobuilder.NewField("desc", protobuilder.FieldTypeBool()))

	err = NewMessageBuildHelper(
		WithTypeMapping(OperationTypeMapping),
		WithForceOptional(true),
		WithSingleEdge(true),
		WithSkipFunc(func(opt entproto.FieldOptions) bool { return !opt.Filterable }),
	).Build(ctx, options, node)
	if err != nil {
		return err
	}
	err = NewMessageBuildHelper(
		WithTypeMapping(EntityTypeMapping),
		WithForceOptional(true),
		WithSingleEdge(true),
		WithSkipFunc(func(opt entproto.FieldOptions) bool { return opt.Immutable }),
		WithSkipID(true),
	).Build(ctx, update, node)
	if err != nil {
		return err
	}
	err = NewMessageBuildHelper(
		WithEdgeName(func(node *gen.Type) protoreflect.Name { return protoreflect.Name(node.Name) }),
	).Build(ctx, create, node)
	if err != nil {
		return err
	}

	file.AddMessage(options).AddMessage(update).AddMessage(create).AddMessage(orderMessage).AddEnum(orderEnum)
	return err
}

func (b *ServiceBuilder) methods(ctx *Context, file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	apiOpts, err := entproto.GetAPIOptions(node.Annotations)
	if err != nil {
		return err
	}

	m := map[entproto.APIMethod]func(ctx *Context, noe *gen.Type, request, response *protobuilder.MessageBuilder) error{
		entproto.GET:    b.methodGetMessages,
		entproto.LIST:   b.methodListMessages,
		entproto.CREATE: b.methodCreateMessages,
		entproto.UPDATE: b.methodUpdateMessages,
		entproto.DELETE: b.methodDeleteMessages,
	}

	for _, method := range apiOpts.Method.Methods() {
		request, response := ctx.NewMethod(service, method.Name(), method.Name()+node.Name, b.methodOptions(node, method))
		file.AddMessage(request).AddMessage(response)
		h, ok := m[method]
		if !ok {
			return fmt.Errorf("unknown method %s", method)
		}
		if err = h(ctx, node, request, response); err != nil {
			return err
		}
	}
	return nil
}

func (b *ServiceBuilder) methodEdges(ctx *Context, file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	for _, edge := range node.Edges {
		var method = "List"
		if edge.Unique {
			method = "Get"
		}

		opts := b.methodEdgeOptions(node, edge)
		request, response := ctx.NewMethod(service, method+text.ProtoPascal(edge.Name), fmt.Sprintf("%s%s%s", method, node.Name, text.ProtoPascal(edge.Name)), opts)
		file.AddMessage(request).AddMessage(response)

		data := ctx.NewMessage(edge.Type.Name)
		if edge.Unique {
			request.
				AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type])).
				AddField(protobuilder.NewField(protoreflect.Name(strcase.ToSnake(edge.Type.Name)+"_id"), EntityTypeMapping[edge.Type.IDType.Type]))
			response.AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)))
		} else {
			options := ctx.NewMessage(fmt.Sprintf("%sOptions", edge.Type.Name))
			order := ctx.NewMessage(fmt.Sprintf("List%sOrder", edge.Type.Name))
			request.
				AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type])).
				AddField(protobuilder.NewField("paginator", TypePaginator))
			request.
				AddField(protobuilder.NewField("options", protobuilder.FieldTypeMessage(options))).
				AddField(protobuilder.NewField("orders", protobuilder.FieldTypeMessage(order)).SetRepeated())
			response.AddField(protobuilder.NewField("paginator", TypePaginator))
			response.AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)).SetRepeated())
		}
	}
	return nil
}

func (b *ServiceBuilder) methodSet(ctx *Context, file *protobuilder.FileBuilder, service *protobuilder.ServiceBuilder, node *gen.Type) error {
	for _, v := range node.Fields {
		fieldOpt, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return err
		}
		if fieldOpt.Immutable || !fieldOpt.Settable {
			continue
		}
		opts := b.methodSetOptions(node, v.Name)
		request, response := ctx.NewMethod(service, "Set"+text.ProtoPackage(v.Name), "Set"+node.Name+text.ProtoPackage(v.Name), opts)
		file.AddMessage(request).AddMessage(response)
		request.
			AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type])).
			AddField(protobuilder.NewField("set", EntityTypeMapping[v.Type.Type]))
	}
	return nil
}

func (b *ServiceBuilder) methodGetMessages(ctx *Context, node *gen.Type, request, response *protobuilder.MessageBuilder) error {
	data, err := ctx.GetMessage(protoreflect.Name(node.Name))
	if err != nil {
		return err
	}
	request.AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type]))
	response.AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)))
	return nil
}
func (b *ServiceBuilder) methodListMessages(ctx *Context, node *gen.Type, request, response *protobuilder.MessageBuilder) error {
	data, err := ctx.GetMessage(protoreflect.Name(node.Name))
	if err != nil {
		return err
	}
	options, err := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%sOptions", node.Name)))
	if err != nil {
		return err
	}
	order, err := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("List%sOrder", node.Name)))
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
func (b *ServiceBuilder) methodCreateMessages(ctx *Context, node *gen.Type, request, response *protobuilder.MessageBuilder) error {
	create, err := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%sCreate", node.Name)))
	if err != nil {
		return err
	}
	data, err := ctx.GetMessage(protoreflect.Name(node.Name))
	if err != nil {
		return err
	}
	request.AddField(protobuilder.NewField("create", protobuilder.FieldTypeMessage(create)))
	response.AddField(protobuilder.NewField("data", protobuilder.FieldTypeMessage(data)))
	return nil
}
func (b *ServiceBuilder) methodUpdateMessages(ctx *Context, node *gen.Type, request, response *protobuilder.MessageBuilder) error {
	update, err := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%sUpdate", node.Name)))
	if err != nil {
		return err
	}
	request.AddField(protobuilder.NewField("update", protobuilder.FieldTypeMessage(update)))
	return nil
}
func (b *ServiceBuilder) methodDeleteMessages(ctx *Context, node *gen.Type, request, response *protobuilder.MessageBuilder) error {
	request.AddField(protobuilder.NewField("id", EntityTypeMapping[node.IDType.Type]))
	return nil
}

func (b *ServiceBuilder) methodOptions(node *gen.Type, m entproto.APIMethod) *descriptor.MethodOptions {
	apiOpt, err := entproto.GetAPIOptions(node.Annotations)
	if err != nil {
		return nil
	}
	if apiOpt.Pattern == "" {
		return nil
	}
	rule, err := m.Rule(path.Join(apiOpt.Pattern, strings.ToLower(node.Name)))
	if err != nil {
		return nil
	}
	properties := &descriptor.MethodOptions{}
	proto.SetExtension(properties, googleapi.E_Http, rule)
	return properties
}

func (b *ServiceBuilder) methodEdgeOptions(node *gen.Type, edge *gen.Edge) *descriptor.MethodOptions {
	apiOpt, err := entproto.GetAPIOptions(node.Annotations)
	if err != nil {
		return nil
	}
	if apiOpt.Pattern == "" {
		return nil
	}
	var p string
	if edge.Unique {
		p = path.Join(apiOpt.Pattern, strings.ToLower(node.Name), "{id}", strings.ToLower(strcase.ToSnake(edge.Name)), fmt.Sprintf("{%s_id}", strcase.ToSnake(edge.Type.Name)))
	} else {
		p = path.Join(apiOpt.Pattern, strings.ToLower(node.Name), "{id}", strings.ToLower(strcase.ToSnake(edge.Name)))
	}

	rule := &googleapi.HttpRule{Pattern: &googleapi.HttpRule_Get{Get: p}}
	properties := &descriptor.MethodOptions{}
	proto.SetExtension(properties, googleapi.E_Http, rule)
	return properties
}

func (b *ServiceBuilder) methodSetOptions(node *gen.Type, fieldName string) *descriptor.MethodOptions {
	apiOpt, err := entproto.GetAPIOptions(node.Annotations)
	if err != nil {
		return nil
	}
	if apiOpt.Pattern == "" {
		return nil
	}
	rule := &googleapi.HttpRule{Pattern: &googleapi.HttpRule_Put{
		Put: path.Join(apiOpt.Pattern, strings.ToLower(node.Name), "{id}", strings.ToLower(strcase.ToCamel(fieldName))),
	}}
	properties := &descriptor.MethodOptions{}
	proto.SetExtension(properties, googleapi.E_Http, rule)
	return properties
}

package entproto

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	googleapi "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"path"
	"strings"
)

func MethodSet() ProtoMethodBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MethodBuilder, error) {
		apiOpt, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil {
			return nil, err
		}
		var methods = make([]*protobuilder.MethodBuilder, 0, len(node.Fields))
		for _, field := range node.Fields {
			opt, err := entproto.GetFieldOptions(field.Annotations)
			if err != nil {
				return nil, err
			}
			if !opt.Settable {
				continue
			}
			name := strcase.ToCamel(field.Name)
			request := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("Set%s%sRequest", node.Name, name)))
			if request == nil {
				return nil, fmt.Errorf("message Set%sRequest not found", name)
			}
			response := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("Set%s%sResponse", node.Name, name)))
			if response == nil {
				return nil, fmt.Errorf("message Set%s%sResponse not found", node.Name, name)
			}
			method := protobuilder.NewMethod(
				protoreflect.Name(fmt.Sprintf("Set%s", name)),
				protobuilder.RpcTypeMessage(request, false),
				protobuilder.RpcTypeMessage(response, false),
			)
			if apiOpt.Pattern != "" {
				pattern := path.Join(apiOpt.Pattern, strings.ToLower(node.Name), "{id}", strings.ToLower(field.Name))
				rule := &googleapi.HttpRule{Pattern: &googleapi.HttpRule_Put{Put: pattern}, Body: "*"}
				properties := &descriptor.MethodOptions{}
				proto.SetExtension(properties, googleapi.E_Http, rule)
				method.SetOptions(properties)
			}
			methods = append(methods, method)
		}
		return methods, nil
	}
}

func MethodListEdges() ProtoMethodBuildFunc {
	return func(ctx *FileContext, node *gen.Type) ([]*protobuilder.MethodBuilder, error) {
		apiOpt, err := entproto.GetAPIOptions(node.Annotations)
		if err != nil {
			return nil, err
		}

		var methods = make([]*protobuilder.MethodBuilder, 0, len(node.Edges))
		for _, ed := range node.Edges {
			if ed.Unique {
				continue
			}
			name := strcase.ToCamel(ed.Name)
			request := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("List%s%sRequest", node.Name, name)))
			if request == nil {
				return nil, fmt.Errorf("message List%s%sRequest not found", node.Name, name)
			}
			response := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("List%sResponse", ed.Type.Name)))
			if response == nil {
				return nil, fmt.Errorf("message List%sResponse not found", ed.Type.Name)
			}
			method := protobuilder.NewMethod(
				protoreflect.Name(fmt.Sprintf("List%s", name)),
				protobuilder.RpcTypeMessage(request, false),
				protobuilder.RpcTypeMessage(response, false),
			)
			if apiOpt.Pattern != "" {
				pattern := path.Join(
					apiOpt.Pattern,
					strings.ToLower(node.Name),
					"{id}",
					strings.ToLower(strcase.ToCamel(ed.Name)),
				)
				rule := &googleapi.HttpRule{Pattern: &googleapi.HttpRule_Get{Get: pattern}}
				properties := &descriptor.MethodOptions{}
				proto.SetExtension(properties, googleapi.E_Http, rule)
				method.SetOptions(properties)
			}
			methods = append(methods, method)
		}
		return methods, nil
	}
}

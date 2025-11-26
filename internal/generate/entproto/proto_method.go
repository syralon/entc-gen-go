package entproto

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/v2/protobuilder"
	googleapi "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"path"
	"strings"
)

type method struct {
	pattern Pattern
	withAPI bool
	prefix  string
}

func (m *method) Build(ctx *FileContext, node *gen.Type) ([]*protobuilder.MethodBuilder, error) {
	name := m.pattern.Name()
	request := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%s%sRequest", name, node.Name)))
	if request == nil {
		return nil, fmt.Errorf("message %s%sRequest not found", name, node.Name)
	}
	response := ctx.GetMessage(protoreflect.Name(fmt.Sprintf("%s%sResponse", name, node.Name)))
	if response == nil {
		return nil, fmt.Errorf("message %s%sResponse not found", name, node.Name)
	}
	mm := protobuilder.NewMethod(
		protoreflect.Name(name),
		protobuilder.RpcTypeMessage(request, false),
		protobuilder.RpcTypeMessage(response, false),
	)
	if m.withAPI {
		rule, err := m.pattern.Rule(path.Join(m.prefix, strings.ToLower(node.Name)))
		if err != nil {
			return nil, err
		}
		properties := &descriptor.MethodOptions{}
		proto.SetExtension(properties, googleapi.E_Http, rule)
		mm.SetOptions(properties)
	}
	return []*protobuilder.MethodBuilder{mm}, nil
}

type MethodOption func(*method)

func WithAPI(prefix string) MethodOption {
	return func(m *method) {
		if prefix == "" {
			return
		}
		m.withAPI = true
		m.prefix = prefix
	}
}

func NewMethod(pattern Pattern, opts ...MethodOption) ProtoMethodBuilder {
	mm := &method{
		pattern: pattern,
	}
	for _, opt := range opts {
		opt(mm)
	}
	return mm
}

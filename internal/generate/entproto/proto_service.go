package entproto

import (
	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type serviceBuilder struct {
	name    protoreflect.Name
	methods []ProtoMethodBuilder
}

func NewService(name protoreflect.Name, methods ...ProtoMethodBuilder) ProtoServiceBuilder {
	return &serviceBuilder{name: name, methods: methods}
}

func (s *serviceBuilder) AddMethod(m ProtoMethodBuilder) {
	s.methods = append(s.methods, m)
}

func (s *serviceBuilder) Build(ctx *FileContext, node *gen.Type) ([]*protobuilder.ServiceBuilder, error) {
	service := protobuilder.NewService(s.name)
	for _, mb := range s.methods {
		methods, err := mb.Build(ctx, node)
		if err != nil {
			return nil, err
		}
		for _, mt := range methods {
			service.AddMethod(mt)
		}
	}
	return []*protobuilder.ServiceBuilder{service}, nil
}

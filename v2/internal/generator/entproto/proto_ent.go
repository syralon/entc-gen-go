package entproto

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"strings"
)

type EntBuilder struct {
	builder
}

func (b *EntBuilder) Build(_ *Context, node *gen.Type) (*File, error) {
	filename := fmt.Sprintf("ent_%s.proto", strings.ToLower(node.Name))
	file := b.NewFile(filename)
	message, delay, err := NewMessageBuildHelper().Build(node)
	if err != nil {
		return nil, err
	}
	file.AddMessage(message)
	return WithFileBuilder(file, delay), nil
}

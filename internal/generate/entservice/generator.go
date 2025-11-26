package entservice

import (
	"context"
	"entgo.io/ent/entc/gen"
	"github.com/dave/jennifer/jen"
)

type Builder interface {
	Build(ctx context.Context, node *gen.Type) (*jen.File, error)
}

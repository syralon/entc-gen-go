package entcgen

import (
	"context"

	"entgo.io/ent/entc/gen"
)

type Generator interface {
	Generate(ctx context.Context, graph *gen.Graph) error
}

type Generators []Generator

func (g Generators) Generate(ctx context.Context, graph *gen.Graph) error {
	for _, v := range g {
		if err := v.Generate(ctx, graph); err != nil {
			return err
		}
	}
	return nil
}

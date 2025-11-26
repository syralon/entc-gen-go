package entservice

import (
	"context"
	"entgo.io/ent/entc/gen"
	"github.com/dave/jennifer/jen"
)

type HelperBuilder struct{}

func (b *HelperBuilder) Build(_ context.Context, _ *gen.Type) (*jen.File, error) {
	// func Trans[A, B any](a []A, fn func(A) B) []B {
	//	b := make([]B, 0, len(a))
	//	for _, v := range a {
	//		b = append(b, fn(v))
	//	}
	//	return b
	// }
	file := jen.NewFile("service")
	file.Func().Id("Trans").
		Types(jen.List(jen.Id("A"), jen.Id("B")).Any()).
		Params(
			jen.Id("a").Op("[]").Id("A"),
			jen.Id("fn").Func().Params(jen.Id("A")).Id("B"),
		).
		Op("[]").Id("B").
		Block(
			jen.Id("b").Op(":=").Id("make").Call(jen.Op("[]").Id("B"), jen.Id("0"), jen.Id("len").Call(jen.Id("a"))),
			jen.For(
				jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Id("range").Id("a").Block(
					jen.Id("b").Op("=").Id("append").Call(jen.Id("b"), jen.Id("fn").Call(jen.Id("v"))),
				),
			),
			jen.Return(jen.Id("b")),
		)
	return file, nil
}

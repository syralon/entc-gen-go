package entservice

import "github.com/dave/jennifer/jen"

func chain(names ...string) *jen.Statement {
	var st *jen.Statement
	for i, name := range names {
		if i == 0 {
			st = jen.Id(name)
		} else {
			st = st.Dot(name)
		}
	}
	return st
}

func list(names ...string) *jen.Statement {
	var st []jen.Code
	for _, name := range names {
		st = append(st, jen.Id(name))
	}
	return jen.List(st...)
}

func assign(names ...string) *jen.Statement {
	return list(names...).Op("=")
}

func define(names ...string) *jen.Statement {
	return list(names...).Op(":=")
}

func calls(codes ...*jen.Statement) *jen.Statement {
	st := jen.Op("(").Id("\n")
	for _, code := range codes {
		st = st.Add(code.Op(",").Id("\n"))
	}
	return st.Op(")")
}

func ifErr() *jen.Statement {
	return jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
		jen.Return(jen.Id("nil"), jen.Id("err")),
	)
}

func structPtr(name, key, val *jen.Statement) *jen.Statement {
	return jen.Op("&").Add(name).Op("{").Add(key).Op(":").Add(val).Op("}")
}

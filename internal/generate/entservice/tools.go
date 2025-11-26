package entservice

import "github.com/dave/jennifer/jen"

func chain(names ...string) *jen.Statement {
	var st *jen.Statement
	for i, name := range names {
		if i == 0 {
			st = jen.Id(name)
		} else {
			st = st.Op(".").Id(name)
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

func join(st *jen.Statement, args ...jen.Code) *jen.Statement {
	*st = append(*st, args...)
	return st
}

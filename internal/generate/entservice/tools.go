package entservice

import (
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/dave/jennifer/jen"
)

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

func EntID(n *gen.Type, v *jen.Statement) *jen.Statement {
	switch n.ID.Type.Type {
	case field.TypeInt8:
		return jen.Int8().Call(v)
	case field.TypeInt16:
		return jen.Int16().Call(v)
	case field.TypeInt32:
		return v
	case field.TypeInt:
		return jen.Int().Call(v)
	case field.TypeInt64:
		return v
	case field.TypeUint8:
		return jen.Uint8().Call(v)
	case field.TypeUint16:
		return jen.Uint16().Call(v)
	case field.TypeUint32:
		return v
	case field.TypeUint:
		return jen.Uint().Call(v)
	case field.TypeUint64:
		return v
	case field.TypeString:
		return v
	case field.TypeUUID:
		return v
	default:
		//panic(fmt.Errorf("unsupported ent id type: %s", n.ID.Type.Type))
		return v
	}
}

func ProtoID(n *gen.Type, v *jen.Statement) *jen.Statement {
	switch n.ID.Type.Type {
	case field.TypeInt8, field.TypeInt16, field.TypeInt32:
		return jen.Int32().Call(v)
	case field.TypeInt, field.TypeInt64:
		return jen.Int64().Call(v)
	case field.TypeUint8, field.TypeUint32:
		return jen.Uint32().Call(v)
	case field.TypeUint, field.TypeUint64:
		return jen.Uint64().Call(v)
	case field.TypeString:
		return v
	case field.TypeUUID:
		return v.Dot("String").Call()
	default:
		//panic(fmt.Errorf("unsupported ent id type: %s", n.ID.Type.Type))
		return v
	}
}

func ternary(v bool, a, b *jen.Statement) *jen.Statement {
	if v {
		return a
	}
	return b
}

package entservice

import (
	"context"
	"fmt"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"

	"github.com/syralon/entc-gen-go/internal/tools/text"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
)

type serviceBuilder struct {
	entPackage   string
	protoPackage string
}

func (s *serviceBuilder) Build(_ context.Context, node *gen.Type) (*jen.File, error) {
	opt, err := entproto.GetAPIOptions(node.Annotations)
	if err != nil {
		return nil, err
	}

	file := jen.NewFile("service")

	file.ImportAlias(s.protoPackage, "pb")

	s.orderMapping(file, node)
	s.funcToProto(file, node)
	s.funcFromProto(file, node)
	s.serviceStruct(file, node)

	s.funcSet(file, node)
	s.funcListEdge(file, node)

	for _, m := range opt.Method.Methods() {
		switch m {
		case entproto.GET:
			s.funcGet(file, node)
		case entproto.LIST:
			s.funcList(file, node)
		case entproto.CREATE:
			s.funcCreate(file, node)
		case entproto.UPDATE:
			s.funcUpdate(file, node)
		case entproto.DELETE:
			s.funcDelete(file, node)
		default:
		}
	}
	return file, nil
}

func (s *serviceBuilder) funcToProto(file *jen.File, node *gen.Type) {
	defer file.Line()

	var fields []jen.Code
	for _, fi := range node.Fields {
		v := jen.Id(text.ProtoPascal(fi.Name)).Op(":")
		if fi.Type.Type == field.TypeTime {
			v = v.Qual(timestamppb, "New").
				Call(jen.Id("data").Op(".").Id(text.EntPascal(fi.Name))).Op(",")
		} else {
			v = v.Id("data").Op(".").Id(text.EntPascal(fi.Name)).Op(",")
		}
		fields = append(fields, v)
	}

	file.Func().
		Id(fmt.Sprintf("%sToProto", node.Name)).
		Params(jen.Id("data").Op("*").Qual(s.entPackage, text.EntPascal(node.Name))).
		Op("*").Qual(s.protoPackage, text.ProtoPascal(node.Name)).
		Block(
			jen.Return(
				jen.Op("&").Qual(s.protoPackage, text.ProtoPascal(node.Name)).Block(fields...),
			),
		)
}

func (s *serviceBuilder) funcFromProto(file *jen.File, node *gen.Type) {
	defer file.Line()

	var fields []jen.Code
	for _, fi := range node.Fields {
		v := jen.Id(text.EntPascal(fi.Name)).Op(":").Id("data").Op(".").Id(text.ProtoPascal(fi.Name))
		if fi.Type.Type == field.TypeTime {
			v = v.Op(".").Id("AsTime()")
		}
		fields = append(fields, v.Op(","))
	}

	file.Func().
		Id(fmt.Sprintf("%sFromProto", node.Name)).
		Params(jen.Id("data").Op("*").Qual(s.protoPackage, text.ProtoPascal(node.Name))).
		Op("*").Qual(s.entPackage, text.EntPascal(node.Name)).
		Block(
			jen.Return(
				jen.Op("&").Qual(s.entPackage, text.EntPascal(node.Name)).Block(fields...),
			),
		)
}

func (s *serviceBuilder) serviceStruct(file *jen.File, node *gen.Type) {
	defer file.Line()

	file.Type().Id(fmt.Sprintf("%sService", node.Name)).Struct(
		jen.Qual(s.protoPackage, fmt.Sprintf("Unimplemented%sServiceServer", node.Name)),
		jen.Id("client").Op("*").Qual(s.entPackage, fmt.Sprintf("%sClient", node.Name)),
	)
	file.Line()
	file.Func().Id(fmt.Sprintf("New%sService", node.Name)).
		Params(
			jen.Id("client").Op("*").Qual(s.entPackage, "Client"),
		).Op("*").Id(fmt.Sprintf("%sService", node.Name)).
		Block(jen.Return(
			jen.Op("&").Id(fmt.Sprintf("%sService", node.Name)).Block(
				jen.Id("client").Op(":").Id("client").Op(".").Id(node.Name).Op(","),
			),
		))
}

func (s *serviceBuilder) funcGet(file *jen.File, node *gen.Type) {
	defer file.Line()

	fn := s.serviceFunc(file, node.Name, "Get", fmt.Sprintf("Get%sRequest", node.Name), fmt.Sprintf("Get%sResponse", node.Name))
	fn = fn.Block(
		// data, err := s.client.Get(ctx, int(request.GetId()))
		jen.List(jen.Id("data"), jen.Id("err")).Op(":=").Id("s").Op(".").Id("client").Op(".").Id("Get").
			Call(jen.Id("ctx"), jen.Int().Op("(").Id("request").Op(".").Id("GetId()").Op(")")),
		// if err != nil {}
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		),
		jen.Return(
			structPtr(jen.Qual(s.protoPackage, fmt.Sprintf("Get%sResponse", node.Name)), jen.Id("Data"), jen.Id(fmt.Sprintf("%sToProto(data)", node.Name))),
			jen.Nil(),
		),
	)
}

func (s *serviceBuilder) funcList(file *jen.File, node *gen.Type) {
	defer file.Line()
	var fields = make([]*jen.Statement, 0, len(node.Fields))
	for _, item := range node.Fields {
		fields = append(
			fields,
			chain("request", "Options", text.ProtoPascal(item.Name), "Selector").
				Call(jen.Qual(path.Join(s.entPackage, strings.ToLower(node.Name)), fmt.Sprintf("Field%s", text.EntPascal(item.Name)))),
		)
	}

	var edges jen.Statement
	for _, edge := range node.Edges {

		var edgeFields []*jen.Statement
		for _, ef := range edge.Type.Fields {
			edgeFields = append(
				edgeFields,
				chain("e", text.ProtoPascal(ef.Name), "Selector").Call(jen.Qual(path.Join(s.entPackage, strings.ToLower(edge.Type.Name)), fmt.Sprintf("Field%s", text.EntPascal(ef.Name)))),
			)
		}
		edges = append(
			edges,
			// if e := request.Options.EdgeName; e != nil
			jen.If(define("e").Add(chain("request", "Options", text.EntPascal(edge.Name)).Op(";").Id("e").Op("!=").Nil())).Block(
				chain("query", fmt.Sprintf("With%s", text.EntPascal(edge.Name))).
					Call(jen.Func().Params(jen.Id("eq").Op("*").Id("ent").Dot(fmt.Sprintf("%sQuery", edge.Type.Name))).Block(
						chain("eq", "Where").Call(jen.Qual(pkgEntproto, "Selectors").Index(chain("predicate", edge.Type.Name)).Add(calls(edgeFields...)).Op("...")),
					)),
			),
		)
	}

	s.serviceFunc(file, node.Name, "List", fmt.Sprintf("List%sRequest", node.Name), fmt.Sprintf("List%sResponse", node.Name)).
		Block(
			jen.Id("conditions").Op(":=").Qual(pkgEntproto, "Selectors").Index(
				jen.Qual(path.Join(s.entPackage, "predicate"), node.Name),
			).Add(calls(fields...)),
			jen.Id("query").Op(":=").Id("s").Op(".").Id("client").Op(".").Id("Query").Call(),
			jen.Id("query").Op("=").Id("query").Op(".").Id("Where").Call(jen.Id("conditions").Op("...")),

			jen.Line(),
			&edges,
			jen.Line(),
			s.buildOrder(node.Name),
			s.buildPaginator(node.Name),
			s.buildListResponse(node.Name),
		)
}

func (s *serviceBuilder) funcListEdge(file *jen.File, node *gen.Type) {
	for _, edge := range node.Edges {
		opts, err := entproto.GetAPIOptions(edge.Annotations)
		if err != nil {
			return
		}
		if opts.DisableEdge {
			continue
		}
		var fields []*jen.Statement
		for _, v := range edge.Type.Fields {
			fields = append(fields, chain("request", "Options", text.ProtoPascal(v.Name), "Selector").Call(
				jen.Qual(path.Join(s.entPackage, strings.ToLower(edge.Type.Name)), fmt.Sprintf("Field%s", text.EntPascal(v.Name))),
			))
		}
		s.serviceFunc(file, node.Name,
			fmt.Sprintf("List%s", text.ProtoPascal(edge.Name)),
			fmt.Sprintf("List%s%sRequest", node.Name, text.ProtoPascal(edge.Name)),
			fmt.Sprintf("List%sResponse", edge.Type.Name),
		).Block(
			define("query").Id("s").Dot("client").Dot("Query").Call().Dot("Where").Call(
				jen.Qual(path.Join(s.entPackage, strings.ToLower(node.Name)), "ID").Call(jen.Int().Call(chain("request", "Id"))),
			).Dot(fmt.Sprintf("Query%s", text.EntPascal(edge.Name))).Call().Dot("Where").Call(
				jen.Qual(pkgEntproto, "Selectors").Index(jen.Qual(path.Join(s.entPackage, "predicate"), edge.Type.Name)).Add(calls(fields...)).Op("..."),
			),
			file.Line(),
			s.buildOrder(edge.Type.Name),
			s.buildPaginator(edge.Type.Name),
			s.buildListResponse(edge.Type.Name),
		)
	}
}

func (s *serviceBuilder) serviceFunc(file *jen.File, name, method, request, response string) *jen.Statement {
	return file.Func().Op("(").Id("s").Op("*").Id(fmt.Sprintf("%sService", name)).Op(")").Id(method).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("request").Op("*").Qual(s.protoPackage, request),
		).
		Call(
			jen.Op("*").Qual(s.protoPackage, response),
			jen.Id("error"),
		)
}

func (s *serviceBuilder) buildPaginator(name string) jen.Code {
	return jen.If(define("paginator").Id("request").Dot("GetPaginator").Call().Op(";").Id("paginator").Op("!=").Nil()).Block(
		jen.Switch(define("page").Id("paginator").Dot("GetPaginator").Call().Op(".").Call(jen.Type())).Block(
			jen.Case(jen.Op("*").Qual(pkgEntproto, "Paginator_Classical")),
			assign("query").Id("query").
				Dot("Order").Call(chain("page", "Classical", "OrderSelector").Call()).
				Dot("\nOffset").Call(jen.Int().Call(
				jen.Id("page").Dot("Classical").Dot("GetLimit()").Op("*").Call(jen.Id("page").Dot("Classical").Dot("GetPage()").Op("-").Id("1")),
			)).
				Dot("\nLimit").Call(jen.Int().Call(jen.Id("page").Dot("Classical").Dot("GetLimit()"))),
			jen.Case(jen.Op("*").Qual(pkgEntproto, "Paginator_Infinite")),
			assign("query").Id("query").
				Dot("Order").Call(jen.Qual(path.Join(s.entPackage, strings.ToLower(name)), "ByID").Call()).
				Dot("\nLimit").Call(jen.Int().Call(chain("page", "Infinite", "GetLimit()"))),
			jen.If(define("sequence").Id("page").Dot("Infinite").Dot("GetSequence()").Op(";").Id("sequence").Op(">").Id("0")).Block(
				assign("query").Id("query").Dot("Where").Call(jen.Qual(path.Join(s.entPackage, strings.ToLower(name)), "IDLT").Call(
					jen.Int().Call(chain("page", "Infinite", "GetSequence()")),
				)),
			),
		),
	).Line()
}

func (s *serviceBuilder) buildListResponse(name string) *jen.Statement {
	codes := jen.Statement{
		define("data", "err").Id("query").Dot("All").Call(jen.Id("ctx")),
		jen.Line(),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		),
		jen.Line(),
		jen.Return(
			structPtr(jen.Qual(s.protoPackage, fmt.Sprintf("List%sResponse", name)), jen.Id("Data"), jen.Id("Trans").Call(jen.Id("data"), jen.Id(fmt.Sprintf("%sToProto", name)))),
			jen.Nil()),
	}
	return &codes
}

func (s *serviceBuilder) orderMapping(file *jen.File, node *gen.Type) *jen.Statement {
	mapName := fmt.Sprintf("%sOrderFields", strcase.ToLowerCamel(node.Name))

	byID := fmt.Sprintf("%sOrder_%s_ORDER_BY_ID", node.Name, strcase.ToScreamingSnake(node.Name))
	var fields = []jen.Code{
		jen.Qual(s.protoPackage, byID).Op(":").Qual(path.Join(s.entPackage, strings.ToLower(node.Name)), "FieldID").Op(","),
	}
	for _, v := range node.Fields {
		opts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			panic(err)
		}
		if !opts.Orderable {
			continue
		}
		enumName := fmt.Sprintf("%sOrder_%s_ORDER_BY_%s", node.Name, strcase.ToScreamingSnake(node.Name), strcase.ToScreamingSnake(v.Name))
		fields = append(
			fields,
			jen.Qual(s.protoPackage, enumName).Op(":").Qual(path.Join(s.entPackage, strings.ToLower(node.Name)), fmt.Sprintf("Field%s", text.EntPascal(v.Name))).Op(","),
		)
	}

	return file.Var().Id(mapName).Op("=").
		Map(jen.Qual(s.protoPackage, fmt.Sprintf("%sOrder", node.Name))).String().Block(fields...)
}

func (s *serviceBuilder) buildOrder(name string) *jen.Statement {
	return jen.For(define("_", "order").Range().Add(chain("request", "GetOrders()"))).Block(
		jen.If(jen.Id("order").Op("==").Nil()).Block(jen.Continue()),
		jen.Var().Id("opts").Index().Qual(entsql, "OrderTermOption"),
		jen.If(chain("order", "GetDesc()")).Block(
			assign("opts").Append(jen.Id("opts"), jen.Qual(entsql, "OrderDesc()")),
		),
		assign("query").Id("query").Dot("Order").Call(
			jen.Qual(entsql, "OrderByField").Call(
				jen.Id(fmt.Sprintf("%sOrderFields", strcase.ToLowerCamel(name))).Index(chain("order", "GetBy()")),
				jen.Id("opts").Op("...")).Dot("ToFunc()"),
		),
	).Line()
}

func (s *serviceBuilder) funcCreate(file *jen.File, node *gen.Type) {
	defer file.Line()

	create := define("create").Add(chain("s", "client", "Create()"))
	for _, v := range node.Fields {
		if v.Name == "created_at" || v.Name == "updated_at" {
			continue
		}
		create = create.Op(".").Id("\n").Id(fmt.Sprintf("Set%s", text.EntPascal(v.Name))).Call(jen.Id("request").Dot(fmt.Sprintf("Get%s()", text.ProtoPascal(v.Name))))
	}
	//fn := s.serviceFunc(file, "Create", node.Name)
	s.serviceFunc(file, node.Name, "Create", fmt.Sprintf("Create%sRequest", node.Name), fmt.Sprintf("Create%sResponse", node.Name)).
		Block(
			create,
			define("data", "err").Id("create").Dot("Save").Call(jen.Id("ctx")),
			ifErr(),
			jen.Return(
				structPtr(jen.Qual(s.protoPackage, fmt.Sprintf("Create%sResponse", node.Name)), jen.Id("Id"), jen.Int64().Call(chain("data", "ID"))),
				jen.Nil(),
			),
		)
}

func (s *serviceBuilder) funcUpdate(file *jen.File, node *gen.Type) {
	defer file.Line()

	var fields []jen.Code
	for _, v := range node.Fields {
		if v.Name == "created_at" || v.Name == "updated_at" {
			continue
		}
		if v.Immutable {
			continue
		}
		fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			panic(err)
		}
		if fieldOpts.Immutable {
			continue
		}
		fields = append(
			fields,
			jen.If(jen.Id("request").Dot("GetUpdate()").Dot(text.ProtoPascal(v.Name)).Op("!=").Nil()).Block(
				jen.Id("update").Dot(fmt.Sprintf("Set%s", text.EntPascal(v.Name))).Call(
					jen.Id("request").Dot("GetUpdate()").Dot(fmt.Sprintf("Get%s()", text.ProtoPascal(v.Name))),
				),
			).Line(),
		)
	}
	s.serviceFunc(file, node.Name, "Update", fmt.Sprintf("Update%sRequest", node.Name), fmt.Sprintf("Update%sResponse", node.Name)).
		Block(
			define("update").Id("s").Dot("client").Dot("UpdateOneID").Call(jen.Int().Call(chain("request", "GetId()"))),
			jen.Add(fields...),
			define("_", "err").Id("update").Dot("Save").Call(jen.Id("ctx")),
			ifErr(),
			jen.Return(jen.Op("&").Qual(s.protoPackage, fmt.Sprintf("Update%sResponse", node.Name)).Block(), jen.Nil()),
		)
}

func (s *serviceBuilder) funcSet(file *jen.File, node *gen.Type) {
	for _, v := range node.Fields {
		if v.Name == "created_at" || v.Name == "updated_at" {
			continue
		}
		if v.Immutable {
			continue
		}
		fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			panic(err)
		}
		if fieldOpts.Immutable || !fieldOpts.Settable {
			continue
		}
		s.serviceFunc(file, node.Name,
			fmt.Sprintf("Set%s", text.ProtoPascal(v.Name)),
			fmt.Sprintf("Set%s%sRequest", node.Name, text.ProtoPascal(v.Name)),
			fmt.Sprintf("Set%s%sResponse", node.Name, text.ProtoPascal(v.Name)),
		).Block(
			define("_", "err").Id("s").Dot("client").
				Dot("UpdateOneID").Call(jen.Int().Call(chain("request", "GetId()"))).
				Dot(fmt.Sprintf("Set%s", text.EntPascal(v.Name))).Call(jen.Id("request").Dot(fmt.Sprintf("Get%s()", text.ProtoPascal(v.Name)))).
				Dot("Save").Call(jen.Id("ctx")),
			ifErr(),
			jen.Return(jen.Op("&").Qual(s.protoPackage, fmt.Sprintf("Set%s%sResponse", node.Name, text.ProtoPascal(v.Name))).Block(), jen.Nil()),
		).Line()

	}
}

func (s *serviceBuilder) funcDelete(file *jen.File, node *gen.Type) {
	defer file.Line()
	s.serviceFunc(file, node.Name, "Delete", fmt.Sprintf("Delete%sRequest", node.Name), fmt.Sprintf("Delete%sResponse", node.Name)).
		Block(
			define("err").Id("s").Dot("client").
				Dot("DeleteOneID").Call(jen.Int().Call(chain("request", "GetId()"))).
				Dot("Exec").Call(jen.Id("ctx")),
			ifErr(),
			jen.Return(jen.Op("&").Qual(s.protoPackage, fmt.Sprintf("Delete%sResponse", node.Name)).Block(), jen.Nil()),
		)
}

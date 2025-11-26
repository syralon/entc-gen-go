package entservice

import (
	"context"
	"fmt"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/dave/jennifer/jen"
	"github.com/syralon/entc-gen-go/internal/tools/text"
)

const (
	entproto    = "github.com/syralon/entc-gen-go/proto/syralon/entproto"
	timestamppb = "google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceBuilder struct {
	entPackage   string
	protoPackage string
}

type ServiceOption func(*ServiceBuilder)

func WithEntPackage(pkg string) ServiceOption {
	return func(builder *ServiceBuilder) {
		builder.entPackage = pkg
	}
}

func WithProtoPackage(pkg string) ServiceOption {
	return func(builder *ServiceBuilder) {
		builder.protoPackage = pkg
	}
}

func NewServiceBuilder(opts ...ServiceOption) *ServiceBuilder {
	sb := &ServiceBuilder{}
	for _, opt := range opts {
		opt(sb)
	}
	return sb
}

func (s *ServiceBuilder) Build(_ context.Context, node *gen.Type) (*jen.File, error) {
	file := jen.NewFile("service")

	//"google.golang.org/protobuf/types/known/timestamppb"
	//"github.com/syralon/entc-gen-go/proto/syralon/entproto"
	file.ImportAlias(s.protoPackage, "pb")
	s.funcToProto(file, node)
	s.funcFromProto(file, node)
	s.serviceStruct(file, node)
	s.funcGet(file, node)
	s.funcList(file, node)
	return file, nil
}

func (s *ServiceBuilder) funcToProto(file *jen.File, node *gen.Type) {
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

func (s *ServiceBuilder) funcFromProto(file *jen.File, node *gen.Type) {
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

func (s *ServiceBuilder) serviceStruct(file *jen.File, node *gen.Type) {
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

func (s *ServiceBuilder) funcGet(file *jen.File, node *gen.Type) {
	defer file.Line()

	fn := s.serviceFunc(file, "Get", node.Name)
	fn = fn.Block(
		// data, err := s.client.Get(ctx, int(request.GetId()))
		jen.List(jen.Id("data"), jen.Id("err")).Op(":=").Id("s").Op(".").Id("client").Op(".").Id("Get").
			Call(jen.Id("ctx"), jen.Int().Op("(").Id("request").Op(".").Id("GetId()").Op(")")),
		// if err != nil {}
		jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
			jen.Return(jen.Id("nil"), jen.Id("err")),
		),
		jen.Return(
			jen.Op("&").Qual(s.protoPackage, fmt.Sprintf("Get%sResponse", node.Name)).Block(
				jen.Id("Data").Op(":").Id(fmt.Sprintf("%sToProto(data)", node.Name)).Op(","),
			),
			jen.Id("nil"),
		),
	)
}

func (s *ServiceBuilder) funcList(file *jen.File, node *gen.Type) {
	var fields = make([]jen.Code, 0, len(node.Fields))
	for _, item := range node.Fields {
		fields = append(
			fields,
			chain("request", "Options", text.ProtoPascal(item.Name), "Selector").
				Call(jen.Qual(path.Join(s.entPackage, strings.ToLower(node.Name)), fmt.Sprintf("Field%s", text.EntPascal(item.Name)))),
		)
	}

	var edges jen.Statement
	for _, edge := range node.Edges {

		var edgeFields []jen.Code
		for _, ef := range edge.Type.Fields {
			edgeFields = append(
				edgeFields,
				chain("e", text.ProtoPascal(ef.Name), "Selector").Call(jen.Qual(path.Join(s.entPackage, strings.ToLower(edge.Type.Name)), fmt.Sprintf("Field%s", text.EntPascal(ef.Name)))),
			)
		}
		edges = append(
			edges,
			// if e := request.Options.EdgeName; e != nil
			jen.If(join(define("e"), chain("request", "Options", text.EntPascal(edge.Name)).Op(";").Id("e").Op("!=").Nil())).Block(
				chain("query", fmt.Sprintf("With%s", text.EntPascal(edge.Name))).
					Call(jen.Func().Params(jen.Id("eq").Op("*").Id("ent").Dot(fmt.Sprintf("%sQuery", edge.Type.Name))).Block(
						chain("eq", "Where").Call(chain("entproto", "Selectors").Index(chain("predicate", edge.Type.Name)).Call(edgeFields...).Op("...")),
					)),
			),
		)
	}

	s.serviceFunc(file, "List", node.Name).
		Block(
			jen.Id("conditions").Op(":=").Qual(entproto, "Selectors").Index(
				jen.Qual(path.Join(s.entPackage, "predicate"), node.Name),
			).Call(fields...),
			jen.Id("query").Op(":=").Id("s").Op(".").Id("client").Op(".").Id("Query").Call(),
			jen.Id("query").Op("=").Id("query").Op(".").Id("Where").Call(jen.Id("conditions").Op("...")),

			jen.Line(),
			&edges,
			jen.Line(),
			jen.If(define("paginator").Id("request").Dot("GetPaginator").Call().Op(";").Id("paginator").Op("!=").Nil()).Block(
				jen.Switch(define("page").Id("paginator").Dot("GetPaginator").Call().Op(".").Call(jen.Type())).Block(
					jen.Case(jen.Op("*").Qual(entproto, "Paginator_Classical")),
					assign("query").Id("query").
						Dot("Order").Call(chain("page", "Classical", "OrderSelector").Call()).
						Dot("Offset").Call(jen.Int().Call(
						jen.Id("page").Dot("Classical").Dot("GetLimit()").Op("*").Call(jen.Id("page").Dot("Classical").Dot("GetPage()").Op("-").Id("1")),
					)).
						Dot("Limit").Call(jen.Int().Call(jen.Id("page").Dot("Classical").Dot("GetLimit()"))),
					jen.Case(jen.Op("*").Qual(entproto, "Paginator_Infinite")),
					assign("query").Id("query").
						Dot("Order").Call(jen.Qual(path.Join(s.entPackage, strings.ToLower(node.Name)), "ByID").Call()).
						Dot("Limit").Call(jen.Int().Call(chain("page", "Infinite", "GetLimit()"))),
					jen.If(define("sequence").Id("page").Dot("Infinite").Dot("GetSequence()").Op(";").Id("sequence").Op(">").Id("0")).Block(
						assign("query").Id("query").Dot("Where").Call(jen.Qual(path.Join(s.entPackage, strings.ToLower(node.Name)), "IDLT").Call(
							jen.Int().Call(chain("page", "Infinite", "GetSequence()")),
						)),
					),
				),
			),
			jen.Line(),

			jen.List(jen.Id("data"), jen.Id("err")).Op(":=").Id("query").Op(".").Id("All").Call(jen.Id("ctx")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return(jen.Id("nil"), jen.Id("err")),
			),

			jen.Return(
				jen.Op("&").Qual(s.protoPackage, fmt.Sprintf("List%sResponse", node.Name)).Block(
					jen.Id("Data").Op(":").Id("Trans").Call(jen.Id("data"), jen.Id(fmt.Sprintf("%sToProto", node.Name))).Op(","),
				),
				jen.Id("nil"),
			),
		)

}

func (s *ServiceBuilder) serviceFunc(file *jen.File, method string, name string) *jen.Statement {
	return file.Func().Op("(").Id("s").Op("*").Id(fmt.Sprintf("%sService", name)).Op(")").Id(method).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("request").Op("*").Qual(s.protoPackage, fmt.Sprintf("%s%sRequest", method, name)),
		).
		Call(
			jen.Op("*").Qual(s.protoPackage, fmt.Sprintf("%s%sResponse", method, name)),
			jen.Id("error"),
		)
}

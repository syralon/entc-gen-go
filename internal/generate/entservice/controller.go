package entservice

import (
	"context"
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/syralon/entc-gen-go/internal/tools/text"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	"path"
	"strings"
)

type controllerBuilder struct {
	pkg        string
	entPackage string
}

func (s *controllerBuilder) Build(_ context.Context, node *gen.Type) (*jen.File, error) {
	file := jen.NewFile(fmt.Sprintf("%sservice", strings.ToLower(node.Name)))

	name := fmt.Sprintf("%sService", node.Name)

	file.Type().Id(name).Struct(
		jen.Op("*").Qual(path.Join(s.pkg, "internal/service"), name),
	).Line()

	file.Func().Id("New" + name).
		Params(jen.Id("client").Op("*").Qual(s.entPackage, "Client")).Op("*").Id(name).
		Block(
			jen.Return(jen.Op("&").Id(name).Block(
				jen.Id(name).Op(":").Qual(path.Join(s.pkg, "internal/service"), "New"+name).Call(jen.Id("client")).Op(","),
			)),
		)
	return file, nil
}

func customControllerMethod(protoPackage string, node *gen.Type) (map[string]*jen.File, error) {
	files := make(map[string]*jen.File)
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
		file := rpcMethod(
			protoPackage, node.Name,
			fmt.Sprintf("Set%s", text.ProtoPascal(v.Name)),
			fmt.Sprintf("Set%s%sRequest", node.Name, text.ProtoPascal(v.Name)),
			fmt.Sprintf("Set%s%sResponse", node.Name, text.ProtoPascal(v.Name)),
		)
		filename := fmt.Sprintf("internal/controller/%sservice/set%s.go", strings.ToLower(node.Name), strings.ToLower(v.Name))
		files[filename] = file
	}
	for _, v := range node.Edges {
		var method string
		if v.Unique {
			method = "Get"
		} else {
			method = "List"
		}
		file := rpcMethod(
			protoPackage, node.Name,
			fmt.Sprintf("%s%s", method, text.ProtoPascal(v.Name)),
			fmt.Sprintf("%s%s%sRequest", method, node.Name, text.ProtoPascal(v.Name)),
			fmt.Sprintf("%s%s%sResponse", method, node.Name, text.ProtoPascal(v.Name)),
		)
		filename := fmt.Sprintf("internal/controller/%sservice/%s%s.go", strings.ToLower(node.Name), strings.ToLower(method), strings.ToLower(v.Name))
		files[filename] = file
	}
	opt, err := entproto.GetAPIOptions(node.Annotations)
	if err != nil {
		return nil, err
	}
	for _, m := range opt.Method.Methods() {
		method := text.ProtoPascal(m.String())
		file := rpcMethod(
			protoPackage, node.Name,
			method,
			fmt.Sprintf("%s%sRequest", method, node.Name),
			fmt.Sprintf("%s%sResponse", method, node.Name),
		)
		filename := fmt.Sprintf("internal/controller/%sservice/%s.go", strings.ToLower(node.Name), strings.ToLower(method))
		files[filename] = file
	}
	return files, nil
}

func rpcMethod(protoPackage, name, method, request, response string) *jen.File {
	file := jen.NewFile(fmt.Sprintf("%sservice", strings.ToLower(name)))
	file.ImportAlias(protoPackage, "pb")
	file.Comment(fmt.Sprintf("uncomment codes in this file to rewrite auto generated '%s' method", method))
	file.Func().Params(jen.Id("s").Op("*").Id(fmt.Sprintf("%sService", name))).Id(method).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("request").Op("*").Qual(protoPackage, request),
		).
		Params(
			jen.Id("response").Op("*").Qual(protoPackage, response),
			jen.Id("err").Error(),
		).
		Block(
			jen.Return(jen.Id("s").Dot(fmt.Sprintf("%sService", name))).Dot(method).Call(jen.Id("ctx"), jen.Id("request")),
		)
	return file
}

func controllerProvider(pkg, entPackage, protoPackage string, graph *gen.Graph) (*jen.File, error) {
	file := jen.NewFile("controller")
	file.ImportAlias(protoPackage, "pb")
	f1 := make([]jen.Code, 0)
	f2 := make([]jen.Code, 0)
	for _, node := range graph.Nodes {
		name := fmt.Sprintf("%sService", node.Name)
		imp := fmt.Sprintf("%s/internal/controller/%sservice", pkg, strings.ToLower(node.Name))
		f1 = append(
			f1,
			jen.Id(fmt.Sprintf("%sService", node.Name)).Op("*").Qual(imp, name),
		)
		f2 = append(
			f2,
			jen.Id(name).Op(":").Qual(imp, fmt.Sprintf("New%sService", node.Name)).Call(jen.Id("client")).Op(","),
		)
	}
	file.Type().Id("Services").Struct(f1...).Line()
	file.Func().Id("NewServices").Params(jen.Id("client").Op("*").Qual(entPackage, "Client")).Op("*").Id("Services").Block(
		jen.Return(jen.Op("&").Id("Services").Block(f2...)),
	).Line()

	file.Var().Id("ProviderSet").Op("=").Qual(pkgWire, "NewSet").Add(calls(jen.Id("NewServices"))).Line()
	return file, nil
}

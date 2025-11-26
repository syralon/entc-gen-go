package main

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

func main() {
	fset := token.NewFileSet()

	// 创建文件
	file := &ast.File{
		Name: ast.NewIdent("service"),
		Decls: []ast.Decl{
			buildImportDecl(),
			buildUserFromEntFunc(),
			buildUserServiceStruct(),
			buildGetMethod(),
		},
	}

	printer.Fprint(os.Stdout, fset, file)
}

func buildImportDecl() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"context"`}},
			&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"github.com/example/example/ent"`}},
			&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"github.com/example/example/proto/example"`}},
			&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"google.golang.org/protobuf/types/known/timestamppb"`}},
		},
	}
}

func buildUserFromEntFunc() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("UserFromEnt"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("user")},
						Type:  &ast.StarExpr{X: selector("ent", "User")},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.StarExpr{X: selector("example", "User")}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: selector("example", "User"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("Id"),
										Value: &ast.CallExpr{Fun: ast.NewIdent("int64"), Args: []ast.Expr{selector("user", "ID")}},
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("Name"),
										Value: selector("user", "Name"),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("CreatedAt"),
										Value: &ast.CallExpr{Fun: selector("timestamppb", "New"), Args: []ast.Expr{selector("user", "CreatedAt")}},
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("UpdatedAt"),
										Value: &ast.CallExpr{Fun: selector("timestamppb", "New"), Args: []ast.Expr{selector("user", "UpdatedAt")}},
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("Status"),
										Value: selector("user", "Status"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildUserServiceStruct() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("UserService"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: selector("example", "UnimplementedUserServiceServer"),
							},
							{
								Names: []*ast.Ident{ast.NewIdent("client")},
								Type:  &ast.StarExpr{X: selector("ent", "UserClient")},
							},
						},
					},
				},
			},
		},
	}
}

func buildGetMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: ast.NewIdent("UserService")},
				},
			},
		},
		Name: ast.NewIdent("Get"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{Names: []*ast.Ident{ast.NewIdent("ctx")}, Type: selector("context", "Context")},
					{Names: []*ast.Ident{ast.NewIdent("req")}, Type: &ast.StarExpr{X: selector("example", "GetUserRequest")}},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.StarExpr{X: selector("example", "GetUserResponse")}},
					{Type: ast.NewIdent("error")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// user, err := s.client.Get(ctx, int(req.GetId()))
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("user"), ast.NewIdent("err")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: selector2(selector("s", "client"), "Get"),
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.CallExpr{Fun: ast.NewIdent("int"), Args: []ast.Expr{
									&ast.CallExpr{Fun: selector("req", "GetId")},
								}},
							},
						},
					},
				},
				// if err != nil { return nil, err }
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{ast.NewIdent("nil"), ast.NewIdent("err")},
							},
						},
					},
				},
				// return &example.GetUserResponse{Data: UserFromEnt(user)}, nil
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: selector("entpb", "GetUserResponse"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("Data"),
										Value: &ast.CallExpr{Fun: ast.NewIdent("UserFromEnt"), Args: []ast.Expr{ast.NewIdent("user")}},
									},
								},
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

// 辅助函数：构造 a.b
func selector(pkg, name string) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: ast.NewIdent(pkg), Sel: ast.NewIdent(name)}
}

// 辅助函数：构造 (obj).method
func selector2(x ast.Expr, method string) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: x, Sel: ast.NewIdent(method)}
}

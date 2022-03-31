package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alxbse/toiletpaper/pkg/steps/file"
	"github.com/alxbse/toiletpaper/pkg/utils"
)

func main() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		panic("no build info")
	}
	fmt.Printf("buildInfo: %#v\n", buildInfo)

	fset := token.NewFileSet()

	imports := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("path"),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("sync"),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("os/user"),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("github.com/alxbse/toiletpaper/pkg/steps/file"),
			},
		},
		&ast.ImportSpec{
			Name: ast.NewIdent("_"),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("embed"),
			},
		},
	}

	globals := []ast.Spec{}

	steps := []ast.Stmt{}

	filepath.WalkDir(".", func(p string, d fs.DirEntry, err error) error {
		fmt.Printf("path: %s, isDir: %t\n", p, d.IsDir())
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(p, ".git") {
			return nil
		}
		if p == "main.go" {
			return nil
		}
		if p == "go.mod" {
			return nil
		}

		text := fmt.Sprintf("//go:embed %s", strconv.Quote(p))
		ident := utils.SnakeCasePath(p)

		valueSpec := &ast.ValueSpec{
			Doc: &ast.CommentGroup{
				List: []*ast.Comment{
					&ast.Comment{
						Text: text,
					},
				},
			},
			Names: []*ast.Ident{
				ast.NewIdent(ident),
			},
			Type: &ast.ArrayType{
				Elt: ast.NewIdent("byte"),
			},
		}
		globals = append(globals, valueSpec)

		deferStmt := &ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("wg"),
					Sel: ast.NewIdent("Done"),
				},
			},
		}

		compositeLit := file.AstFromFile(p, true)
		assignStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("step"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				compositeLit,
			},
		}
		c := &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("step"),
						Sel: ast.NewIdent("Do"),
					},
					Args:     nil,
					Ellipsis: token.NoPos,
				},
			},
		}
		ifStmt := &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("err"),
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						&ast.CallExpr{
							Fun: ast.NewIdent("panic"),
							Args: []ast.Expr{
								ast.NewIdent("err"),
							},
						},
					},
				},
			},
		}
		g := ast.GoStmt{
			Call: &ast.CallExpr{
				Fun: &ast.FuncLit{
					Type: &ast.FuncType{},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{deferStmt, assignStmt, c, ifStmt},
					},
				},
			},
		}
		steps = append(steps, &g)

		return nil
	})

	stmts := []ast.Stmt{
		&ast.DeclStmt{
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("wg"),
						},
						Type: ast.NewIdent("sync.WaitGroup"),
					},
				},
			},
		},
		&ast.ExprStmt{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("wg"),
					Sel: ast.NewIdent("Add"),
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.INT,
						Value: strconv.Itoa(len(steps)),
					},
				},
			},
		},
		&ast.ExprStmt{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("user"),
					Sel: ast.NewIdent("Current"),
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("baseDir"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				ast.NewIdent("\"/tmp\""),
			},
		},
	}

	stmts = append(stmts, steps...)

	w := ast.ExprStmt{
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("wg"),
				Sel: ast.NewIdent("Wait"),
			},
		},
	}
	stmts = append(stmts, &w)

	i := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: imports,
	}

	mainFunc := &ast.FuncDecl{
		Name: ast.NewIdent("main"),
		Body: &ast.BlockStmt{
			List: stmts,
		},
		Doc:  nil,
		Recv: nil,
		Type: &ast.FuncType{},
	}

	g := &ast.GenDecl{
		Tok:   token.VAR,
		Specs: globals,
	}

	f := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{i, g, mainFunc},
	}
	ast.SortImports(fset, f)
	ast.Print(fset, f)

	sourceFile, err := os.Create("main.go")
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(sourceFile, os.Stdout)

	err = format.Node(multiWriter, fset, f)
	if err != nil {
		panic(err)
	}
}

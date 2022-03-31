package file

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"

	"github.com/alxbse/toiletpaper/pkg/types"
)

func TestFileStep(t *testing.T) {
	tt := []struct {
		name string
		step FileStep
	}{
		{
			name: "hello",
			step: FileStep{
				Contents:    []byte("hello"),
				Destination: "/tmp/test_filestep",
				Hash:        "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			},
		},
		{
			name: "test missing parent directory",
			step: FileStep{
				Contents:    []byte("test"),
				Destination: "/tmp/does_not_exist/test",
				Hash:        "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.step.Do()
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}
		})
	}
}

func TestAstFromStep(t *testing.T) {
	tt := []struct {
		name     string
		step     types.Step
		expected *ast.CompositeLit
	}{
		{
			name: "test empty step",
			step: types.Step{
				Type: "file",
			},
			expected: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "file",
					},
					Sel: &ast.Ident{
						Name: "FileStep",
					},
				},
				Elts: []ast.Expr{},
			},
		},
		{
			name: "test embed step",
			step: types.Step{
				Type:  "file",
				Embed: "/tmp/test",
			},
			expected: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "file",
					},
					Sel: &ast.Ident{
						Name: "FileStep",
					},
				},
				Elts: []ast.Expr{
					&ast.KeyValueExpr{
						Key: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "Contents",
						},
						Value: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "/tmp/test",
						},
					},
				},
			},
		},
		{
			name: "test destination step",
			step: types.Step{
				Type:        "file",
				Destination: "/root/test",
			},
			expected: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "file",
					},
					Sel: &ast.Ident{
						Name: "FileStep",
					},
				},
				Elts: []ast.Expr{
					&ast.KeyValueExpr{
						Key: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "Destination",
						},
						Value: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"/root/test\"",
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := AstFromStep(&tc.step)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

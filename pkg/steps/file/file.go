package file

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/alxbse/toiletpaper/pkg/utils"
)

type FileStep struct {
	Contents    []byte
	Destination string
	Hash        string
}

func (f *FileStep) Do() error {
	dir, _ := path.Split(f.Destination)
	err := os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(f.Destination, os.O_RDWR, 0)
	if os.IsNotExist(err) {
		file, err = os.Create(f.Destination)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	fmt.Printf("file: %#v\n", file)
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return err
	}
	d, err := hex.DecodeString(f.Hash)
	if err != nil {
		return err
	}
	if bytes.Compare(hash.Sum(nil), d) != 0 {
		fmt.Printf("actual: %s, expected: %s\n", hex.EncodeToString(hash.Sum(nil)), f.Hash)
		_, err := file.WriteAt(f.Contents, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func AstFromFile(filename string, embed bool) *ast.CompositeLit {
	elts := []ast.Expr{}
	if embed {
		e := ast.KeyValueExpr{
			Key: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "Contents",
			},
			Value: &ast.BasicLit{
				Kind:  token.STRING,
				Value: utils.SnakeCasePath(filename),
			},
		}
		elts = append(elts, &e)
	}
	if filename != "" {
		dotfile := fmt.Sprintf(".%s", filename)
		e := ast.KeyValueExpr{
			Key: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "Destination",
			},
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "path",
					},
					Sel: &ast.Ident{
						Name: "Join",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "baseDir",
					},
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(dotfile),
					},
				},
			},
		}
		elts = append(elts, &e)
	}
	compositeLit := ast.CompositeLit{
		Type: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "file",
			},
			Sel: &ast.Ident{
				Name: "FileStep",
			},
		},
		Elts: elts,
	}
	return &compositeLit
}

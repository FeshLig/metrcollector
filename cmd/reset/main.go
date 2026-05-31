package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type StructInfo struct {
	Name   string
	Fields []*ast.Field
}

func main() {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		structs := findResetStructs(pkg)

		if len(structs) == 0 {
			continue
		}

		if err := generatePackage(pkg, structs); err != nil {
			panic(err)
		}
	}
}

func findResetStructs(pkg *packages.Package) []StructInfo {
	var result []StructInfo

	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			if !hasResetMarker(genDecl) {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				st, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				result = append(result, StructInfo{
					Name:   typeSpec.Name.Name,
					Fields: st.Fields.List,
				})
			}
		}
	}

	return result
}

func hasResetMarker(decl *ast.GenDecl) bool {
	if decl.Doc == nil {
		return false
	}

	for _, c := range decl.Doc.List {
		if strings.Contains(c.Text, "generate:reset") {
			return true
		}
	}

	return false
}

func generatePackage(pkg *packages.Package, structs []StructInfo) error {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "package %s\n\n", pkg.Name)

	for _, s := range structs {
		generateStructReset(&buf, s)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	dir := filepath.Dir(pkg.GoFiles[0])

	return os.WriteFile(
		filepath.Join(dir, "reset.gen.go"),
		src,
		0644,
	)
}

func generateStructReset(buf *bytes.Buffer, s StructInfo) {
	receiver := strings.ToLower(string(s.Name[0]))

	fmt.Fprintf(buf, "func (%s *%s) Reset() {\n", receiver, s.Name)
	fmt.Fprintf(buf, "\tif %s == nil {\n", receiver)
	fmt.Fprintf(buf, "\t\treturn\n")
	fmt.Fprintf(buf, "\t}\n\n")

	for _, field := range s.Fields {
		for _, name := range field.Names {
			generateFieldReset(
				buf,
				receiver+"."+name.Name,
				field.Type,
			)
		}
	}

	fmt.Fprintf(buf, "}\n\n")
}

func generateFieldReset(
	buf *bytes.Buffer,
	fieldExpr string,
	typ ast.Expr,
) {
	switch t := typ.(type) {

	case *ast.Ident:
		switch t.Name {
		case "string":
			fmt.Fprintf(buf,
				"\t%s = \"\"\n",
				fieldExpr,
			)

		case "bool":
			fmt.Fprintf(buf,
				"\t%s = false\n",
				fieldExpr,
			)

		case "int",
			"int8",
			"int16",
			"int32",
			"int64",
			"uint",
			"uint8",
			"uint16",
			"uint32",
			"uint64",
			"uintptr",
			"float32",
			"float64":
			fmt.Fprintf(buf,
				"\t%s = 0\n",
				fieldExpr,
			)

		default:
			fmt.Fprintf(buf,
				"\tif r, ok := any(&%s).(interface{ Reset() }); ok {\n",
				fieldExpr,
			)
			fmt.Fprintf(buf,
				"\t\tr.Reset()\n",
			)
			fmt.Fprintf(buf,
				"\t}\n",
			)
		}

	case *ast.ArrayType:
		fmt.Fprintf(buf,
			"\tif %s != nil {\n",
			fieldExpr,
		)
		fmt.Fprintf(buf,
			"\t\t%s = %s[:0]\n",
			fieldExpr,
			fieldExpr,
		)
		fmt.Fprintf(buf,
			"\t}\n",
		)

	case *ast.MapType:
		fmt.Fprintf(buf,
			"\tclear(%s)\n",
			fieldExpr,
		)

	case *ast.StarExpr:
		fmt.Fprintf(buf,
			"\tif %s != nil {\n",
			fieldExpr,
		)

		generatePointerReset(
			buf,
			"*"+fieldExpr,
			t.X,
		)

		fmt.Fprintf(buf,
			"\t}\n",
		)
	}
}

func generatePointerReset(
	buf *bytes.Buffer,
	expr string,
	typ ast.Expr,
) {
	switch t := typ.(type) {

	case *ast.Ident:
		switch t.Name {
		case "string":
			fmt.Fprintf(buf,
				"\t\t%s = \"\"\n",
				expr,
			)

		case "bool":
			fmt.Fprintf(buf,
				"\t\t%s = false\n",
				expr,
			)

		case "int",
			"int8",
			"int16",
			"int32",
			"int64",
			"uint",
			"uint8",
			"uint16",
			"uint32",
			"uint64",
			"uintptr",
			"float32",
			"float64":
			fmt.Fprintf(buf,
				"\t\t%s = 0\n",
				expr,
			)

		default:
			fmt.Fprintf(buf,
				"\t\tif r, ok := any(%s).(interface{ Reset() }); ok {\n",
				expr,
			)
			fmt.Fprintf(buf,
				"\t\t\tr.Reset()\n",
			)
			fmt.Fprintf(buf,
				"\t\t}\n",
			)
		}

	case *ast.ArrayType:
		fmt.Fprintf(buf,
			"\t\t%s = %s[:0]\n",
			expr,
			expr,
		)

	case *ast.MapType:
		fmt.Fprintf(buf,
			"\t\tclear(%s)\n",
			expr,
		)

	case *ast.StarExpr:
		fmt.Fprintf(buf,
			"\t\tif %s != nil {\n",
			expr,
		)

		generatePointerReset(
			buf,
			"*"+expr,
			t.X,
		)

		fmt.Fprintf(buf,
			"\t\t}\n",
		)
	}
}

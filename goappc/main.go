package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"io/ioutil"
	"crypto/md5"
)

const assetsImport = `"github.com/sysr-q/assets"`

type assetsVisitor struct {
	importName string
	rewritten map[*ast.Ident][]byte // const ident -> contents
}

func (v *assetsVisitor) assetsCallExpr(node ast.Node) bool {
	n, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}

	if len(n.Args) != 1	{
		return false
	}

	fun, ok := n.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if x, ok := fun.X.(*ast.Ident); !ok || x.Name != v.importName {
		return false
	}

	if fun.Sel.Name != "MustRead" && fun.Sel.Name != "Read" {
		return false
	}

	return true
}

func (v *assetsVisitor) rewriteCallExpr(node ast.Expr) ast.Expr {
	// We assume this is only called on nodes which pass v.assetsCallExpr(node)
	n, ok := node.(*ast.CallExpr)
	if !ok {
		return node
	}

	if len(n.Args) != 1 {
		return node
	}

	filename := n.Args[0].(*ast.BasicLit).Value
	filename = filename[1:len(filename)-1] // Strip quotes

	hash := fmt.Sprintf("%x", md5.Sum([]byte(filename)))
	ident := ast.NewIdent("Asset_" + hash[:16]) // 16 bytes ought to be enough.

	// TODO(sysr-q): This is rubbish.
	for i := range v.rewritten {
		if i.Name != ident.Name {
			continue
		}
		return ident
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1) // RIP in peace.
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1) // RIP^2
	}

	v.rewritten[ident] = b
	return ident
}

func (v *assetsVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CallExpr:
		for i, arg := range n.Args {
			if !v.assetsCallExpr(arg) {
				continue
			}
			n.Args[i] = v.rewriteCallExpr(arg)
		}
	case *ast.AssignStmt:
		for i, r := range n.Rhs {
			if !v.assetsCallExpr(r) {
				continue
			}
			n.Rhs[i] = v.rewriteCallExpr(r)
		}
	}
	return v
}

// Insert takes an *ast.File, walks through the ast of the file, replacing
// "magic" calls to /assets.(Must)?Read?/, replacing them with references
// to a constant containing the contents of that file.
func (v *assetsVisitor) Insert(file *ast.File) *ast.File {
	for _, s := range file.Imports {
		if s.Path.Value != assetsImport {
			continue
		}

		if s.Name != nil {
			v.importName = s.Name.Name
		} else {
			v.importName = "assets"
		}
	}

	// assetsImport is never imported?
	if v.importName == "" {
		return file
	}

	ast.Walk(v, file)


	var decls []ast.Decl
	decls = append(decls, file.Decls[0])

	// I'm so sorry about this entire for loop. Blame go/ast
	for name, content := range v.rewritten {
		elts := make([]ast.Expr, 0)
		for _, b := range content {
			elts = append(elts, &ast.BasicLit{
				ValuePos: token.NoPos,
				Kind: token.INT,
				Value: fmt.Sprintf("%#.2x", b),
			})
		}

		val := &ast.ValueSpec{
			Names: []*ast.Ident{name},
			Values: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.ArrayType{
						Lbrack: token.NoPos,
						Elt: ast.NewIdent("byte"),
					},
					Lbrace: token.NoPos,
					Elts: elts,
					Rbrace: token.NoPos,
				},
			},
		}

		decl := &ast.GenDecl{
			Doc: nil,
			TokPos: token.NoPos,
			Tok: token.VAR,
			Lparen: token.NoPos,
			Specs: []ast.Spec{val},
			Rparen: token.NoPos,
		}

		decls = append(decls, decl)
	}

	// var _ = assets.Read (to silence unused imports)
	underscore := &ast.GenDecl{
		Doc: nil,
		TokPos: token.NoPos,
		Tok: token.VAR,
		Lparen: token.NoPos,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("_")},
				Values: []ast.Expr{
					&ast.SelectorExpr{
						X: ast.NewIdent(v.importName),
						Sel: ast.NewIdent("Read"),
					},
				},
			},
		},
		Rparen: token.NoPos,
	}

	decls = append(decls, underscore)
	decls = append(decls, file.Decls[1:]...)
	file.Decls = decls // put our new []ast.Decl in place

	return file
}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, os.Args[1], nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	av := assetsVisitor{
		rewritten: make(map[*ast.Ident][]byte),
	}

	f = av.Insert(f) // Where the magic happens.
	printer.Fprint(os.Stdout, fset, f)

/*
	pkgs, err := parser.ParseDir(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	for pkgName, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			
		}
	}
*/
}

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"github.com/dchest/uniuri"
)

const assetsImport = `"github.com/sysr-q/assets"`

type assetNode struct {
	Contents []byte
	Const *ast.Ident
	Filename string
}

type assetsVisitor struct {
	found map[ast.Node]assetNode // filename -> contents
	importName string
	looking bool
	lookingPos ast.Node
	pass bool // true = first, false = second
}

func (v *assetsVisitor) firstPass(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Ident:
		if n.Name == v.importName {
			v.looking = true
			// This is the node we found
			v.lookingPos = n
			break
		}
		if !v.looking || n.Name != "Read" && n.Name != "MustRead" {
			v.looking = false
			v.lookingPos = nil
			break
		}
	case *ast.BasicLit:
		if n.Kind != token.STRING || !v.looking {
			break
		}
		v.found[v.lookingPos] = assetNode{
			Contents: []byte(n.Value), // TODO(sysr-q): Actually add the contents of the file.
			Const: ast.NewIdent("Assets_" + uniuri.New()),
			Filename: n.Value,
		}
		fmt.Printf("found: %#v starting: %#v\n", n, v.lookingPos)
		v.looking = false
		v.lookingPos = nil
	}
	return v
}

func (v *assetsVisitor) secondPass(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	}
	return v
}

func (v *assetsVisitor) Visit(node ast.Node) ast.Visitor {
	if v.pass {
		return v.firstPass(node)
	}
	return v.secondPass(node)
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

	// First pass: identify ast nodes to replace
	v.pass = true
	ast.Walk(v, file)

	// Second pass: walk through, replacing nodes
	v.pass = false
	ast.Walk(v, file)

	fmt.Printf("%#v\n", v.found)
	return file
}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	av := assetsVisitor{
		found: make(map[ast.Node]assetNode),
		importName: "",
	}

	f = av.Insert(f) // Where the magic happens.
	//printer.Fprint(os.Stdout, fset, f)
	_ = printer.Fprint
	ast.Print(fset, f)

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

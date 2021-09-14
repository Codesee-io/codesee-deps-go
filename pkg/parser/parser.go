package parser

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/pkg/errors"
)

type ParsedDir struct {
	// FileSet is the token.FileSet that was used to parse the directory. This
	// is used to get filenames from an AST node.
	FileSet *token.FileSet
	// ModulePath is the name of the module that is defined in a go.mod file.
	// This is used to resolve imports within the same module.
	ModulePath string
	ModuleRoot string
	// Packages is the return value of parser.ParseDir, where the map key is the
	// package name and the map value is the AST of the whole package (which is
	// a directory in Go).
	Packages map[string]*ast.Package
}

type Parser struct {
	root  string
	cache map[string]*ParsedDir
}

func New(root string) *Parser {
	return &Parser{
		root:  root,
		cache: map[string]*ParsedDir{},
	}
}

func (p *Parser) Parse(dir string) (*ParsedDir, error) {
	// First, we check the cache to see if we've already parsed this file, and
	// if we have, return the cached version instead.
	if parsedDir, ok := p.cache[dir]; ok {
		return parsedDir, nil
	}

	modulePath, moduleRoot, err := recursiveModulePath(p.root, dir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		// If we encounter an error when parsing, then it's probably not a
		// valid Go file, so we just skip it.
		p.cache[dir] = nil
		return nil, nil
	}

	p.cache[dir] = &ParsedDir{
		FileSet:    fset,
		ModulePath: modulePath,
		ModuleRoot: moduleRoot,
		Packages:   pkgs,
	}
	return p.cache[dir], nil
}

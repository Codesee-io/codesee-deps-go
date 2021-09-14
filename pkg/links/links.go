package links

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/Codesee-io/codesee-deps-go/pkg/parser"
	"github.com/pkg/errors"
)

type Link struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// This type aliases are only used to make some maps a bit more readable. They
// aren't actually necessary to work correctly.
type (
	// PackagePath e.g. github.com/Codesee-io/codesee-deps-go/pkg/parser
	PackagePath string
	// PackageName e.g. parser
	PackageName string
	// Identifier e.g. New or ParsedDir
	Identifier string
	// Filename e.g. /root/codesee-deps-go/pkg/parser/parser.go
	Filename string
)

// DetermineLinks takes in a root directory and generates all the links between
// the Go files in this directory, relative from this root directory. The order
// of links is not guaranteed to be deterministic to make it faster. If you're
// asserting equality for the links (e.g. in a test), make sure you sort it
// before your assertion.
func DetermineLinks(root string) ([]Link, error) {
	absRoot, err := filepath.Abs(root)
	dirs, err := determineGoDirectories(absRoot)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	linksSet := map[string]struct{}{}
	links := []Link{}

	p := parser.New(absRoot)

	// We determine the links for a project by making 2 passes over the
	// directories.

	// This is a mapping from package path to a package name. This is needed to
	// help generate the reverse mapping (name to path) for a specific file.
	// More details about why we need to do this can be found in the second
	// pass.
	pkgPathToPkgName := map[PackagePath]PackageName{}
	// This is a mapping from package path and object name in scope to the
	// filename that it's defined in. So this will look something like this:
	// {
	//   "github.com/Codesee-io/codesee-deps-go/pkg/parser": {
	//     "New": "/root/codesee-deps-go/pkg/parser/parser.go"
	//   }
	// }
	identifierToFilename := map[PackagePath]map[Identifier]Filename{}

	// This first pass populates pkgPathToPkgName and identifierToFilename.
	for _, dir := range dirs {
		parsedDir, err := p.Parse(dir)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if parsedDir == nil {
			// This package wasn't able to be parsed correctly, so we just skip
			// it.
			continue
		}

		// In most cases, there's only one package per directly, though this
		// isn't guaranteed.
		for pkgName, pkg := range parsedDir.Packages {
			// Add this package in our mapping from package path to package
			// name.
			pkgPath := PackagePath(strings.Replace(dir, parsedDir.ModuleRoot, parsedDir.ModulePath, -1))
			pkgPathToPkgName[pkgPath] = PackageName(pkgName)

			for _, file := range pkg.Files {
				pos := parsedDir.FileSet.Position(file.Pos())
				filename := Filename(pos.Filename)

				// For each file, go through all the objects that are in the
				// global scope (e.g. types, functions, const and var
				// declarations, etc.) and add them to our mapping from
				// identifier to filename.
				for name := range file.Scope.Objects {
					if _, ok := identifierToFilename[pkgPath]; !ok {
						identifierToFilename[pkgPath] = map[Identifier]Filename{}
					}
					identifierToFilename[pkgPath][Identifier(name)] = filename
				}
			}
		}
	}

	// This is a mapping from filename to the package path and object name that
	// is being used in that file. This is how we'll know exactly what is being
	// used in the imported package. This is necessary to be able to map back to
	// where a specific object is defined. So this will look something like
	// this:
	// {
	//   "/root/codesee-deps-go/pkg/links/links.go": {
	//     "github.com/Codesee-io/codesee-deps-go/pkg/parser": {
	//       "New": {}
	//     }
	//   }
	// }
	filenameToIdentifierUsed := map[Filename]map[PackagePath]map[Identifier]struct{}{}

	// This second pass populates filenameToIdentifierUsed.
	for _, dir := range dirs {
		parsedDir, err := p.Parse(dir)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if parsedDir == nil {
			// This package wasn't able to be parsed correctly, so we just skip
			// it.
			continue
		}

		for _, pkg := range parsedDir.Packages {
			pkgPath := PackagePath(strings.Replace(dir, parsedDir.ModuleRoot, parsedDir.ModulePath, -1))

			for _, file := range pkg.Files {
				// This is on a per-file basis since each file can have an alias
				// for an import.
				pkgNameToPkgPath := map[PackageName]PackagePath{}

				pos := parsedDir.FileSet.Position(file.Pos())
				filename := Filename(pos.Filename)

				dotImports := []PackagePath{}

				// Go through all the imports in this file.
				for _, importSpec := range file.Imports {
					// The path value is wrapped in quotes, so we need to trim
					// them.
					importedPkgPath := strings.Trim(importSpec.Path.Value, "\"")

					if !strings.Contains(importedPkgPath, parsedDir.ModulePath) {
						continue
					}

					// The import spec's name is only defined if it's been
					// aliased to a different name, like this:
					// import (
					//   c "github.com/a/b"
					// )
					// In this example, the name would be "c". If there is no
					// alias, then name is nil.
					if importSpec.Name == nil {
						// If there isn't a custom package name, then we need to
						// use the package's assigned name. While this is
						// usually the final segment in the package path, this
						// isn't guaranteed. So that's why we need to use the
						// mapping from package path to package name that we
						// generated in the first pass to fill in the default
						// package name.
						if name, ok := pkgPathToPkgName[PackagePath(importedPkgPath)]; ok {
							// The imported package path is found in our
							// mapping, which means this is an internal import,
							// not an external dependency.
							pkgNameToPkgPath[name] = PackagePath(importedPkgPath)
						}
					} else {
						if importSpec.Name.String() == "." {
							// If the name is ".", then all of that package's
							// identifiers are accessible without needing to
							// qualify it with a package name. Here's an
							// example:
							// import (
							//   . "fmt"
							// )
							// With this, we can then use Println and Printf
							// instead of fmt.Println and fmt.Printf.
							dotImports = append(dotImports, PackagePath(importedPkgPath))
						} else {
							pkgNameToPkgPath[PackageName(importSpec.Name.String())] = PackagePath(importedPkgPath)
						}
					}
				}

				// These are all the packages that we should check for
				// unresolved identifiers. If an identifier is being used
				// without a package name, that means it's either defined in its
				// own package, or it was imported with a ".".
				packagePathsWithUnqualifiedIdentifiers := append([]PackagePath{pkgPath}, dotImports...)

				for _, ident := range file.Unresolved {
					for _, pkgPath := range packagePathsWithUnqualifiedIdentifiers {
						if toFilename, ok := identifierToFilename[pkgPath][Identifier(ident.String())]; ok {
							setKey := fmt.Sprintf("%s:%s", filename, toFilename)

							if _, ok := linksSet[setKey]; !ok {
								links = append(links, Link{
									From: strings.Replace(string(filename), absRoot+"/", "", -1),
									To:   strings.Replace(string(toFilename), absRoot+"/", "", -1),
								})
								linksSet[setKey] = struct{}{}
							}
						}
					}
				}

				// We've gotten all the intra-package resolutions, but we can't
				// rely on the Unresolved portion of the file AST for all of
				// them because it doesn't show the fully unresolved path e.g.
				// for parser.New, it will only tell us that parser is
				// unresolved, so we don't know what in parser was actually
				// used. To find those, we walk the AST to find all selector
				// expressions.
				ast.Inspect(file, func(n ast.Node) bool {
					// A selector expression is an expression in the format of
					// "X.Selector" (e.g. parser.New, p.Parse, parser.ParsedDir,
					// etc.). This is the main way that we'll determine how an
					// imported package is being used.
					selectorExpr, ok := n.(*ast.SelectorExpr)
					if !ok {
						return true
					}

					// X in our cause will be the package name.
					xIdent, ok := selectorExpr.X.(*ast.Ident)
					if !ok {
						return true
					}

					usedPkgName := PackageName(xIdent.String())
					// Sel in our cause is the identifier.
					usedIdentifier := Identifier(selectorExpr.Sel.String())

					usedPkgPath, ok := pkgNameToPkgPath[usedPkgName]
					if !ok {
						// This used package name is not found in our mapping,
						// which means this is not an internal import, but an
						// external dependency instead.
						return true
					}

					if _, ok := filenameToIdentifierUsed[filename]; !ok {
						filenameToIdentifierUsed[filename] = map[PackagePath]map[Identifier]struct{}{}
					}
					if _, ok := filenameToIdentifierUsed[filename][usedPkgPath]; !ok {
						filenameToIdentifierUsed[filename][usedPkgPath] = map[Identifier]struct{}{}
					}
					filenameToIdentifierUsed[filename][usedPkgPath][usedIdentifier] = struct{}{}

					return true
				})
			}
		}
	}

	// Now that we've pulled all the necessary data out of all the Go ASTs, we
	// can piece together a comprehensive list of links from file to file.
	for fromFilename, pkgPathsUsed := range filenameToIdentifierUsed {
		for pkgPath, identifiersUsed := range pkgPathsUsed {
			identifiersDefined, ok := identifierToFilename[pkgPath]
			if !ok {
				// We don't have the identifiers for this package path. This is
				// probably an external dependency.
				continue
			}

			for identifierUsed := range identifiersUsed {
				toFilename, ok := identifiersDefined[identifierUsed]
				if !ok {
					// We found an identifier being used by this package, but
					// that identifier isn't defined in this package. This could
					// be a Go file that wouldn't compile, or it could mean that
					// we missed adding it. Either way, we don't want it
					// interfering with all the other links, so we just skip it.
					continue
				}
				setKey := fmt.Sprintf("%s:%s", fromFilename, toFilename)

				if _, ok := linksSet[setKey]; !ok {
					links = append(links, Link{
						From: strings.Replace(string(fromFilename), absRoot+"/", "", -1),
						To:   strings.Replace(string(toFilename), absRoot+"/", "", -1),
					})
					linksSet[setKey] = struct{}{}
				}
			}
		}
	}

	return links, nil
}

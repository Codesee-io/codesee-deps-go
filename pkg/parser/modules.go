package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
)

// recursiveModulePath takes in the root of the project and a directory within
// that root, and it will search all directories starting with dir and ending
// with root to find a go.mod file. It returns the module path (which is
// retrieved from the go.mod), and the directory where the go.mod was found,
// which is the module root.
func recursiveModulePath(root, dir string) (string, string, error) {
	modFilePath := dir + "/go.mod"
	_, err := os.Stat(modFilePath)
	if err != nil && !os.IsNotExist(err) {
		return "", "", errors.WithStack(err)
	}

	if err == nil {
		// A go.mod file exists in this directory.
		mod, err := ioutil.ReadFile(modFilePath)
		if err != nil {
			return "", "", errors.WithStack(err)
		}
		return modfile.ModulePath(mod), dir, nil
	}

	if dir == root {
		// This means that we didn't find a go.mod file anywhere in the
		// directory tree, so this project might not be using Go modules.
		// Behavior without a go.mod is not fully tested. We could either throw
		// an error or just try to run it and see what happens. Sometimes it
		// does work (e.g. with the golang/go repo).
		return "", "", nil
	}

	// If we didn't find a go.mod in this directory, and we're not at the root
	// yet, go up one directory and look for a go.mod file there.
	return recursiveModulePath(root, filepath.Dir(dir))
}

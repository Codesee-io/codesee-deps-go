package links

import (
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

func determineGoDirectories(root string) ([]string, error) {
	dirSet := map[string]struct{}{}

	err := godirwalk.Walk(root, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			segments := strings.Split(path, "/")
			for _, segment := range segments {
				// Skip over the entire .git directory to speed up the walking.
				if segment == ".git" {
					return godirwalk.SkipThis
				}
				// Skip over any vendored dependencies since we don't care about
				// external dependencies.
				if segment == "vendor" {
					return godirwalk.SkipThis
				}
			}

			// We're looking for Go files, so if this is a directory, skip over
			// it.
			if de.IsDir() {
				return nil
			}
			// If this file isn't a Go file, skip it.
			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			// We've found a Go file, so we should add its directory to the set
			// of directories.
			dir := filepath.Dir(path)
			dirSet[dir] = struct{}{}

			return nil
		},
		Unsorted: true,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Convert the set to a slice.
	dirs := make([]string, 0, len(dirSet))
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}

	return dirs, nil
}

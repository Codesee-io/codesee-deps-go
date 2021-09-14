package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecursiveModulePath(t *testing.T) {
	t.Run("finds go.mod if it's in the directory", func(tt *testing.T) {
		root := "../testdata/simple-repo"
		dir := "../testdata/simple-repo"

		modulePath, moduleRoot, err := recursiveModulePath(root, dir)
		require.NoError(tt, err)

		assert.Equal(tt, "simple-repo", modulePath)
		assert.Equal(tt, "../testdata/simple-repo", moduleRoot)
	})

	t.Run("finds go.mod if it's in a parent directory", func(tt *testing.T) {
		root := "../testdata/simple-repo"
		dir := "../testdata/simple-repo/cmd/api"

		modulePath, moduleRoot, err := recursiveModulePath(root, dir)
		require.NoError(tt, err)

		assert.Equal(tt, "simple-repo", modulePath)
		assert.Equal(tt, "../testdata/simple-repo", moduleRoot)
	})

	t.Run("returns an empty string if there's no go.mod file", func(tt *testing.T) {
		root := "../testdata/simple-repo/cmd"
		dir := "../testdata/simple-repo/cmd/api"

		modulePath, moduleRoot, err := recursiveModulePath(root, dir)
		require.NoError(tt, err)

		assert.Equal(tt, "", modulePath)
		assert.Equal(tt, "", moduleRoot)
	})
}

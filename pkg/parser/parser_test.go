package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	t.Run("parses a Go file", func(tt *testing.T) {
		root := "../testdata/simple-repo"
		dir := "../testdata/simple-repo/cmd/api"
		p := New(root)

		parsedDir, err := p.Parse(dir)
		require.NoError(tt, err)

		require.NotNil(tt, parsedDir)
		assert.NotNil(tt, parsedDir.FileSet)
		assert.NotNil(tt, parsedDir.Packages)
		assert.Len(tt, parsedDir.Packages, 1)
		assert.NotNil(tt, parsedDir.Packages["main"])
	})

	t.Run("returns a cached version if we've parsed the file before", func(tt *testing.T) {
		root := "../testdata/simple-repo"
		dir := "../testdata/simple-repo/cmd/api"
		p := New(root)

		firstParsedDir, err := p.Parse(dir)
		require.NoError(tt, err)
		require.NotNil(tt, firstParsedDir)

		secondParsedDir, err := p.Parse(dir)
		require.NoError(tt, err)
		require.NotNil(tt, firstParsedDir)

		// This asserts that the pointers are the same.
		assert.Equal(tt, firstParsedDir, secondParsedDir)
	})

	t.Run("returns nil for an invalid file", func(tt *testing.T) {
		root := "../testdata/simple-repo"
		dir := "../testdata/simple-repo/pkg/invalid"
		p := New(root)

		parsedDir, err := p.Parse(dir)
		require.NoError(tt, err)

		assert.Nil(tt, parsedDir)
	})
}

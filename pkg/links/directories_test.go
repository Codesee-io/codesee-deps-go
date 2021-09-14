package links

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineGoDirectories(t *testing.T) {
	t.Run("handles a simple repo", func(tt *testing.T) {
		root := filepath.Clean("../testdata/simple-repo")

		dirs, err := determineGoDirectories(root)
		require.NoError(tt, err)

		// Sort the slice since its order isn't deterministic.
		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i] < dirs[j]
		})
		assert.Equal(tt, []string{
			"../testdata/simple-repo/cmd/api",
			"../testdata/simple-repo/pkg/handlers",
			"../testdata/simple-repo/pkg/invalid",
			"../testdata/simple-repo/pkg/server",
			"../testdata/simple-repo/pkg/signals",
		}, dirs)
	})
}

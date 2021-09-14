package links

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineLinks(t *testing.T) {
	t.Run("determines links for a simple repo", func(tt *testing.T) {
		root := "../testdata/simple-repo"

		links, err := DetermineLinks(root)
		require.NoError(tt, err)

		// Sort the slice since its order isn't deterministic.
		sort.Slice(links, func(i, j int) bool {
			if links[i].From == links[j].From {
				return links[i].To < links[j].To
			}
			return links[i].From < links[j].From
		})
		assert.Equal(tt, []Link{
			{From: "cmd/api/main.go", To: "pkg/server/server.go"},
			{From: "cmd/api/main.go", To: "pkg/signals/signals.go"},
			{From: "pkg/server/server.go", To: "pkg/handlers/handlers.go"},
			{From: "pkg/signals/signals_test.go", To: "pkg/signals/signals.go"},
		}, links)
	})
}

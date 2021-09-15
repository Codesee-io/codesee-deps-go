package links

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineLinks(t *testing.T) {
	t.Run("determines links for a simple repo", func(tt *testing.T) {
		root := "../testdata/simple-repo"

		links, err := DetermineLinks(root)
		require.NoError(tt, err)

		assert.Equal(tt, []Link{
			{From: "cmd/api/main.go", To: "pkg/server/server.go"},
			{From: "cmd/api/main.go", To: "pkg/signals/signals.go"},
			{From: "pkg/server/server.go", To: "pkg/handlers/handlers.go"},
			{From: "pkg/signals/signals_test.go", To: "pkg/signals/signals.go"},
		}, links)
	})
}

package keyboard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModesAreUniqueAndNamed(t *testing.T) {
	t.Parallel()
	seen := make(map[LightingMode]struct{})
	for _, mode := range LightingMode(0).All() {
		require.NotEmpty(t, mode.String())
		_, exists := seen[mode]
		require.False(t, exists)
		seen[mode] = struct{}{}
	}
	require.Len(t, seen, 20)
}

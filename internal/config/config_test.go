package config

import (
	"testing"

	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
	"github.com/stretchr/testify/require"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Setenv("AK820PRO_CONFIG_DIR", t.TempDir())
	cfg := &Config{
		LastMode: keyboard.ModeConfig{
			Mode:       keyboard.Static,
			R:          12,
			G:          34,
			B:          56,
			Rainbow:    true,
			Brightness: 78,
			Speed:      90,
			Direction:  keyboard.Right,
		},
		Themes: map[string][][3]uint8{
			"default": {
				{1, 2, 3},
				{4, 5, 6},
			},
		},
		LastThemeName: "default",
	}

	require.NoError(t, cfg.Save())
	loaded, err := Load()
	require.NoError(t, err)
	require.Equal(t, cfg, loaded)
}

func TestLoadMissingFileReturnsEmptyConfig(t *testing.T) {
	t.Setenv("AK820PRO_CONFIG_DIR", t.TempDir())
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Empty(t, cfg.LastThemeName)
	require.Len(t, cfg.Themes, 0)
}

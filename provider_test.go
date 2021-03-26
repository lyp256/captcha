package captcha

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDirImage(t *testing.T) {
	dir, err := filepath.Abs("geom/testdata")
	require.NoError(t, err)
	p, err := LoadDirImage(dir)
	require.NoError(t, err)
	img, err := p.Get()
	require.NoError(t, err)
	img.Bounds()
}

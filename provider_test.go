package captcha

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDirImage(t *testing.T) {
	ctx := context.Background()
	dir, err := filepath.Abs("geom/testdata")
	require.NoError(t, err)
	p, err := LoadDirImage(dir)
	require.NoError(t, err)
	imgName, err := p.Random(ctx)
	require.NoError(t, err)
	img, err := p.Get(ctx, imgName)
	img.Bounds()
}

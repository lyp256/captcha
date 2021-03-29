package kv

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapCURD(t *testing.T) {
	c := NewCURD(time.Millisecond * 10)
	ctx := context.Background()

	err := c.Set(ctx, "keyA", "imgA", 1)
	require.NoError(t, err)
	err = c.Set(ctx, "keyB", "imgA", 1)
	require.NoError(t, err)

	img, rad, err := c.Get(ctx, "keyA")
	require.NoError(t, err)
	require.Equal(t, 1., rad)
	require.Equal(t, "imgA", img)

	err = c.Set(ctx, "keyA", "imgA", 2)
	require.NoError(t, err)

	img, rad, err = c.Get(ctx, "keyA")
	require.NoError(t, err)
	require.Equal(t, 2., rad)

	time.Sleep(time.Millisecond * 11)

	_, _, err = c.Get(ctx, "keyA")
	require.Equal(t, ErrNotExist, err)

	_, _, err = c.Delete(ctx, "imgB")
	require.Equal(t, ErrNotExist, err)

	img, rad, err = c.Get(ctx, "keyB")
	require.Equal(t, ErrNotExist, err)

	err = c.Set(ctx, "keyA", "imgA", 2)
	require.NoError(t, err)

	img, rad, err = c.Delete(ctx, "keyA")
	require.NoError(t, err)
	require.Equal(t, 2., rad)

	_, _, err = c.Delete(ctx, "keyA")
	require.Equal(t, ErrNotExist, err)

}

func TestMapCacheClean(t *testing.T) {
	ctx := context.Background()
	c := NewCURD(time.Millisecond * 2).(*mapCache)
	_ = c.Set(ctx, "keyA", "imgA", 1)
	time.Sleep(time.Millisecond * 3)
	_ = c.Set(ctx, "keyA", "imgB", 1)
	assert.Len(t, c.data, 1)
}

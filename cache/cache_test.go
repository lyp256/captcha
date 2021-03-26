package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapCURD(t *testing.T) {
	c := NewCURD(time.Millisecond * 10)

	err := c.Set("a", 1)
	require.NoError(t, err)
	err = c.Set("b", 1)
	require.NoError(t, err)

	v, err := c.Get("a")
	require.NoError(t, err)
	require.Equal(t, 1., v)

	err = c.Set("a", 2)
	require.NoError(t, err)

	v, err = c.Get("a")
	require.NoError(t, err)
	require.Equal(t, 2., v)

	time.Sleep(time.Millisecond * 11)

	_, err = c.Get("a")
	require.Equal(t, ErrNotExist, err)

	_, err = c.Delete("b")
	require.Equal(t, ErrNotExist, err)

	v, err = c.Get("b")
	require.Equal(t, ErrNotExist, err)

	err = c.Set("a", 2)
	require.NoError(t, err)

	v, err = c.Delete("a")
	require.NoError(t, err)
	require.Equal(t, 2., v)

	_, err = c.Delete("a")
	require.Equal(t, ErrNotExist, err)

}

func TestMapCacheClean(t *testing.T) {
	c := NewCURD(time.Millisecond * 2).(*mapCache)
	_ = c.Set("a", 1)
	time.Sleep(time.Millisecond * 3)
	_ = c.Set("b", 1)
	assert.Len(t, c.data, 1)
}

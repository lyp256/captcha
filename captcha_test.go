package captcha

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_randRadian(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		rad := randRadian()
		require.True(t, rad >= 0.14)
		require.True(t, rad < 2*math.Pi-0.14)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

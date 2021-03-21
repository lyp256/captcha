package geom

import (
	"image"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCenterPoint(t *testing.T) {
	assert.Equal(t,
		image.Point{X: 5, Y: 5},
		CenterPoint(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 10, Y: 10}}),
	)
	assert.Equal(t,
		image.Point{X: 15, Y: 5},
		CenterPoint(image.Rectangle{Min: image.Point{X: 10, Y: 0}, Max: image.Point{X: 20, Y: 10}}),
	)

}

func TestInCircle(t *testing.T) {
	assert.True(t, InCircle(image.Point{X: 5, Y: 6}, image.Point{X: 5, Y: 5}, 2))
	assert.False(t, InCircle(image.Point{X: 5, Y: 7}, image.Point{X: 5, Y: 5}, 2))
}

func TestMaxSquare(t *testing.T) {
	assert.Equal(t,
		image.Rectangle{
			Min: image.Point{X: 0, Y: 50},
			Max: image.Point{X: 100, Y: 150},
		},
		MaxSquareRectangle(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: 100, Y: 200},
		}),
	)

	assert.Equal(t,
		image.Rectangle{
			Min: image.Point{X: 50, Y: 0},
			Max: image.Point{X: 150, Y: 100},
		},
		MaxSquareRectangle(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: 200, Y: 100},
		}),
	)

	assert.Equal(t,
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: 100, Y: 100},
		},
		MaxSquareRectangle(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: 100, Y: 100},
		}),
	)
}

func TestRotatePoint(t *testing.T) {
	assert.Equal(t,
		image.Point{X: 20, Y: 10},
		RotatePoint(image.Point{X: 0, Y: 10}, image.Point{X: 10, Y: 10}, math.Pi),
	)
}

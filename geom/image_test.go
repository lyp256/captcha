package geom

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/png"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/test.png
var testData []byte

func newRectangle() image.Image {
	return image.Rectangle{
		Min: image.Point{X: 100, Y: 200},
		Max: image.Point{X: 200, Y: 400},
	}
}

func TestPng(t *testing.T) {
	old, _ := png.Decode(bytes.NewReader(testData))
	sq := MaxSquareRectangle(old.Bounds())
	n := image.NewRGBA(sq)
	n2 := image.NewRGBA(sq)
	var center = CenterPoint(sq)
	var point image.Point
	for point.X = sq.Min.X; point.X < sq.Max.X; point.X++ {
		for point.Y = sq.Min.Y; point.Y < sq.Max.Y; point.Y++ {
			n.Set(point.X, point.Y, old.At(point.X, point.Y))
			if InCircle(point, center, sq.Dx()/2) {
				dst := RotatePoint(point, center, math.Pi)
				n2.Set(dst.X, dst.Y, old.At(point.X, point.Y))
			}
		}
	}

}

func TestSubSquare(t *testing.T) {
	{
		src, _ := png.Decode(bytes.NewReader(testData))
		dst := SubSquare(src)
		assert.Equal(t, dst.Bounds().Dx(), dst.Bounds().Dy())
	}
	{
		src := newRectangle()
		dst := SubSquare(src)
		assert.Equal(t, dst.Bounds().Dx(), dst.Bounds().Dy())
		var point image.Point
		bound := dst.Bounds()
		for point.X = bound.Min.X; point.X < bound.Max.X; point.X++ {
			for point.Y = bound.Min.Y; point.Y < bound.Max.Y; point.Y++ {
				assert.True(t, colorEqual(src.At(point.X, point.Y), dst.At(point.X, point.Y)))
			}
		}
	}
}

func TestCircle(t *testing.T) {
	{
		src, _ := png.Decode(bytes.NewReader(testData))
		dst := Circle(src)
		assert.Equal(t, dst.Bounds().Dx(), dst.Bounds().Dy())
		assert.True(t,
			colorEqual(color.Transparent, dst.At(src.Bounds().Min.X, src.Bounds().Min.Y)))
	}
	{
		src := newRectangle()
		dst := Circle(src)
		assert.Equal(t, dst.Bounds().Dx(), dst.Bounds().Dy())
		assert.True(t,
			colorEqual(color.Transparent, dst.At(src.Bounds().Min.X, src.Bounds().Min.Y)))
	}

}

func TestCircleRotate(t *testing.T) {
	src, _ := png.Decode(bytes.NewReader(testData))
	dst := CircleAndRotate(src, math.Pi)
	assert.Equal(t, dst.Bounds().Dx(), dst.Bounds().Dy())
	assert.True(t,
		colorEqual(color.Transparent, dst.At(src.Bounds().Min.X, src.Bounds().Min.Y)))
}

func colorEqual(a, b color.Color) bool {
	r1, g1, b2, a1 := a.RGBA()
	r2, g2, b1, a2 := b.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

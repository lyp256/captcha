package geom

import (
	"image"
	"image/color"
)

type subImager interface {
	SubImage(image.Rectangle) image.Image
}

// SubSquare 截取最大方形
func SubSquare(img image.Image) image.Image {
	subRge := MaxSquareRectangle(img.Bounds())
	subFace, ok := img.(subImager)
	if ok {
		return subFace.SubImage(subRge)
	}
	return commonSubImage(img, subRge)
}

// Circle 返回圆形内容，多余部分为透明,优先使用原 img 进行修改
func Circle(img image.Image) image.Image {
	rec := img.Bounds().Canon()
	center := CenterPoint(rec)
	R := MaxSquareRectangle(rec).Dx() / 2
	// 可以直接修改
	if v, ok := img.(imageSet); ok {
		cleanCirclePadding(v, rec, center, R)
		return SubSquare(img)
	}
	// 不能修改
	dst := image.NewNRGBA(rec)
	copyCircle(img, dst, rec, center, R)
	return SubSquare(img)
}

// CircleAndRotate Circle 保留最大圆型内容并旋转，不会修改原 img
func CircleAndRotate(img image.Image, rad float64) image.Image {
	// 有效矩形
	rec := MaxSquareRectangle(img.Bounds())
	center := CenterPoint(rec)
	R := rec.Dx() / 2
	dst := image.NewNRGBA(rec)
	var point image.Point
	for point.X = rec.Min.X; point.X < rec.Max.X; point.X++ {
		for point.Y = rec.Min.Y; point.Y < rec.Max.Y; point.Y++ {
			if InCircle(point, center, R) {
				// 使用反向映射 解决噪点问题
				rPoint := RotatePoint(point, center, rad*-1)
				dst.Set(point.X, point.Y, img.At(rPoint.X, rPoint.Y))
			}
		}
	}
	return dst
}

func cleanCirclePadding(dst imageSet, rec image.Rectangle, center image.Point, R int) {
	var point image.Point
	for point.X = rec.Min.X; point.X < rec.Max.X; point.X++ {
		for point.Y = rec.Min.Y; point.Y < rec.Max.Y; point.Y++ {
			if !InCircle(point, center, R) {
				dst.Set(point.X, point.Y, image.Transparent)
			}
		}
	}
}

func copyCircle(src imageAt, dst imageSet, rec image.Rectangle, center image.Point, R int) {
	var point image.Point
	for point.X = rec.Min.X; point.X < rec.Max.X; point.X++ {
		for point.Y = rec.Min.Y; point.Y < rec.Max.Y; point.Y++ {
			if InCircle(point, center, R) {
				copyImagePoint(src, dst, point)
			}
		}
	}
}

func commonSubImage(img image.Image, rec image.Rectangle) image.Image {
	subImg := image.NewNRGBA(rec)
	var point image.Point
	for point.X = rec.Min.X; point.X < rec.Max.X; point.X++ {
		for point.Y = rec.Min.Y; point.Y < rec.Max.Y; point.Y++ {
			copyImagePoint(img, subImg, point)
		}
	}
	return subImg
}

type imageAt interface {
	At(x, y int) color.Color
}

type imageSet interface {
	Set(x, y int, c color.Color)
}

func copyImagePoint(src imageAt, dst imageSet, point image.Point) {
	dst.Set(point.X, point.Y, src.At(point.X, point.Y))
}

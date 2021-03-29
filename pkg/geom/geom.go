package geom

import (
	"image"
	"math"
)

// InCircle 判断一个点是否在圆内
func InCircle(point, center image.Point, r int) bool {
	x := point.X - center.X
	y := point.Y - center.Y
	return x*x+y*y < r*r
}

// MaxSquareRectangle 返回矩形中心的最大正方形
func MaxSquareRectangle(r image.Rectangle) image.Rectangle {
	r = r.Canon()
	dx, dy := r.Dx(), r.Dy()
	if dx == dy {
		return r
	}
	if dx > dy {
		offset := (dx - dy) / 2
		r.Min.X += offset
		r.Max.X = r.Min.X + dy
		return r
	}
	offset := (dy - dx) / 2
	r.Min.Y += offset
	r.Max.Y = r.Min.Y + dx
	return r
}

// CenterPoint 返回矩形的中心点
func CenterPoint(r image.Rectangle) image.Point {
	return image.Point{
		X: (r.Min.X + r.Max.X) / 2,
		Y: (r.Min.Y + r.Max.Y) / 2,
	}
}

// RotatePoint 计算 src 绕着 ref 旋转 radian 后的位置
func RotatePoint(src, ref image.Point, radian float64) (dst image.Point) {
	sin, cos := math.Sincos(radian)
	return image.Point{
		X: int(math.Round(float64(src.X-ref.X)*cos-float64(src.Y-ref.Y)*sin)) + ref.X,
		Y: int(math.Round(float64(src.X-ref.X)*sin+float64(src.Y-ref.Y)*cos)) + ref.Y,
	}
}

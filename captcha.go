package captcha

import (
	"image"
	"math"
	"math/rand"

	"github.com/lyp256/captcha/cache"
	"github.com/lyp256/captcha/geom"
)

// Captcha 验证码实例
type Captcha struct {
	p Provider
	c cache.CURD
}

// NewCaptcha 创建验证码实例
func NewCaptcha(p Provider, c cache.CURD) *Captcha {
	return &Captcha{
		p: p,
		c: c,
	}
}

// Rand 返回一个随机旋转的图片验证码
func (i *Captcha) Rand(key string) (image.Image, float64, error) {
	rad := randRadian()
	img, err := i.Draw(key, rad)
	if err != nil {
		return nil, 0, err
	}
	return img, rad, nil
}

// Draw 返回一个指定旋转角度的验证码
func (i *Captcha) Draw(key string, rad float64) (image.Image, error) {
	img, err := i.p.Get()
	if err != nil {
		return nil, err
	}
	err = i.c.Set(key, rad)
	if err != nil {
		return nil, err
	}
	// 逆时针旋转
	geom.CircleAndRotate(img, rad*-1)
	return img, nil
}

// Compare 比较角度
func (i *Captcha) Compare(key string, rad float64) (threshold float64, err error) {
	src, err := i.c.Get(key)
	if err != nil {
		return -1, err
	}
	return math.Abs(math.Abs(rad) - math.Abs(src)), nil
}

// return [0.14,2π-0.14]
func randRadian() float64 {
	const r = 2*math.Pi - 0.28
	return 0.14 + (rand.Float64() * r)
}

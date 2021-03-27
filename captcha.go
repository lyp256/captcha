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

// Generate 生成一个验证请求
func (i *Captcha) Generate(key string) error {
	rad := randRadian()
	imageName, err := i.p.Random()
	if err != nil {
		return err
	}
	return i.c.Set(key, imageName, rad)
}

// DrawCaptcha 生成一个验证码图片
func (i *Captcha) DrawCaptcha(key string) (image.Image, error) {
	imageName, rad, err := i.c.Get(key)
	if err != nil {
		return nil, err
	}
	img, err := i.Draw(imageName, rad)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Draw 根据 imageKey 绘制一个验证码
func (i *Captcha) Draw(imageKey string, rad float64) (image.Image, error) {
	img, err := i.p.Get(imageKey)
	if err != nil {
		return nil, err
	}
	// 逆时针旋转
	img = geom.CircleAndRotate(img, rad)
	return img, nil
}

// Compare 比较角度
func (i *Captcha) Compare(key string, rad float64) (diff float64, err error) {
	_, src, err := i.c.Get(key)
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

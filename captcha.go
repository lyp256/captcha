package captcha

import (
	"context"
	"image"
	"math"
	"math/rand"
	"runtime"

	"github.com/lyp256/captcha/pkg/geom"
	cache2 "github.com/lyp256/captcha/pkg/kv"
)

// Captcha 验证码实例
type Captcha struct {
	img  Provider
	curd cache2.CURD
	// 用于控制做图的的并发数，防止大量请求占用内存
	tokenBucket chan struct{}
}

// NewCaptcha 创建验证码实例
func NewCaptcha(p Provider, c cache2.CURD) *Captcha {
	return &Captcha{
		img:         p,
		curd:        c,
		tokenBucket: make(chan struct{}, runtime.NumCPU()),
	}
}

// Generate 生成一个验证请求
func (i *Captcha) Generate(ctx context.Context, key string) error {
	rad := randRadian()
	imageName, err := i.img.Random(ctx)
	if err != nil {
		return err
	}
	return i.curd.Set(ctx, key, imageName, rad)
}

// DrawCaptcha 生成一个验证码图片
func (i *Captcha) DrawCaptcha(ctx context.Context, key string) (image.Image, error) {
	imageName, rad, err := i.curd.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	img, err := i.Draw(ctx, imageName, rad)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Draw 根据 imageKey 绘制一个验证码
func (i *Captcha) Draw(ctx context.Context, imageKey string, rad float64) (image.Image, error) {
	// 并发控制希望获取一个令牌
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case i.tokenBucket <- struct{}{}:
	}
	defer func() { <-i.tokenBucket }()

	img, err := i.img.Get(ctx, imageKey)
	if err != nil {
		return nil, err
	}
	// 逆时针旋转
	img = geom.CircleAndRotate(img, rad)
	return img, nil
}

// Compare 比较角度
func (i *Captcha) Compare(ctx context.Context, key string, rad float64) (diff float64, err error) {
	_, src, err := i.curd.Get(ctx, key)
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

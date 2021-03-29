package captcha

import (
	"context"
	"image"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"path/filepath"
)

// Provider 图片提供接口
type Provider interface {
	Get(ctx context.Context, key string) (image.Image, error)
	Random(ctx context.Context) (string, error)
}

// LoadDirImage 加载目录中的图片作为 Provider
func LoadDirImage(dirPath string) (Provider, error) {
	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var p = fsProvider{
		dir:   dirPath,
		items: make([]fs.DirEntry, 0, len(dir)),
	}
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		p.items = append(p.items, file)
	}
	return &p, err
}

type fsProvider struct {
	dir   string
	items []fs.DirEntry
}

// Get 获取一个已知图片
func (f *fsProvider) Get(ctx context.Context, fileName string) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

// Random 随机获取一个图片
func (f *fsProvider) Random(_ context.Context) (string, error) {
	index := rand.Intn(len(f.items))
	item := f.items[index]
	fileName := path.Join(f.dir, item.Name())
	return fileName, nil
}

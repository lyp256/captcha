package captcha

import (
	"image"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"path/filepath"
)

// Provider 图片提供接口
type Provider interface {
	Get() (image.Image, error)
}

type fsProvider struct {
	dir   string
	items []fs.DirEntry
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

func (f *fsProvider) Get() (image.Image, error) {
	index := rand.Intn(len(f.items))
	item := f.items[index]
	file, err := os.Open(path.Join(f.dir, item.Name()))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

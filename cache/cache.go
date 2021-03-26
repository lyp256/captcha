package cache

import (
	"errors"
	"sync"
	"time"
)

// ErrNotExist 不存在
var ErrNotExist = errors.New("not exist")

// CURD cache 的基本操作定义
type CURD interface {
	Set(key string, rad float64) error
	Get(Key string) (rad float64, err error)
	Delete(Key string) (rad float64, err error)
}

type item struct {
	val      float64
	createAt int64
}

type mapCache struct {
	data map[string]item
	mux  sync.Mutex
	ttl  int64 // Duration

	lastClean int64
}

// NewCURD 创建一个 基于 map 的 CURD
func NewCURD(TTL time.Duration) CURD {
	return &mapCache{
		data: make(map[string]item),
		ttl:  int64(TTL),
	}
}

func (m *mapCache) Set(key string, rad float64) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	unixNano := time.Now().UnixNano()
	m.data[key] = item{
		val:      rad,
		createAt: unixNano,
	}
	if unixNano-m.lastClean > (m.ttl / 10) {
		m.lastClean = unixNano
		m.clean(1000)
	}
	return nil
}

func (m *mapCache) clean(n int) {
	if m.ttl < 1 {
		return
	}
	for key, val := range m.data {
		if time.Now().UnixNano()-val.createAt > m.ttl {
			delete(m.data, key)
		}
		if n--; n < 0 {
			break
		}
	}
}

func (m *mapCache) Get(Key string) (rad float64, err error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	v, ok := m.data[Key]
	if !ok {
		return 0, ErrNotExist
	}
	if m.ttl > 0 && time.Now().UnixNano()-v.createAt > m.ttl {
		delete(m.data, Key)
		return 0, ErrNotExist
	}
	return v.val, nil
}

func (m *mapCache) Delete(Key string) (rad float64, err error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	v, ok := m.data[Key]
	if !ok {
		return 0, ErrNotExist
	}
	delete(m.data, Key)

	if m.ttl > 0 && time.Now().UnixNano()-v.createAt > m.ttl {
		return 0, ErrNotExist
	}
	return v.val, nil
}

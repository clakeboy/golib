package components

import (
	"ck_go_lib/utils"
	"errors"
	"fmt"
	"sync"
	"time"
)
//GO本地内存数据缓存
type MemCache struct {
	CacheIn
	lock    sync.RWMutex
	content map[string]*MemCacheData
}

type MemCacheData struct {
	data   interface{}
	expire int64
}

//数据是否过期
func (mc *MemCacheData) IsExpire() bool {
	if mc.expire == -1 {
		return false
	}
	return time.Now().Unix() > mc.expire
}

//新创建一个内存缓存
func NewMemCache() *MemCache {
	return &MemCache{
		lock:    sync.RWMutex{},
		content: make(map[string]*MemCacheData),
	}
}

/**
 * MemCache Functions
 */
func (m *MemCache) Get(key string) (interface{}, error) {
	m.lock.RLock()
	cache, ok := m.content[key]
	m.lock.RUnlock()
	if ok {
		if cache.IsExpire() {
			m.Delete(key)
			return "", errors.New("this key is expire!")
		}

		return cache.data, nil
	}
	return "", errors.New("not this key")
}

//设置数据
func (m *MemCache) Set(key string, v interface{}, expire int64) bool {
	m.lock.RLock()
	cache, ok := m.content[key]
	m.lock.RUnlock()
	if ok {
		cache.data = v
		cache.expire = utils.YN(expire == -1, int64(-1), time.Now().Unix()+expire).(int64)
		return true
	} else {
		cache = &MemCacheData{
			data:   v,
			expire: utils.YN(expire == -1, int64(-1), time.Now().Unix()+expire).(int64),
		}
		m.lock.Lock()
		m.content[key] = cache
		m.lock.Unlock()
		return true
	}
}

//删除缓存
func (m *MemCache) Delete(key string) bool {
	m.lock.Lock()
	delete(m.content, key)
	m.lock.Unlock()
	return true
}

func (m *MemCache) Dump() {
	for k, v := range m.content {
		fmt.Printf("Key: %s, Value: %v\n", k, v.data)
	}

	fmt.Println("Numbers: ", len(m.content))
}

func (m *MemCache) Stat() {
	fmt.Println("Numbers: ", len(m.content))
	fmt.Println()
}

package components

import (
	"sync"
	"errors"
	"time"
	"fmt"
)

type MemCache struct {
	lock sync.RWMutex
	content map[string]*MemCacheData
}

type MemCacheData struct {
	data interface{}
	expire int64
}
//数据是否过期
func (mc *MemCacheData) IsExpire() bool {
	return time.Now().Unix() > mc.expire
}
//新创建一个内存缓存
func NewMemCache() *MemCache {
	return &MemCache{
		lock:sync.RWMutex{},
		content:make(map[string]*MemCacheData),
	}
}

/**
 * MemCache Functions
 */
func (m *MemCache) Get(key string) (interface{},error) {
	m.lock.RLock()
	cache,ok := m.content[key]
	m.lock.RUnlock()
	if ok {
		if cache.IsExpire() {
			m.Delete(key)
			return "",errors.New("this key is expire!")
		}

		return cache.data,nil
	}
	return "",errors.New("not this key")
}
//设置数据
func (m *MemCache) Set(key string,v interface{},expire int64) bool {
	m.lock.RLock()
	cache,ok := m.content[key]
	m.lock.RUnlock()
	if ok {
		cache.data = v
		cache.expire = time.Now().Unix() + expire
		return true
	} else {
		cache = &MemCacheData{
			data:v,
			expire:time.Now().Unix()+expire,
		}
		m.lock.Lock()
		m.content[key] = cache
		m.lock.Unlock()
		return true
	}
}
//删除缓存
func (m *MemCache) Delete(key string) {
	m.lock.Lock()
	delete(m.content,key)
	m.lock.Unlock()
}

func (m *MemCache) Dump() {
	for k,v := range m.content {
		fmt.Printf("Key: %s, Value: %v\n",k,v.data)
	}

	fmt.Println("Numbers: ",len(m.content))
}

func (m *MemCache) Stat() {
	fmt.Println("Numbers: ",len(m.content))
}
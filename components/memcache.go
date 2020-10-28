package components

import (
	"errors"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"regexp"
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

//得到所有key值,按指定条件
// * 为得到所有

func (m *MemCache) Keys(pattern string) ([]string, error) {
	if pattern == "" {
		return nil, fmt.Errorf("param is empty")
	}
	var keys []string
	for k, _ := range m.content {
		if m.condition(pattern, k) {
			keys = append(keys, k)
		}
	}

	return keys, nil
}

func (MemCache) condition(condition string, key string) bool {
	if condition == "*" {
		return true
	}
	var reg *regexp.Regexp
	if condition[:1] == "*" {
		reg = regexp.MustCompile(fmt.Sprintf("%s$", condition[1:]))
	} else if condition[len(condition)-1:] == "*" {
		reg = regexp.MustCompile(fmt.Sprintf("^%s", condition[:len(condition)-1]))
	} else {
		reg = regexp.MustCompile(condition)
	}

	return reg.MatchString(key)
}

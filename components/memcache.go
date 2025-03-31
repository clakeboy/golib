package components

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/clakeboy/golib/utils"
)

// GO本地内存数据缓存
type MemCache struct {
	CacheIn
	content     sync.Map
	tmpName     string
	storeCancel chan bool
}

type MemCacheData struct {
	Data   interface{}
	Expire int64
}

// 数据是否过期
func (mc *MemCacheData) IsExpire() bool {
	if mc.Expire == -1 {
		return false
	}
	return time.Now().Unix() > mc.Expire
}

// 新创建一个内存缓存
func NewMemCache() *MemCache {
	return &MemCache{
		content:     sync.Map{},
		storeCancel: make(chan bool, 1),
	}
}

/**
 * MemCache Functions
 */
func (m *MemCache) Get(key string) (interface{}, error) {
	data, ok := m.content.Load(key)
	if ok {
		cache := data.(*MemCacheData)
		if cache.IsExpire() {
			m.Delete(key)
			return nil, fmt.Errorf("this key is expire")
		}

		return cache.Data, nil
	}
	return nil, errors.New("not this key")
}

// 设置数据
func (m *MemCache) Set(key string, v interface{}, expire int64) bool {
	data, ok := m.content.Load(key)
	if ok {
		cache := data.(*MemCacheData)
		cache.Data = v
		cache.Expire = utils.YN(expire == -1, int64(-1), time.Now().Unix()+expire).(int64)
		return true
	} else {
		cache := &MemCacheData{
			Data:   v,
			Expire: utils.YN(expire == -1, int64(-1), time.Now().Unix()+expire).(int64),
		}
		m.content.Store(key, cache)
		return true
	}
}

// 删除缓存
func (m *MemCache) Delete(key string) bool {
	m.content.Delete(key)
	return true
}

func (m *MemCache) Dump() {
	lens := 0
	m.content.Range(func(key, value any) bool {
		cache := value.(*MemCacheData)
		fmt.Printf("Key: %v, Value: %v\n", key, cache.Data)
		lens++
		return true
	})

	fmt.Println("Numbers: ", lens)
}

func (m *MemCache) Stat() {

}

// 加载本地缓存数据
func (m *MemCache) LoadLocal(name string) error {
	defer func() {
		m.startStore()
	}()
	m.tmpName = name

	fullPath := path.Join("./golib_cache", m.tmpName)
	if !utils.Exist(path.Dir(fullPath)) {
		os.MkdirAll(path.Dir(fullPath), 0755)
		return nil
	}
	if !utils.Exist(fullPath) {
		return fmt.Errorf("%s does not exist", fullPath)
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}
	list := make(map[string]*MemCacheData)
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	err = dec.Decode(&list)
	if err != nil {
		return err
	}
	for k, v := range list {
		m.content.Store(k, v)
	}
	return nil
}

// 开始自动运行本地文件写入
func (m *MemCache) startStore() {
	go func(m *MemCache) {
		for {
			select {
			case <-m.storeCancel:
				return
			default:
				err := m.Store()
				if err != nil {
					fmt.Println(err)
				}
				time.Sleep(5 * time.Second)
			}
		}
	}(m)
}

func (m *MemCache) StopStore() {
	m.storeCancel <- true
}

func (m *MemCache) Store() error {
	list := make(map[string]*MemCacheData)

	m.content.Range(func(key, value any) bool {
		if !value.(*MemCacheData).IsExpire() {
			list[key.(string)] = value.(*MemCacheData)
		}
		return true
	})

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(list)
	if err != nil {
		return err
	}

	fullPath := path.Join("./golib_cache", m.tmpName)
	err = os.WriteFile(fullPath, b.Bytes(), 0755)
	if err != nil {
		return err
	}
	return nil
}

// 得到所有key值,按指定条件
// * 为得到所有
func (m *MemCache) Keys(pattern string) ([]string, error) {
	if pattern == "" {
		return nil, fmt.Errorf("param is empty")
	}
	var keys []string
	m.content.Range(func(key, value any) bool {
		if condition(pattern, key.(string)) {
			keys = append(keys, key.(string))
		}
		return true
	})

	return keys, nil
}

func condition(condition string, key string) bool {
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

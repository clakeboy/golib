package components

import (
	"io/ioutil"
	"ck_res_go/utils"
	"encoding/json"
	"time"
	"os"
	"errors"
)

type CacheIn interface {
	Get(key string) (string,error)
	Set(key string,value interface{},expire int64) bool
	Delete(key string) bool
}

const (
	CACHE_FILE  = 1 + iota
	CACHE_MEM
)

type Cache struct {
	cache_driver CacheIn
}

type CacheData struct {
	Expire int64 `json:"expire"`
	Content string `json:"content"`
}

type FileCache struct {
	CacheIn
	dir string
	prefix string
}

type MemCache struct {
	CacheIn
	server string
}

func NewCache(cache_type int) *Cache{
	var c *Cache
	if cache_type == CACHE_FILE {
		f := &FileCache{dir:"./cache/data/",
			prefix:"ck_"}
		c = &Cache{cache_driver:f}
	} else if cache_type == CACHE_MEM {
		c = nil
	}
	return c
}

/**
 * FileCache Functions
 */
func (this *FileCache) Get(key string) (string,error) {
	file_name := this.GetName(key)

	con,err := ioutil.ReadFile(file_name)
	if err != nil {
		return "",err
	}

	var data = CacheData{}
	err = json.Unmarshal(con,&data)

	if data.Expire <= time.Now().Unix() {
		os.Remove(file_name)
		return "",errors.New("cache expire")
	}

	return data.Content,nil
}

func (this *FileCache)Set(key string,v interface{},expire int64) bool {
	cache_con,err := json.Marshal(v)

	if err != nil  {
		panic(err)
	}

	file_name := this.GetName(key)
	data := CacheData{Expire:time.Now().Unix()+expire,
		Content:string(cache_con)}
	con,err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(this.dir,0755)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(file_name,con,0755)
	if err != nil {
		panic(err)
	}

	return true
}

func (this *FileCache) Delete(key string) bool {
	file_name := this.GetName(key)
	err := os.Remove(file_name)
	if err != nil {
		return false
	}
	return true
}

func (this *FileCache) GetName(key string) string {
	return this.dir+this.prefix+utils.EncodeMD5(key)
}

/**
 * Cache Functions
 */
func (this *Cache) Get(key string) (string,error) {
	return this.cache_driver.Get(key)
}

func (this *Cache)Set(key string,v interface{},expire int64) bool {
	return this.cache_driver.Set(key,v,expire)
}

func (this *Cache) Delete(key string) bool {
	return this.cache_driver.Delete(key)
}
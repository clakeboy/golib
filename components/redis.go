package components

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
	"errors"
)

//redis 配置
type RedisConfig struct {
	RDServer   string `json:"rd_server"`
	RDPort     string `json:"rd_port"`
	RDDb       int    `json:"rd_db"`
	RDPassword string `json:"rd_password"`
	RDListName string `json:"rd_list_name"`
	RDPoolSize int `json:"rd_pool_size"`
	RDIdleSize int `json:"rd_idle_size"`
}

type CKRedis struct {
	rd    redis.Conn
	listName string
	conf     *RedisConfig

}

var CKRedisPool *redis.Pool

func InitRedisPool(cfg *RedisConfig) {
	CKRedisPool = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     cfg.RDIdleSize,
		MaxActive:   cfg.RDPoolSize,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			opt := []redis.DialOption{
				redis.DialDatabase(cfg.RDDb),
				redis.DialPassword(cfg.RDPassword),
			}
			addr := fmt.Sprintf("%s:%s", cfg.RDServer, cfg.RDPort)
			r, err := redis.Dial("tcp", addr, opt...)
			if err != nil {
				return nil, err
			}
			return r,nil
		},
	}
}

//初始一个消息队列
func NewCKRedis() (*CKRedis, error) {
	if CKRedisPool == nil {
		return nil,errors.New("not pool init")
	}

	mq := &CKRedis{rd: CKRedisPool.Get()}

	return mq, nil
}

//设置一个缓存值
func (m *CKRedis) Set(key string, val interface{}, exp int) error {
	_, err := m.rd.Do("SET", key, val,"EX", exp)
	if err != nil {
		return err
	}
	return nil
}

//得到一个缓存值
func (m *CKRedis) Get(key string) (interface{}, error) {
	val, err := m.rd.Do("GET", key)
	if err != nil {
		return nil, err
	}

	return val, nil
}

//删除一个缓存值
func (m *CKRedis) Remove(key ...string) error {
	_,err := m.rd.Do("DEL",key)
	if err != nil {
		return err
	}
	return nil
}

//设置当前操作的消息队列名称
func (m *CKRedis) SetMessageQueueName(name string) {
	m.listName = name
}

//PUSH 消息到队列
func (m *CKRedis) Push(msg string) error {
	_, err := m.rd.Do("LPUSH", m.listName, msg)
	if err != nil {
		return err
	}
	return nil
}

//接收队列消息
func (m *CKRedis) Receive() (string, error) {
	return redis.String(m.rd.Do("RPOP", m.listName))
}

//得到队列消息数
func (m *CKRedis) Count() (int, error) {
	v, err := redis.Int(m.rd.Do("LLEN", m.listName))
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

//设置当前操作的数据ID
func (m *CKRedis) SetDB(db_idx int) {
	m.rd.Do("SELECT",db_idx)
}

//得到KEYS
func (m *CKRedis) Keys(perm string) ([]string,error) {
	return redis.Strings(m.rd.Do("KEYS",perm))
}

//关闭连接
func (m *CKRedis) Close() {
	m.rd.Close()
}
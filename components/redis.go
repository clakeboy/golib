package components

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
	"errors"
	"strconv"
)

//redis 配置
type RedisConfig struct {
	RDServer   string `json:"rd_server" yaml:"rd_server"`
	RDPort     string `json:"rd_port" yaml:"rd_port"`
	RDDb       int    `json:"rd_db" yaml:"rd_db"`
	RDPassword string `json:"rd_password" yaml:"rd_password"`
	RDListName string `json:"rd_list_name" yaml:"rd_list_name"`
	RDPoolSize int `json:"rd_pool_size" yaml:"rd_pool_size"`
	RDIdleSize int `json:"rd_idle_size" yaml:"rd_idle_size"`
}

//GEO地理位置
type RedisGeo struct {
	Member string `json:"member"`
	Lon    float64 `json:"lon"`
	Lat    float64 `json:"lat"`
	Dist   float64 `json:"dist"`
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

//添加一个member有序集合
func (m *CKRedis) ZAdd(key string,score int,member string) (bool,error) {
	return redis.Bool(m.rd.Do("ZADD",key,score,member))
}

//返回有序集合的基数
func (m *CKRedis) ZCard(key string) (int,error) {
	return redis.Int(m.rd.Do("ZCARD",key))
}

//返回指定score 大小之间的成员数量
func (m *CKRedis) ZCount(key string,min string,max string) (int,error) {
	return redis.Int(m.rd.Do("ZCOUNT",key,min,max))
}

//返回指定 score 大小之间的成员列表
func (m *CKRedis) ZRangeByScore(key string,min string,max string,with_score bool) ([]string,error) {
	args := []interface{}{
		key,min,max,
	}

	if with_score {
		args = append(args,"WITHSCORES")
	}

	return redis.Strings(m.rd.Do("ZRANGEBYSCORE",args...))
}
//删除有序集合里面的成员
func (m *CKRedis) ZRem(key ,member string) (bool,error) {
	return redis.Bool(m.rd.Do("ZREM",key,member))
}

//----------------------------------- GEOHASH

//添加一个坐标到 GEOHASH
func (m *CKRedis) GAdd(key string,lon string,lat string,member string) (bool,error) {
	return redis.Bool(m.rd.Do("GEOADD",key,lon,lat,member))
}

//得到一个member 的坐标
func (m *CKRedis) GPos(key string,member string) ([]interface{},error) {
	return redis.Values(m.rd.Do("GEOPOS",key,member))
}

//计算两个位置的距离,返回米
func (m *CKRedis) GDist(key string,member string,member2 string) (int,error) {
	return redis.Int(m.rd.Do("GEODIST",key,member,member2))
}

//得到传入坐标半径的所有所坐标
func (m *CKRedis) GRadius(key string,lon string,lat string,radius int) ([]*RedisGeo,error) {
	list,err := redis.Values(m.rd.Do("GEORADIUS",key,lon,lat,radius,"m","WITHCOORD","WITHDIST","ASC"))
	if err != nil {
		return nil,err
	}
	return m.transGeo(list),nil
}

//转换取得的坐标数据为geo结构体
func (m *CKRedis) transGeo(geo_list []interface{}) []*RedisGeo {
	list := []*RedisGeo{}

	for _,v := range geo_list {
		val := v.([]interface{})
		geo := &RedisGeo{}
		geo.Member = string(val[0].([]byte))
		geo.Dist,_ = strconv.ParseFloat(string(val[1].([]byte)),64)
		geo.Lon,_ = strconv.ParseFloat(string(val[2].([]interface{})[0].([]byte)),64)
		geo.Lat,_ = strconv.ParseFloat(string(val[2].([]interface{})[1].([]byte)),64)
		list = append(list,geo)
	}
	return list
}

//--- HASH ---
//设置一个HASH值
func (m *CKRedis) HSet(key,field string,val interface{}) (bool,error) {
	return redis.Bool(m.rd.Do("HSET",key,field,val))
}

//得到一个HASH值
func (m *CKRedis) HGet(key,field string) (interface{},error) {
	val,err := m.rd.Do("HGET",key,field)
	if err != nil {
		return nil,err
	}
	return val,nil
}
//检查是否存在一个HASH值
func (m *CKRedis) HExists(key,field string) (bool,error) {
	return redis.Bool(m.rd.Do("HEXISTS",key,field))
}

//得到多个HASH 值
func (m *CKRedis) HMGet(key string,field ...interface{}) ([]interface{},error) {
	return redis.Values(m.rd.Do("HMGET",field...))
}

//删除一个 hash 值
func (m *CKRedis) HDel(key ,field string) (bool,error) {
	return redis.Bool(m.rd.Do("HDEL",key,field))
}

//得到所有 hash 键
func (m *CKRedis) HKeys(key string) ([]interface{},error) {
	return redis.Values(m.rd.Do("HKEYS",key))
}
//得到 hash 长度
func (m *CKRedis) HLen(key string) (int,error) {
	return redis.Int(m.rd.Do("HLEN",key))
}

//执行命令
func (m *CKRedis) Do(command string,args ...interface{}) (interface{},error) {
	return m.rd.Do(command,args...)
}

//关闭连接
func (m *CKRedis) Close() {
	m.rd.Close()
}
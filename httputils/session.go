package httputils

import (
	"encoding/json"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/components/task"
	"github.com/clakeboy/golib/utils"
	"time"
)

type SessionType int

const (
	SessionMem SessionType = iota
	SessionRedis
	SessionFile
)

var (
	memDriver   *components.MemCache
	redisDriver *components.CKRedis
	boltDriver  *ckdb.BoltDB
)

var initOptions *SessionOptions
var sessiongc *task.Management

//session 数据
type SessionData struct {
	Key    string  `json:"key"`    //session key
	Value  utils.M `json:"data"`   //session value
	Expire int     `json:"expire"` //过期时间
}

func (m *SessionData) ToJson() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return data
}

func (m *SessionData) ToJsonString() string {
	data := m.ToJson()
	if data == nil {
		return ""
	}
	return string(data)
}

func (m *SessionData) ParseJson(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *SessionData) ParseJsonString(data string) error {
	return m.ParseJson([]byte(data))
}

//session 选项
type SessionOptions struct {
	StorageType  SessionType   //Session 存储类型
	SurvivalTime time.Duration //Session 存活时间
	CookieName   string        //Session cookie name
}

type HttpSession struct {
	options    *SessionOptions
	cookie     *HttpCookie
	sessionKey string
	data       *SessionData
}

func InitSession(options *SessionOptions) {
	if options == nil {
		options = &SessionOptions{
			StorageType:  SessionMem,
			SurvivalTime: time.Minute * 20,
			CookieName:   "CK-SESSION",
		}
	}
	switch options.StorageType {
	case SessionMem:
		memDriver = components.NewMemCache()
	case SessionFile:
		boltDriver = ckdb.NewBoltDB("./session/")
	case SessionRedis:
		redisDriver, _ = components.NewCKRedis()
	}
	initOptions = options
}

func NewHttpSession(cookie *HttpCookie) *HttpSession {
	sessionKey, err := cookie.Get(initOptions.CookieName)
	if err != nil || sessionKey == "" {
		sessionKey = utils.CreateUUID(false)
	}

	return &HttpSession{
		options:    initOptions,
		cookie:     cookie,
		sessionKey: sessionKey,
	}
}

//开始Session
func (s *HttpSession) Start() {
	sData := &SessionData{}

	var data interface{}
	var err error
	switch s.options.StorageType {
	case SessionMem:
		data, err = memDriver.Get(s.sessionKey)
		if err == nil {
			sData.ParseJsonString(data.(string))
		}
	case SessionFile:
		data, err = boltDriver.Get("session", s.sessionKey)
		if err == nil {
			jserr := sData.ParseJson(data.([]byte))
			if jserr != nil {
				err = jserr
				break
			}
			if int64(sData.Expire) < time.Now().Unix() {
				boltDriver.Delete("session", s.sessionKey)
				sData.Value = utils.M{}
				sData.Expire = int(time.Now().Add(s.options.SurvivalTime).Unix())
			}
		}
	case SessionRedis:
		data, err = redisDriver.Get(s.sessionKey)
		if err == nil {
			sData.ParseJsonString(data.(string))
		}
	}

	if err != nil {
		sData.Key = s.sessionKey
		sData.Expire = int(time.Now().Add(s.options.SurvivalTime).Unix())
		sData.Value = utils.M{}
	}

	s.data = sData
	s.cookie.Set(s.options.CookieName, s.sessionKey, 3600*24*365*10)
}

//设置一个Session 值
func (s *HttpSession) Set(name string, val string) {
	s.data.Value[name] = val
}

//设置一个Session 值
func (s *HttpSession) Get(name string) string {
	val, ok := s.data.Value[name].(string)
	if ok {
		return val
	}
	return ""
}

//将SESSION 回写
func (s *HttpSession) Flush() {
	if s.data == nil {
		return
	}
	switch s.options.StorageType {
	case SessionMem:
		memDriver.Set(s.sessionKey, s.data.ToJsonString(), int64(s.options.SurvivalTime.Seconds()))
	case SessionFile:
		boltDriver.Put("session", s.sessionKey, s.data)
	case SessionRedis:
		redisDriver.Set(s.sessionKey, s.data.ToJsonString(), int(s.options.SurvivalTime.Seconds()))
	}
}

////初始化GC
//func (s *HttpSession) InitGc() {
//	if s.options.StorageType != SessionFile || sessiongc != nil {
//		return
//	}
//
//	sessiongc = task.NewManagement()
//	sessiongc.AddTaskString("* */1 * * * *", func(item *task.Item) bool {
//		if boltDriver != nil {
//			var keys []string
//			sessionData := &SessionData{}
//			currentUnix := int(time.Now().Unix())
//			boltDriver.ForEach("session", func(key []byte, value []byte) error {
//				err := sessionData.ParseJson(value)
//				if err == nil {
//					if currentUnix > sessionData.Expire {
//						keys = append(keys, string(key))
//					}
//				}
//				fmt.Println(sessionData)
//				return nil
//			})
//			if len(keys) > 0 {
//				boltDriver.Delete("session", keys...)
//			}
//			fmt.Println("session has been gc")
//		}
//		return false
//	}, nil)
//	sessiongc.Start()
//}

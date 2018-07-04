package components

import (
	"encoding/json"
	"ck_go_lib/ckdb"
	"time"
	"ck_go_lib/utils"
)

type SessionType int
const(
	SessionMem SessionType = iota
	SessionRedis
	SessionFile
)

var (
	memDriver *MemCache
	redisDriver *CKRedis
	boltDriver *ckdb.BoltDB
)

var sessionInit = false

//session 数据
type SessionData struct {
	Key   string `json:"key"`  //session key
	Value utils.M `json:"data"`  //session value
	Expire int `json:"expire"`  //过期时间
}

func (m *SessionData) ToJson() []byte {
	data ,err := json.Marshal(m)
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
	err := json.Unmarshal(data,m)
	if err != nil{
		return err
	}
	return nil
}

func (m *SessionData) ParseJsonString(data string) error {
	return m.ParseJson([]byte(data))
}

//session 选项
type SessionOptions struct {
	StorageType SessionType  //Session 存储类型
	SurvivalTime  time.Duration       //Session 存活时间
	CookieName  string       //Session cookie name
}

type HttpSession struct {
	options *SessionOptions
	cookie *HttpCookie
	sessionKey string
	data   *SessionData
}

func NewHttpSession(cookie *HttpCookie,options *SessionOptions) *HttpSession {
	if options == nil {
		options = &SessionOptions{
			StorageType:SessionFile,
			SurvivalTime:time.Minute*20,
			CookieName:"CK-SESSION",
		}
	}
	if (!sessionInit) {
		switch options.StorageType {
		case SessionMem:
			memDriver = NewMemCache()
		case SessionFile:
			boltDriver = ckdb.NewBoltDB("./session/")
		case SessionRedis:
			redisDriver,_ = NewCKRedis()
		}
		sessionInit = true
	}

	sessionKey,err := cookie.Get(options.CookieName)
	if err != nil || sessionKey == "" {
		sessionKey = utils.CreateUUID(false)
	}

	return &HttpSession{
		options:options,
		cookie:cookie,
		sessionKey:sessionKey,
	}
}
//开始Session
func (s *HttpSession) Start() {
	sData := &SessionData{}

	var data interface{}
	var err error
	switch s.options.StorageType {
	case SessionMem:
		data ,err = memDriver.Get(s.sessionKey)
		if err == nil {
			sData.ParseJsonString(data.(string))
		}
	case SessionFile:
		data ,err = boltDriver.Get("session",s.sessionKey)
		if err == nil {
			sData.ParseJson(data.([]byte))
			if int64(sData.Expire) < time.Now().Unix() {
				boltDriver.Delete("session",s.sessionKey)
				sData.Value = utils.M{}
				sData.Expire = int(time.Now().Add(s.options.SurvivalTime).Unix())
			}
		}
	case SessionRedis:
		data ,err = redisDriver.Get(s.sessionKey)
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
	s.cookie.Set(s.options.CookieName,s.sessionKey,3600*24*365*10)
}

//设置一个Session 值
func (s *HttpSession) Set(name string,val string) {
	s.data.Value[name] = val
}
//设置一个Session 值
func (s *HttpSession) Get(name string) string {
	return s.data.Value[name].(string)
}

//将SESSION 回写
func (s *HttpSession) Flush() {
	if s.data == nil {
		return
	}
	switch s.options.StorageType {
	case SessionMem:
		memDriver.Set(s.sessionKey,s.data.ToJsonString(),int64(s.options.SurvivalTime.Seconds()))
	case SessionFile:
		boltDriver.Put("session",s.sessionKey,s.data.ToJsonString())
	case SessionRedis:
		redisDriver.Set(s.sessionKey,s.data.ToJsonString(),int(s.options.SurvivalTime.Seconds()))
	}
}
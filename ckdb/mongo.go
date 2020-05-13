package ckdb

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

//dsn mongodb://root:WiaQ82n7B3L5Cz*2#10m@172.18.76.150:27017?authSource=admin
type MongoDBConfig struct {
	DBHost     string `json:"db_host" yaml:"db_host"`
	DBPort     string `json:"db_port" yaml:"db_port"`
	DBUser     string `json:"db_user" yaml:"db_user"`
	DBPasswd   string `json:"db_passwd" yaml:"db_passwd"`
	DBAuth     string `json:"db_auth" yaml:"db_auth"`
	DBName     string `json:"db_name" yaml:"db_name"`
	DBPoolSize int    `json:"db_pool_size" yaml:"db_pool_size"`
}

//build dsn string
func (mc *MongoDBConfig) BuildDsn() string {
	if mc.DBAuth == "" {
		return fmt.Sprintf("mongodb://%s:%s", mc.DBHost, mc.DBPort)
	}
	return fmt.Sprintf("mongodb://%s:%s@%s:%s?authSource=%s", mc.DBUser, mc.DBPasswd, mc.DBHost, mc.DBPort, mc.DBAuth)
}

type DBMongo struct {
	is_open  bool
	db_name  string
	session  *mgo.Session
	database *mgo.Database
}

var globalSession *mgo.Session

//new init mongodb
func InitMongo(conf *MongoDBConfig) error {
	return InitDB(conf.BuildDsn(), conf.DBPoolSize)
}

//old init
func InitDB(db_dsn string, pool_size int) error {
	var err error
	globalSession, err = mgo.Dial(db_dsn)
	if err != nil {
		return err
	}
	globalSession.SetPoolLimit(pool_size)
	globalSession.SelectServers()
	return nil
}

func NewDB(db_name string) *DBMongo {
	db := new(DBMongo)
	db.db_name = db_name
	return db
}

func (this *DBMongo) Open(db_name string) *DBMongo {
	if !this.is_open {
		this.session = globalSession.Clone()
		this.is_open = true
		this.database = this.session.DB(db_name)
	} else {
		this.database = this.session.DB(db_name)
	}

	return this
}

func (this *DBMongo) Table(tab_name string) *mgo.Collection {
	if !this.is_open {
		this.Open(this.db_name)
	}

	return this.database.C(tab_name)
}

func (this *DBMongo) ChangeDB(db_name string) {
	if !this.is_open {
		this.Open(this.db_name)
	}
	this.db_name = db_name
	this.database = this.session.DB(db_name)
}

func (this *DBMongo) Close() {
	if this.is_open {
		this.is_open = false
		this.session.Close()
	}
}

func (this *DBMongo) Collection(collection_name string) *CKCollection {
	return NewCollection(this, collection_name)
}

func (this *DBMongo) RunCmd() {

}

package ckdb

import (
	"gopkg.in/mgo.v2"
)

type MongoDBConfig struct {
	DBDsn string `json:"db_dsn"`
	DBName string `json:"db_name"`
	DBPoolSize int `json:"db_pool_size"`
}

type DBMongo struct {
	is_open bool
	db_name string
	session *mgo.Session
	database *mgo.Database
}

var globalSession *mgo.Session

func InitDB(db_dsn string,pool_size int) error {
	var err error
	globalSession,err = mgo.Dial(db_dsn)
	if err != nil {
		return err
	}
	globalSession.SetPoolLimit(pool_size)
	return nil
}

func NewDB(db_name string) *DBMongo{
	db := new(DBMongo)
	db.db_name = db_name
	return db
}

func (this *DBMongo) Open(db_name string) *DBMongo{
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

func (this *DBMongo) Close() {
	if this.is_open {
		this.is_open = false
		this.session.Close()
	}
}

func (this *DBMongo) Collection(collection_name string) *CKCollection {
	return NewCollection(this,collection_name)
}
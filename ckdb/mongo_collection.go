package ckdb

type Collection struct {
	db *DBMongo
}

func NewCollection(mdb *DBMongo) *Collection {
	return &Collection{db:mdb}
}


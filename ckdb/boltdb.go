package ckdb

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"wx_shake/utils"
)
//创建一个新的 BlotDB
func NewBoltDB(filepath string) *BoltDB {
	if !utils.PathExists(filepath) {
		err := os.MkdirAll(filepath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err := bolt.Open(filepath+"mp.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &BoltDB{db}
}

type BoltDB struct {
	db *bolt.DB
}

func (b *BoltDB) Path() string {
	return b.db.Path()
}

func (b *BoltDB) Close() {
	b.db.Close()
}

func (b *BoltDB) Get(bucket_name string, key string) ([]byte, error) {
	var v []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte(bucket_name))
		if bu == nil {
			return nil
		}
		v = bu.Get([]byte(key))
		return nil
	})

	if err != nil {
		return nil, err
	}
	return v, nil
}

func (b *BoltDB) Put(bucket_name string, key string, val interface{}) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bu,err := tx.CreateBucketIfNotExists([]byte(bucket_name))
		if err != nil {
			return err
		}

		//row := bu.Get([]byte(key))

		data, err := json.Marshal(val)
		if err != nil {
			return err
		}

		err = bu.Put([]byte(key), data)
		return err
	})

	if err != nil {
		return err
	}
	return nil
}

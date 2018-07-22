/*
数据库处理
 */
package ckdb

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"ck_go_lib/utils"
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
//bolt DB 数据库处理
type BoltDB struct {
	db *bolt.DB
}
//返回当前 bolt DB 位置
func (b *BoltDB) Path() string {
	return b.db.Path()
}
//关闭数据库
func (b *BoltDB) Close() {
	b.db.Close()
}
//从 bucket 传入 KEY 得到一个值
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
//插一个新的 key value 值到 bucket
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

	return err
}
//删除一个值
func (b *BoltDB) Delete(bucket_name string, key ...string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bu,err := tx.CreateBucketIfNotExists([]byte(bucket_name))
		if err != nil {
			return err
		}

		for _,v := range key {
			err = bu.Delete([]byte(v))
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}
//迭代一个bucket里所有的数据
func (b *BoltDB) ForEach(bucketName string,callback func([]byte,[]byte) error) error {
	err := b.db.View(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte(bucketName))
		if bu == nil {
			return nil
		}

		return bu.ForEach(callback)
	})

	return err
}

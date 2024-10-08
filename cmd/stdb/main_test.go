package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/clakeboy/golib/utils"
	"go.etcd.io/bbolt"
)

type AccountData struct {
	Id           int    `storm:"id,increment" json:"id"` //主键,自增长
	Name         string `storm:"index" json:"name"`      //用户名
	Passwd       string `json:"passwd"`                  //密码，默认密码都是1230123
	Phone        string `json:"phone"`                   //管理员手机
	Manage       int    `json:"manage"`                  //是否管理员
	Init         int    `json:"init"`                    //是否初始化 0，1，如果为0强制修改密码
	CreatedDate  int64  `json:"created_date"`            //创建时间
	ModifiedDate int64  `json:"modified_date"`           //修改时间
}

func TestStorm(t *testing.T) {
	var err error
	db, err = storm.Open("/Users/clakeboy/Documents/go-mod-projects/pingan_insurance_service/frontend/db/sys.db", storm.BoltOptions(0, &bbolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: true,
	}))
	if err != nil {
		t.Error(err)
		return
	}
	db.Bolt.View(func(tx *bbolt.Tx) error {
		var list []string
		tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			list = append(list, string(name))
			return nil
		})

		b := tx.Bucket([]byte("car_policy"))
		// b.Cursor().Bucket().Tx().ForEach(func(name []byte, b *bbolt.Bucket) error {
		// 	list = append(list, string(name))
		// 	return nil
		// })
		bs := b.Stats()
		utils.PrintAny(bs)
		utils.PrintAny(list)
		b = b.Bucket([]byte("CarPolicyData"))
		// b = b.Bucket([]byte("__storm_index_OrderNo"))
		// b = b.Bucket([]byte("storm__ids"))
		// b = b.Bucket([]byte("__storm_metadata"))

		c := b.Cursor()
		prefix := []byte("__storm")
		k, v := c.Seek(prefix)
		fmt.Printf("seek: key=%s, value=%x\n", k, v)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Printf("key=%s, value=%x\n", k, v)
		}
		// b.ForEach(func(k, v []byte) error {
		// 	idx, err := utils.BytesToInt64(k)
		// 	if err != nil {
		// 		fmt.Printf("key: %s, value: %X\n", string(k), v)
		// 	} else {
		// 		fmt.Printf("key: %d, value: %s\n", idx, v)
		// 	}

		// 	// 	// fmt.Printf("key: %X, value: %X\n", k, v)
		// 	return nil
		// })
		return nil
	})
	// var list []map[string]interface{}
	// query := db.From("Account").Select()
	// err = query.Limit(1).Find(&list)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// utils.PrintAny(list)
}

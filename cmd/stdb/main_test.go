package main

import (
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
	db, err = storm.Open("/Users/clakeboy/Documents/pcbx_project/pcbx-btm/db/sys.db", storm.BoltOptions(0, &bbolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: true,
	}))
	if err != nil {
		t.Error(err)
		return
	}
	db.Bolt.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Account"))
		b = b.Bucket([]byte("AccountData"))
		// b = b.Bucket([]byte("__storm_index_Name"))
		b = b.Bucket([]byte("__storm_metadata"))
		// b = b.Bucket([]byte("storm__ids"))
		// c := b.Cursor()

		// prefix := []byte("clake")
		// for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		// 	fmt.Printf("key=%s, value=%x\n", k, v)
		// }
		b.ForEach(func(k, v []byte) error {
			idx, err := utils.BytesToInt64(k)
			if err != nil {
				fmt.Printf("key: %s, value: %s\n", string(k), string(v))
			} else {
				fmt.Printf("key: %d, value: %s\n", idx, string(v))
			}

			fmt.Printf("key: %X, value: %X\n", k, v)
			return nil
		})
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

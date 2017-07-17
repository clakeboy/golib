package ckdb

import (
	"reflect"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"ck_go_lib/utils"
)

//查询结果
type QueryResult struct {
	List  interface{} `json:"list"`
	Count int         `json:"count"`
}

type CKCollection struct {
	db *DBMongo
	tab string
}

func NewCollection(mdb *DBMongo,tab_name string) *CKCollection {
	return &CKCollection{db:mdb,tab:tab_name}
}

//插入数据
func (ck *CKCollection) Insert(rows ...interface{}) error {
	c := ck.db.Table(ck.tab)
	err := c.Insert(rows...)
	if err != nil {
		return err
	}
	return nil
}

//更新数据
func (ck *CKCollection) Update(where bson.M, update bson.M) error {
	c := ck.db.Table(ck.tab)
	err := c.Update(where, update)
	if err != nil {
		return err
	}
	return nil
}

//更新所有条件数据
func (ck *CKCollection) UpdateAll(where bson.M, update bson.M) (*mgo.ChangeInfo,error) {
	c := ck.db.Table(ck.tab)
	return c.UpdateAll(where, update)
}

//删除数据
func (ck *CKCollection) Delete(where bson.M) error {
	c := ck.db.Table(ck.tab)
	_, err := c.RemoveAll(where)
	if err != nil {
		return err
	}
	return nil
}

//查找数据
func (ck *CKCollection) Find(where bson.M, row interface{}) error {
	c := ck.db.Table(ck.tab)
	err := c.Find(where).One(row)
	if err != nil {
		return err
	}
	return nil
}

//查询数据库
func (ck *CKCollection) Query(where bson.M, page int, number int,sort_list []string, struct_type interface{}, format func(interface{})) (*QueryResult, error) {
	c := ck.db.Table(ck.tab)
	var list []interface{}
	var err error
	var count int
	var res_list *mgo.Iter

	if sort_list == nil {
		sort_list = []string{"-_id"}
	}

	if where == nil {
		res_list = c.Find(nil).Sort(sort_list...).Skip((page - 1) * number).Limit(number).Iter()
		count, _ = c.Find(nil).Count()
	} else {
		res_list = c.Find(where).Sort(sort_list...).Skip((page - 1) * number).Limit(number).Iter()
		count, _ = c.Find(where).Count()
	}

	if struct_type == nil {
		struct_type = utils.M{}
	}

	t := reflect.TypeOf(struct_type)
	result := reflect.New(t).Interface()
	for res_list.Next(result) {
		if format != nil {
			format(result)
		}
		list = append(list, result)
		result = reflect.New(t).Interface()
	}

	res := &QueryResult{
		List:  list,
		Count: count,
	}

	return res, err
}

//关闭数据库连接
func (ck *CKCollection) Close() {
	ck.db.Close()
}
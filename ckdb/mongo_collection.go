package ckdb

import (
	"ck_go_lib/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

//查询结果
type QueryResult struct {
	List  interface{} `json:"list"`
	Count int         `json:"count"`
}

type CKCollection struct {
	db  *DBMongo
	tab string
}

func NewCollection(mdb *DBMongo, tab_name string) *CKCollection {
	return &CKCollection{db: mdb, tab: tab_name}
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
func (ck *CKCollection) UpdateAll(where bson.M, update bson.M) (*mgo.ChangeInfo, error) {
	c := ck.db.Table(ck.tab)
	return c.UpdateAll(where, update)
}

func (ck *CKCollection) Upset(where bson.M, update bson.M) error {
	c := ck.db.Table(ck.tab)
	_, err := c.Upsert(where, update)
	if err != nil {
		return err
	}

	return nil
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
//得到所给条件的数据量
func (ck *CKCollection) Count(where bson.M) int {
	c := ck.db.Table(ck.tab)
	count,err := c.Find(where).Count()
	if err != nil {
		return 0
	}
	return count
}

//查询数据库
func (ck *CKCollection) Query(where bson.M, page int, number int, sort_list []string, struct_type interface{}, format func(interface{})) (*QueryResult, error) {
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

	result := ck.getQueryType(struct_type)
	for res_list.Next(result) {
		if format != nil {
			format(result)
		}
		list = append(list, result)
		result = ck.getQueryType(struct_type)
	}
	res := &QueryResult{
		List:  list,
		Count: count,
	}

	return res, err
}

//查询数据库返回列表
func (ck *CKCollection) List(where bson.M, page int, number int, sort_list []string, struct_type interface{}, format func(interface{})) ([]interface{}, error) {
	c := ck.db.Table(ck.tab)
	var list []interface{}
	var err error
	var res_list *mgo.Iter

	if sort_list == nil {
		sort_list = []string{"-_id"}
	}

	if where == nil {
		res_list = c.Find(nil).Sort(sort_list...).Skip((page - 1) * number).Limit(number).Iter()
	} else {
		res_list = c.Find(where).Sort(sort_list...).Skip((page - 1) * number).Limit(number).Iter()
	}

	result := ck.getQueryType(struct_type)
	for res_list.Next(result) {
		if format != nil {
			format(result)
		}
		list = append(list, result)
		result = ck.getQueryType(struct_type)
	}
	return list, err
}
//执行聚合操作
func (ck *CKCollection) Aggregate(pipe ...bson.M) (bson.M,error) {
	c := ck.db.Table(ck.tab)
	p := c.Pipe(pipe)
	resp := bson.M{}
	err := p.One(&resp)
	if err != nil {
		return nil,err
	}

	return resp,nil
}

func (ck *CKCollection) getQueryType(i interface{}) interface{} {
	if i == nil {
		return utils.M{}
	}

	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		return reflect.New(t.Elem()).Interface()
	}
	return reflect.New(t).Interface()
}

//删除 Collection
func (ck *CKCollection) Drop() error {
	return ck.db.Table(ck.tab).DropCollection()
}

//关闭数据库连接
func (ck *CKCollection) Close() {
	ck.db.Close()
}
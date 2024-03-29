package ckdb

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/clakeboy/golib/utils"
)

type TBField struct {
	Column string //字段名
	Alias  string //别名
	Func   string //所用方法
}

type TBJoin struct {
	TableTo    string
	ColumnTo   string
	TableFrom  string
	ColumnFrom string
	Key        string
}

// Table 主结构
type DBATable struct {
	where      utils.M
	where_str  string
	join_str   string
	group_str  string
	field_str  string
	order_str  string
	limit_str  string
	sql_str    string
	table      string
	values     []interface{}
	db         *DBA
	columnType interface{}
}

// field regexp
var fieldReg = regexp.MustCompile(`(.+?)\[(.+)\]`)

// 新建一个table处理类
func NewDBATable(db *DBA, table string) *DBATable {
	return &DBATable{table: table, db: db, field_str: "*"}
}

// 开始事务
func (t *DBATable) BeginTrans() error {
	return t.db.BeginTrans()
}

// 提交事务
func (t *DBATable) Commit() error {
	return t.db.Commit()
}

// 回滚事务
func (t *DBATable) Rollback() error {
	return t.db.Rollback()
}

// 设置要显示的字段
func (t *DBATable) Select(fields utils.M) *DBATable {
	var tmp []string
	for column, column_table := range fields {
		var table string
		if column_table.(string) == "" {
			table = t.db.FormatColumn(t.table)
		} else {
			table = t.db.FormatColumn(column_table.(string))
		}
		tb := t.explainField(column)
		if tb.Alias == "" {
			if tb.Func == "" {
				tmp = append(tmp, fmt.Sprintf("%s.%s", table, t.db.FormatColumn(tb.Column)))
			} else {
				tmp = append(tmp, fmt.Sprintf("%s(%s.%s)", tb.Func, table, t.db.FormatColumn(tb.Column)))
			}
		} else {
			if tb.Func == "" {
				tmp = append(tmp, fmt.Sprintf("%s.%s AS '%s'", table, t.db.FormatColumn(tb.Column), tb.Alias))
			} else {
				tmp = append(tmp, fmt.Sprintf("%s(%s.%s) AS '%s'", tb.Func, table, t.db.FormatColumn(tb.Column), tb.Alias))
			}
		}
	}
	t.field_str = strings.Join(tmp, ",")
	return t
}

// 解释字段内是否有别名
func (t *DBATable) explainField(field string) *TBField {
	//reg := regexp.MustCompile(`(.+?)\[(.+)\]`)
	match := fieldReg.FindStringSubmatch(field)
	var (
		fieldStr string
		funcStr  string
	)

	if len(match) > 0 {
		fieldStr = match[1]
		funcStr = strings.ToUpper(match[2])
	} else {
		fieldStr = field
	}

	tmp := strings.Split(fieldStr, " ")
	tb := &TBField{Column: tmp[0], Func: funcStr}
	if len(tmp) > 1 {
		tb.Alias = tmp[1]
	}
	return tb
}

// 设置WHERE条件
func (t *DBATable) Where(fields utils.M, table string) *DBATable {
	t.where = fields
	if table == "" {
		table = t.table
	}
	if len(fields) <= 0 {
		return t
	}
	where_str, val := t.db.WhereRecursion(fields, "AND", table)
	t.where_str = t.where_str + where_str
	t.values = append(t.values, val...)
	return t
}

// 设置where and
func (t *DBATable) WhereAnd(fields utils.M, table string) *DBATable {
	t.where = fields
	if table == "" {
		table = t.table
	}
	whereStr, val := t.db.WhereRecursion(fields, "AND", table)
	t.where_str = fmt.Sprintf("%s AND %s", t.where_str, whereStr)
	t.values = append(t.values, val...)
	return t
}

// 设置where and
func (t *DBATable) WhereOr(fields utils.M, table string) *DBATable {
	t.where = fields
	if table == "" {
		table = t.table
	}
	whereStr, val := t.db.WhereRecursion(fields, "AND", table)
	t.where_str = fmt.Sprintf("%s OR %s", t.where_str, whereStr)
	t.values = append(t.values, val...)
	return t
}

// 添加多外JOIN
func (t *DBATable) Join(fields [][]string) *DBATable {
	for _, join_str := range fields {
		join_to := strings.Split(join_str[0], ".")
		join_from := strings.Split(join_str[1], ".")
		t.JoinOne(&TBJoin{
			TableTo:    join_to[0],
			ColumnTo:   join_to[1],
			TableFrom:  join_from[0],
			ColumnFrom: join_from[1],
			Key:        join_str[2],
		})
	}
	return t
}

// 添加一条JOIN记录
func (t *DBATable) JoinOne(join *TBJoin) *DBATable {
	tbt := t.explainField(join.TableTo)
	if tbt.Alias != "" {
		t.join_str = t.join_str + fmt.Sprintf(" %s JOIN %s %s ON %s.%s=%s.%s", join.Key, tbt.Column, tbt.Alias, tbt.Alias, join.ColumnTo, join.TableFrom, join.ColumnFrom)
	} else {
		t.join_str = t.join_str + fmt.Sprintf(" %s JOIN %s ON %s.%s=%s.%s", join.Key, join.TableTo, join.TableTo, join.ColumnTo, join.TableFrom, join.ColumnFrom)
	}
	return t
}

// 设置ORDER 排序
func (t *DBATable) Order(orders utils.M) *DBATable {
	var tmp []string
	for column, order_type := range orders {
		field := t.explainField(column)
		if field.Alias == "" {
			tmp = append(tmp, fmt.Sprintf("%s.%s %s", t.db.FormatColumn(t.table), t.db.FormatColumn(field.Column), order_type))
		} else {
			tmp = append(tmp, fmt.Sprintf("%s.%s %s", t.db.FormatColumn(field.Alias), t.db.FormatColumn(field.Column), order_type))
		}
	}

	t.order_str = "ORDER BY " + strings.Join(tmp, ",")
	return t
}

// ORDER IN 条件排序
func (t *DBATable) OrderIn(column string, rule []string) *DBATable {
	t.order_str = fmt.Sprintf("ORDER BY FIND_IN_SET(%s,'%s')", column, strings.Join(rule, ","))
	return t
}

// 设置分页
func (t *DBATable) Limit(number int, page int) *DBATable {
	if page == 0 {
		t.limit_str = fmt.Sprintf("LIMIT %d", number)
	} else {
		curnum := 0
		if page > 1 {
			curnum = (page - 1) * number
		}

		t.limit_str = fmt.Sprintf("LIMIT %d,%d", curnum, number)
	}

	return t
}

// 开始查询
func (t *DBATable) Query() *DBATable {
	var where_str string
	if t.where_str != "" {
		where_str = "WHERE " + t.where_str
	}

	t.sql_str = fmt.Sprintf("SELECT %s FROM %s %s %s %s %s",
		t.field_str,
		t.db.FormatColumn(t.table),
		t.join_str,
		where_str,
		t.order_str,
		t.limit_str,
	)

	return t
}

// 得到所有列表结果集
func (t *DBATable) Result() ([]utils.M, error) {
	defer t.Clear()
	return t.db.Query(t.sql_str, t.values...)
}

// 得到所有列表结果集,以传入类型返回该类型数组
func (t *DBATable) ResultStruct(i interface{}) ([]interface{}, error) {
	defer t.Clear()
	return t.db.QueryStruct(i, t.sql_str, t.values...)
}

// 只得到一条记录
func (t *DBATable) Find() (utils.M, error) {
	defer t.Clear()
	return t.db.QueryOne(t.sql_str, t.values...)
}

// 只得到一条记录,以传入类型返回该类型数组
func (t *DBATable) FindStruct(i interface{}) (interface{}, error) {
	defer t.Clear()
	return t.db.QueryOneStruct(i, t.sql_str, t.values...)
}

// 得到记录条数
func (t *DBATable) Rows() int {
	var where_str string
	if t.where_str != "" {
		where_str = "WHERE " + t.where_str
	}

	sql_str := fmt.Sprintf("SELECT count(*) FROM %s %s %s",
		t.db.FormatColumn(t.table),
		t.join_str,
		where_str,
	)

	row := t.db.QueryRow(sql_str, t.values...)
	var length int
	err := row.Scan(&length)
	if err != nil {
		return 0
	}
	return length
}

// 清除所有查询条件
func (t *DBATable) Clear() {
	t.where_str = ""
	t.join_str = ""
	t.group_str = ""
	t.field_str = "*"
	t.order_str = ""
	t.limit_str = ""
	t.sql_str = ""
	t.where = nil
	t.values = []interface{}{}
}

// 插入数据
func (t *DBATable) Insert(data interface{}) (int, bool) {
	return t.db.Insert(t.table, data)
}

// 插入多条数据
func (t *DBATable) InsertMulti(dataList []interface{}) (int, bool) {
	return t.db.InsertMulti(t.table, dataList)
}

// 更新数据
func (t *DBATable) Update(data utils.M) bool {
	defer t.Clear()
	err := t.db.Update(data, t.where, t.table)
	if err != nil {
		return false
	}

	return true
}

// 更新整条记录
func (t *DBATable) UpdateAny(data interface{}) bool {
	defer t.Clear()
	err := t.db.UpdateAny(data, t.where, t.table)
	if err != nil {
		return false
	}

	return true
}

func (t *DBATable) One(rowStruct interface{}) error {
	defer t.Clear()
	return nil
}

// 删除数据
func (t *DBATable) Delete() bool {
	defer t.Clear()
	_, err := t.db.Delete(t.where, t.table)
	if err != nil {
		return false
	}

	return true
}

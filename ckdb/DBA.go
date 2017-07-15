package ckdb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"regexp"
	"strings"
	"ck_go_lib/utils"
	"errors"
)

type DBA struct {
	db        *sql.DB
	table     string
	debug     bool
	LastSql  string
	LastArgs []interface{}
}

//数据库配置
type DBConfig struct {
	DBServer   string `json:"db_server"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBPoolSize int    `json:"db_pool_size"`
	DBIdleSize int    `json:"db_Idle_size"`
	DBDebug    bool   `json:"db_debug"`
}

type DBColumn struct {
	Field string
	Icon  string
}

//DBA专用数据
type DM map[string]interface{}

//新创建DBA操作库
func NewDBA(db_conf *DBConfig) (*DBA, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		db_conf.DBUser,
		db_conf.DBPassword,
		db_conf.DBServer,
		db_conf.DBPort,
		db_conf.DBName,
	)

	db_driver, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db_driver.SetMaxOpenConns(db_conf.DBPoolSize)
	db_driver.SetMaxIdleConns(db_conf.DBIdleSize)
	err = db_driver.Ping()
	if err != nil {
		return nil, err
	}
	dba := &DBA{db: db_driver, debug: db_conf.DBDebug}

	return dba, nil
}

//设置操作的表名
func (d *DBA) Table(table_name string) *DBATable {
	return NewDBATable(d,table_name)
}

//插入记录到数据库
func (d *DBA) Insert(table string, org_data interface{}) (int, bool) {
	var columns []string
	var values []interface{}
	var valmask []string

	data,err := d.ConvertData(org_data)
	if err != nil {
		return 0,false
	}

	for i, v := range data {
		columns = append(columns, d.FormatColumn(i))
		values = append(values, v)
		valmask = append(valmask, "?")
	}

	sql_str := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, strings.Join(columns, ","), strings.Join(valmask, ","))
	res, err := d.Exec(sql_str, values...)
	if err != nil {
		return 0, false
	}

	id, _ := res.LastInsertId()

	return int(id), true
}

//更新数据库数据
func (d *DBA) Update(data DM, where DM, table string) error {
	var values []interface{}
	var tmp []string
	for i, v := range data {
		field := d.explainColumn(i)
		if field.Icon == "+" || field.Icon == "-" {
			tmp = append(tmp, fmt.Sprintf("%s = %s %s ?", field.Field, field.Field, field.Icon))
		} else {
			tmp = append(tmp, fmt.Sprintf("%s %s ?", field.Field, field.Icon))
		}

		values = append(values, v)
	}

	sql_str := fmt.Sprintf("UPDATE %s SET %s", d.FormatColumn(table), strings.Join(tmp, ","))

	if where != nil {
		where_str, where_val := d.WhereRecursion(where, "AND", table)
		values = append(values, where_val...)
		sql_str = fmt.Sprintf("%s WHERE %s", sql_str, where_str)
	}

	_, err := d.Exec(sql_str, values...)
	if err != nil {
		return err
	}

	return nil
}

//条件删除数据
func (d *DBA) Delete(where DM,table string) (int,error){
	var values []interface{}
	where_str, where_val := d.WhereRecursion(where, "AND", table)
	values = append(values, where_val...)
	sql_str := fmt.Sprintf("DELETE FROM %s WHERE %s",table, where_str)
	res, err := d.Exec(sql_str, values...)
	if err != nil {
		return 0,err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return 0,err
	}

	return int(rows),err
}

//查询数据库
func (d *DBA) Query(sql_str string, args ...interface{}) ([]map[string]interface{},error) {
	rows,err := d.db.Query(sql_str, args...)
	d.LastSql = sql_str
	d.LastArgs = args
	if err != nil {
		if d.debug {
			d.HaltError(err)
		}
		return nil,err
	}
	defer rows.Close()

	return d.FetchAll(rows)
}

//执行SQL语句
func (d *DBA) Exec(sql_str string, args ...interface{}) (sql.Result, error) {
	res, err := d.db.Exec(sql_str, args...)

	d.LastSql = sql_str
	d.LastArgs = args

	if err != nil {
		if d.debug {
			d.HaltError(err)
		}
		return nil, err
	}
	return res, err
}

func (d *DBA) QueryOne(sql_str string, args ...interface{}) (map[string]interface{},error) {
	rows,err := d.db.Query(sql_str, args...)
	d.LastSql = sql_str
	d.LastArgs = args
	if err != nil {
		if d.debug {
			d.HaltError(err)
		}
		return nil,err
	}
	defer rows.Close()

	list,err := d.FetchAll(rows)
	if err != nil {
		return nil,err
	}

	if len(list) == 0 {
		return nil,nil
	}

	return list[0],nil
}

func (d *DBA) QueryRow(sql_str string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(sql_str, args...)
}

func (d *DBA) FetchAll(query *sql.Rows) ([]map[string]interface{},error) {
	column, _ := query.Columns()
	values := make([]interface{}, len(column))
	scans := make([]interface{}, len(column))
	for i := range values {
		scans[i] = &values[i]
	}

	results := []map[string]interface{}{}

	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			return nil,err
		}
		row := make(map[string]interface{})
		for k, v := range values {
			key := column[k]
			switch v.(type) {
			case []byte:
				row[key] = string(v.([]byte))
			default:
				row[key] = v
			}

		}

		results = append(results,row)
	}

	return results,nil
}

//处理where条件列表
func (d *DBA) WhereRecursion(where DM, icon string, table string) (string, []interface{}) {
	var where_strings []string
	var values []interface{}
	for i, v := range where {
		if i == "AMD" || i == "OR" {
			tmp_where, val := d.WhereRecursion(v.(DM), i, table)
			where_strings = append(where_strings, tmp_where)
			values = append(values, val...)
		} else {
			vtype := reflect.TypeOf(v).Kind()
			if vtype == reflect.Slice || vtype == reflect.Array {
				values = append(values,v.([]interface{})...)
				//values = append(values, v)
				where_strings = append(where_strings, d.formatWhere(i, table, len(v.([]interface{}))))
			} else {
				values = append(values, v)
				where_strings = append(where_strings, d.formatWhere(i, table, 0))
			}
		}
	}
	where_prefix := fmt.Sprintf(" %s ", icon)
	where_str := fmt.Sprintf("(%s)", strings.Join(where_strings, where_prefix))

	return where_str, values
}

//格式化WHERE条件
func (d *DBA) formatWhere(column string, table string, length int) string {
	field := d.explainColumn(column)

	column_str := d.FormatColumn(field.Field)
	icon := field.Icon

	var format_str string
	if length > 0 {
		var mask_args []string
		where_icon := "IN"
		for i:=0;i<length;i++ {
			mask_args = append(mask_args,"?")
		}
		if icon == "!" {
			where_icon = "NOT IN"
		}
		format_str = fmt.Sprintf("%s.%s %s (%s)", d.FormatColumn(table),column_str,where_icon,strings.Join(mask_args,","))
	} else {
		format_str = fmt.Sprintf("%s.%s %v ?", d.FormatColumn(table), column_str, utils.YN(icon == "!","!=",icon))
	}
	return format_str
}

//解释字段名
func (d *DBA) explainColumn(column string) *DBColumn {
	reg := regexp.MustCompile(`(.+?)\[(\+|-|!|>|<|like)\]`)
	match := reg.FindStringSubmatch(column)
	field := &DBColumn{}
	if len(match) > 0 {
		field.Field = match[1]
		field.Icon = match[2]
	} else {
		field.Field = column
		field.Icon = "="
	}
	return field
}

//错误处理
func (d *DBA) HaltError(err error) {
	fmt.Println(d.LastSql)
	fmt.Println(d.LastArgs)
	fmt.Println(err)
}

//格式化字段名
func (d *DBA) FormatColumn(column string) string {
	return fmt.Sprintf("`%s`", column)
}

func (d *DBA) Close() {
	d.db.Close()
}

func (d *DBA) ConvertData(org_data interface{}) (DM,error) {
	t := reflect.TypeOf(org_data)
	switch t.Kind() {
	case reflect.Map:
		if t.Name() == "DM" {
			return org_data.(DM),nil
		}
		return DM(org_data.(map[string]interface{})),nil
	case reflect.Ptr:
		fallthrough
	case reflect.Struct:
		return DM(utils.Struct2Map(org_data,nil)),nil
	default:
		return nil,errors.New("not support this data")
	}
}

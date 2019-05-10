package ckdb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"reflect"
	"regexp"
	"strings"
)

type DBA struct {
	db             *sql.DB
	table          string
	debug          bool
	LastSql        string
	LastArgs       []interface{}
	LastError      error
	queryInterface interface{}
	tx             *sql.Tx
}

const (
	DB_DRIVER_MYSQL  = "mysql"
	DB_DRIVER_SQLITE = "sqlite"
)

//数据库配置
type DBConfig struct {
	DBDriver   string `json:"db_driver" yaml:"db_driver"`
	DBServer   string `json:"db_server" yaml:"db_server"`
	DBPort     string `json:"db_port" yaml:"db_port"`
	DBName     string `json:"db_name" yaml:"db_name"`
	DBUser     string `json:"db_user" yaml:"db_user"`
	DBPassword string `json:"db_password" yaml:"db_password"`
	DBPoolSize int    `json:"db_pool_size" yaml:"db_pool_size"`
	DBIdleSize int    `json:"db_Idle_size" yaml:"db_Idle_size"`
	DBDebug    bool   `json:"db_debug" yaml:"db_debug"`
}

type DBColumn struct {
	Field string
	Icon  string
}

//reg

var columnReg = regexp.MustCompile(`(.+?)\[(\+|-|!|>|<|<=|>=|like)]`)

//DBA专用数据
type DM map[string]interface{}

func InitDb(conf *DBConfig) (*sql.DB, error) {
	switch conf.DBDriver {
	case DB_DRIVER_MYSQL:
		return InitMysqlDb(conf)
	case DB_DRIVER_SQLITE:
		return InitSqliteDb(conf)
	default:
		return InitMysqlDb(conf)
	}
}

//新创建DBA操作库
func NewDBA(dbConf *DBConfig) (*DBA, error) {
	driver, err := InitDb(dbConf)
	if err != nil {
		return nil, err
	}

	dba := &DBA{db: driver, debug: dbConf.DBDebug}

	return dba, nil
}

//设置操作的表名
func (d *DBA) Table(tableName string) *DBATable {
	return NewDBATable(d, tableName)
}

//开始事务
func (d *DBA) BeginTrans() error {
	var err error
	d.tx, err = d.db.Begin()
	if err != nil {
		return err
	}
	return nil
}

//提交事务
func (d *DBA) Commit() error {
	err := d.tx.Commit()
	if err != nil {
		return err
	}
	d.tx = nil
	return nil
}

//回滚事务
func (d *DBA) Rollback() error {
	err := d.tx.Rollback()
	if err != nil {
		return err
	}
	d.tx = nil
	return nil
}

//插入记录到数据库
func (d *DBA) Insert(table string, orgData interface{}) (int, bool) {
	var columns []string
	var values []interface{}
	var valMask []string

	data, err := d.ConvertData(orgData)
	if err != nil {
		return 0, false
	}

	for i, v := range data {
		columns = append(columns, d.FormatColumn(i))
		values = append(values, v)
		valMask = append(valMask, "?")
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, strings.Join(columns, ","), strings.Join(valMask, ","))
	res, err := d.Exec(sqlStr, values...)
	if err != nil {
		return 0, false
	}

	id, _ := res.LastInsertId()

	return int(id), true
}

//插入多条记录
func (d *DBA) InsertMulti(table string, dataList []interface{}) (int, bool) {
	var columns []string
	var values []interface{}
	var valMask []string
	var keys []string

	for rowIdx, row := range dataList {
		data, err := d.ConvertData(row)
		if err != nil {
			return 0, false
		}
		var mask []string

		if rowIdx == 0 {
			keys = utils.MapKeys(data)
		}

		for _, k := range keys {
			if rowIdx == 0 {
				columns = append(columns, d.FormatColumn(k))
			}
			values = append(values, data[k])
			mask = append(mask, "?")
		}

		valMask = append(valMask, fmt.Sprintf("(%s)", strings.Join(mask, ",")))
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s(%s) VALUES %s", table, strings.Join(columns, ","), strings.Join(valMask, ","))
	res, err := d.Exec(sqlStr, values...)
	if err != nil {
		return 0, false
	}

	rows, _ := res.RowsAffected()
	return int(rows), true
}

//更新数据库数据
func (d *DBA) Update(data utils.M, where utils.M, table string) error {
	var values []interface{}
	var tmp []string
	for i, v := range data {
		field := d.explainColumn(i)
		if field.Icon == "+" || field.Icon == "-" {
			tmp = append(tmp, fmt.Sprintf("%s = %s %s ?", d.FormatColumn(field.Field), d.FormatColumn(field.Field), field.Icon))
		} else {
			tmp = append(tmp, fmt.Sprintf("%s %s ?", d.FormatColumn(field.Field), field.Icon))
		}

		values = append(values, v)
	}

	sqlStr := fmt.Sprintf("UPDATE %s SET %s", d.FormatColumn(table), strings.Join(tmp, ","))

	if where != nil {
		whereStr, whereVal := d.WhereRecursion(where, "AND", table)
		values = append(values, whereVal...)
		sqlStr = fmt.Sprintf("%s WHERE %s", sqlStr, whereStr)
	}

	_, err := d.Exec(sqlStr, values...)
	if err != nil {
		return err
	}

	return nil
}

//条件删除数据
func (d *DBA) Delete(where utils.M, table string) (int, error) {
	var values []interface{}
	whereStr, whereVal := d.WhereRecursion(where, "AND", table)
	values = append(values, whereVal...)
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereStr)
	res, err := d.Exec(sqlStr, values...)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rows), err
}

//查询数据库
func (d *DBA) Query(sqlStr string, args ...interface{}) ([]utils.M, error) {
	rows, err := d.db.Query(sqlStr, args...)
	d.LastSql = sqlStr
	d.LastArgs = args
	if err != nil {
		if d.debug {
			d.HaltError(err)
		}
		return nil, err
	}
	defer rows.Close()

	return d.FetchAll(rows)
}

//查询数据库并返回结果 (传入结构体)
func (d *DBA) QueryStruct(structInterFace interface{}, sqlStr string, args ...interface{}) ([]interface{}, error) {
	rows, err := d.db.Query(sqlStr, args...)
	d.LastSql = sqlStr
	d.LastArgs = args
	if err != nil {
		if d.debug {
			d.HaltError(err)
		}
		return nil, err
	}
	defer rows.Close()

	return d.FetchAllOfStruct(rows, structInterFace)
}

//执行SQL语句
func (d *DBA) Exec(sqlStr string, args ...interface{}) (sql.Result, error) {
	var res sql.Result
	var err error
	if d.tx != nil {
		res, err = d.tx.Exec(sqlStr, args...)
	} else {
		res, err = d.db.Exec(sqlStr, args...)
	}

	d.LastSql = sqlStr
	d.LastArgs = args

	if err != nil {
		if d.debug {
			d.HaltError(err)
		}
		return nil, err
	}
	return res, err
}

//查询一条记录
func (d *DBA) QueryOne(sqlStr string, args ...interface{}) (utils.M, error) {
	list, err := d.Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	//rows, err := d.db.Query(sql_str, args...)
	//d.LastSql = sql_str
	//d.LastArgs = args
	//if err != nil {
	//	if d.debug {
	//		d.HaltError(err)
	//	}
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//list, err := d.FetchAll(rows)
	//if err != nil {
	//	return nil, err
	//}

	if len(list) == 0 {
		return nil, nil
	}

	return list[0], nil
}

//查询一条记录返回结构体
func (d *DBA) QueryOneStruct(structInterFace interface{}, sqlStr string, args ...interface{}) (interface{}, error) {
	list, err := d.QueryStruct(structInterFace, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list[0], nil
}

func (d *DBA) QueryRow(sqlStr string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(sqlStr, args...)
}

//取得所有数据到结构体,没传结构体为默认 utils.M
func (d *DBA) FetchAllOfStruct(query *sql.Rows, i interface{}) ([]interface{}, error) {
	columns, _ := query.Columns()
	scans := make([]interface{}, len(columns))

	var resultList []interface{}

	for query.Next() {
		result := d.scanType(scans, columns, i)
		if err := query.Scan(scans...); err != nil {
			return nil, err
		}
		resultList = append(resultList, result)
	}

	return resultList, nil
}

//取得所有数据到结构体,没传结构体为默认 utils.M
func (d *DBA) FetchAllOfStructV2(query *sql.Rows, i interface{}) ([]interface{}, error) {
	columns, _ := query.Columns()
	scans := make([]interface{}, len(columns))

	var resultList []interface{}

	for query.Next() {
		result := d.scanMap(scans, columns)
		if err := query.Scan(scans...); err != nil {
			return nil, err
		}

		obj := reflect.New(reflect.TypeOf(i)).Interface()

		_ = utils.Map2Struct(result, obj)

		resultList = append(resultList, obj)
	}

	return resultList, nil
}

//取得所有数据
func (d *DBA) FetchAll(query *sql.Rows) ([]utils.M, error) {
	column, _ := query.Columns()
	values := make([]interface{}, len(column))
	scans := make([]interface{}, len(column))
	for i := range values {
		scans[i] = &values[i]
	}

	var results []utils.M

	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			return nil, err
		}
		row := utils.M{}
		for k, v := range values {
			key := column[k]
			switch v.(type) {
			case []byte:
				row[key] = string(v.([]byte))
			default:
				row[key] = v
			}

		}

		results = append(results, row)
	}

	return results, nil
}

//扫描数据到传入的类型
func (d *DBA) scanType(scans []interface{}, columns []string, i interface{}) interface{} {
	if i == nil {
		return d.scanMap(scans, columns)
	}
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.Ptr:
		return d.scanStruct(t.Elem(), scans, columns)
	case reflect.Struct:
		return d.scanStruct(t, scans, columns)
	case reflect.Map:
		fallthrough
	default:
		return d.scanMap(scans, columns)
	}
}

//扫描数据到结构体
func (d *DBA) scanStruct(t reflect.Type, scans []interface{}, columns []string) interface{} {
	obj := reflect.New(t).Interface()
	objV := reflect.ValueOf(obj).Elem()
	for i, colName := range columns {
		idx := d.findTagOfStruct(t, colName)
		if idx != -1 {
			scans[i] = objV.Field(idx).Addr().Interface()
		} else {
			var empty interface{}
			scans[i] = &empty
		}
	}
	return obj
}

//在结构体查找TAG值是否存在
func (d *DBA) findTagOfStruct(t reflect.Type, colName string) int {
	for i := 0; i < t.NumField(); i++ {
		val, ok := t.Field(i).Tag.Lookup("json")
		if ok && val == colName {
			return i
		}
	}
	return -1
}

//扫描数据到MAP 默认 utils.M
func (d *DBA) scanMap(scans []interface{}, columns []string) interface{} {
	obj := utils.M{}
	for i, v := range columns {
		var val interface{}
		obj[v] = val
		scans[i] = &val
	}
	return obj
}

//处理where条件列表
func (d *DBA) WhereRecursion(where utils.M, icon string, table string) (string, []interface{}) {
	var whereStrings []string
	var values []interface{}
	for i, v := range where {
		if i == "AND" || i == "OR" {
			tmpWhere, val := d.WhereRecursion(v.(utils.M), i, table)
			whereStrings = append(whereStrings, tmpWhere)
			values = append(values, val...)
		} else {
			vtype := reflect.TypeOf(v).Kind()
			if vtype == reflect.Slice || vtype == reflect.Array {
				values = append(values, v.([]interface{})...)
				//values = append(values, v)
				whereStrings = append(whereStrings, d.formatWhere(i, table, len(v.([]interface{}))))
			} else {
				values = append(values, v)
				whereStrings = append(whereStrings, d.formatWhere(i, table, 0))
			}
		}
	}
	wherePrefix := fmt.Sprintf(" %s ", icon)
	whereStr := fmt.Sprintf("(%s)", strings.Join(whereStrings, wherePrefix))

	return whereStr, values
}

//格式化WHERE条件
func (d *DBA) formatWhere(column string, table string, length int) string {
	field := d.explainColumn(column)

	columnStr := d.FormatColumn(field.Field)
	icon := field.Icon

	var formatStr string
	if length > 0 {
		var maskArgs []string
		whereIcon := "IN"
		for i := 0; i < length; i++ {
			maskArgs = append(maskArgs, "?")
		}
		if icon == "!" {
			whereIcon = "NOT IN"
		}
		formatStr = fmt.Sprintf("%s.%s %s (%s)", d.FormatColumn(table), columnStr, whereIcon, strings.Join(maskArgs, ","))
	} else {
		formatStr = fmt.Sprintf("%s.%s %v ?", d.FormatColumn(table), columnStr, utils.YN(icon == "!", "!=", icon))
	}
	return formatStr
}

//解释字段名
func (d *DBA) explainColumn(column string) *DBColumn {
	//reg := regexp.MustCompile(`(.+?)\[(\+|-|!|>|<|<=|>=|like)\]`)
	match := columnReg.FindStringSubmatch(column)
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
	d.LastError = err
}

//得到最后一条错误
func (d *DBA) GetLastError() error {
	return d.LastError
}

//格式化字段名
func (d *DBA) FormatColumn(column string) string {
	return fmt.Sprintf("`%s`", column)
}

func (d *DBA) Close() {
	_ = d.db.Close()
}

func (d *DBA) ConvertData(orgData interface{}) (DM, error) {
	t := reflect.TypeOf(orgData)
	switch t.Kind() {
	case reflect.Map:
		if t.Name() == "DM" {
			return orgData.(DM), nil
		} else if t.Name() == "M" {
			return DM(orgData.(utils.M)), nil
		}
		return DM(orgData.(map[string]interface{})), nil
	case reflect.Ptr:
		fallthrough
	case reflect.Struct:
		return DM(utils.Struct2Map(orgData, nil)), nil
	default:
		return nil, errors.New("not support this data")
	}
}

func (d *DBA) SetQueryInterface(i interface{}) {
	d.queryInterface = i
}

//输出表结构为GO struct
func BuildTableStruct(tableName, dbName string, dbConf *DBConfig) {
	types := map[string]string{
		"int":      "int",
		"tinyint":  "int",
		"varchar":  "string",
		"char":     "string",
		"text":     "string",
		"tinytext": "string",
		"double":   "float64",
		"float":    "float64",
		"smallint": "int",
	}

	dba, err := NewDBA(dbConf)
	if err != nil {
		panic(err)
	}

	res, err := dba.Table("COLUMNS").Where(utils.M{"TABLE_NAME": tableName, "TABLE_SCHEMA": dbName}, "").Order(utils.M{"ORDINAL_POSITION": "ASC"}).Query().Result()
	if err != nil {
		panic(err)
	}
	var tmp []string
	var (
		columnName    string
		columnType    string
		columnComment string
	)
	for _, row := range res {
		columnName = row["COLUMN_NAME"].(string)
		columnType = row["DATA_TYPE"].(string)
		columnComment = row["COLUMN_COMMENT"].(string)
		tmp = append(tmp, fmt.Sprintf("\t%v %v `json:\"%v\" bson:\"%v\"` //%v", utils.Under2Hump(columnName), types[columnType], columnName, columnName, columnComment))
	}

	fmt.Println(fmt.Sprintf("type %s struct {", tableName))
	for _, v := range tmp {
		fmt.Println(v)
	}
	fmt.Println("}")
}

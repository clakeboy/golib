package ckdb

import (
	"database/sql"
	"fmt"

	// _ "github.com/go-sql-driver/mysql"
	"time"
)

var MysqlDrivers = make(map[string]*sql.DB)

func InitMysqlDb(conf *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		conf.DBUser,
		conf.DBPassword,
		conf.DBServer,
		conf.DBPort,
		conf.DBName,
	)

	if db, ok := MysqlDrivers[dsn]; ok {
		return db, nil
	}

	MysqlDb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if conf.MaxLife > 0 {
		MysqlDb.SetConnMaxLifetime(time.Duration(conf.MaxLife) * time.Second)
	} else {
		MysqlDb.SetConnMaxLifetime(300 * time.Second)
	}
	MysqlDb.SetMaxOpenConns(conf.DBPoolSize)
	MysqlDb.SetMaxIdleConns(conf.DBIdleSize)

	err = MysqlDb.Ping()
	if err != nil {
		return nil, err
	}

	MysqlDrivers[dsn] = MysqlDb

	return MysqlDb, nil
}

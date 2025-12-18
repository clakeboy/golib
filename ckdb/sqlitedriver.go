package ckdb

import (
	"database/sql"
	"fmt"
	//_ "github.com/mattn/go-sqlite3"
)

var SqliteDrivers = make(map[string]*sql.DB)

func InitSqliteDb(conf *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"file:%s",
		conf.DBName,
	)

	// if db, ok := SqliteDrivers[dsn]; ok {
	// 	return db, nil
	// }

	sqliteDb, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	sqliteDb.SetMaxOpenConns(conf.DBPoolSize)
	sqliteDb.SetMaxIdleConns(conf.DBIdleSize)

	err = sqliteDb.Ping()
	if err != nil {
		return nil, err
	}

	// SqliteDrivers[dsn] = sqliteDb

	return sqliteDb, nil
}

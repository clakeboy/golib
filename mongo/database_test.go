package mongo

import (
	"fmt"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	cfg := &Config{
		Host:     "168.168.0.10",
		Port:     "27017",
		PoolSize: 100,
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	err = db.Open()
	if err != nil {
		t.Error(err)
		return
	}
	err = db.Ping()
	if err != nil {
		t.Error(err)
		return
	}
	list, err := db.ListDatabase()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(list)
}

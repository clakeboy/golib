package components

import (
	"fmt"
	"testing"
	"time"

	"github.com/clakeboy/golib/utils"
)

func TestFindKeys(t *testing.T) {
	mem := NewMemCache()
	for i := 0; i < 10; i++ {
		mem.Set(fmt.Sprintf("ck_df_%d", i), i, -1)
	}

	keys, _ := mem.Keys("ck_df")
	utils.PrintAny(keys)

	mem.LoadLocal("test_store")
	fmt.Println("停止10秒等待写入")
	time.Sleep(10 * time.Second)
}

func TestLoad(t *testing.T) {
	mem := NewMemCache()
	mem.LoadLocal("test_store")
	ss, err := mem.Get("ck_df_0")
	fmt.Println("read key:", "ck_df_0", ss, err)
}

func TestTrick(t *testing.T) {
	ti := time.NewTicker(time.Second)
	for {
		select {
		case t := <-ti.C:
			fmt.Println(t.UnixMicro())
		}
	}
}

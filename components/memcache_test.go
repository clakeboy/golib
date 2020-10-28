package components

import (
	"fmt"
	"github.com/clakeboy/golib/utils"
	"testing"
)

func TestFindKeys(t *testing.T) {
	mem := NewMemCache()
	for i := 0; i < 10; i++ {
		mem.Set(fmt.Sprintf("ck_df_%d", i), i, -1)
	}

	keys, _ := mem.Keys("ck_df")
	utils.PrintAny(keys)
}

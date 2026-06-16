package components

import (
	"fmt"
	"testing"

	"github.com/clakeboy/golib/utils"
)

var logger, _ = NewSlogFile("default.log")

func TestMain(t *testing.T) {
	fmt.Println(utils.RandStr(16, nil))
	logger.Info("this is info message", "key", "value", "randstr", utils.RandStr(16, nil))
	logger.Error("this is error message", "key", "value", "randstr", utils.RandStr(16, nil))
}

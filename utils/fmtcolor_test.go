package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestFmtColor(t *testing.T) {
	fmt.Println(FmtColor("clake is good!", FRED, BRED, SHGH))
	fmt.Fprintln(os.Stdout, FmtColor("clake is good!", FRED, BRED, SHGH))

}

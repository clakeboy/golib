package utils

import (
	"testing"
	"fmt"
	"os"
)

func TestFmtColor(t *testing.T) {
	fmt.Println(FmtColor("clake is good!",FRED,BRED,SHGH))
	fmt.Fprintln(os.Stdout,FmtColor("clake is good!",FRED,BRED,SHGH))

}
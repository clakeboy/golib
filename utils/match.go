package utils

import (
	"fmt"
	"strconv"
)

//保留小数位
func Round(d float64, position int) float64 {
	formatStr := "%." + strconv.Itoa(position) + "f"
	s := fmt.Sprintf(formatStr, d)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

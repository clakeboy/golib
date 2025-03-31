package utils

import (
	"time"
	//"fmt"
	"strconv"
	"strings"
)

var month = map[string]string{
	"January":   "1",
	"February":  "2",
	"March":     "3",
	"April":     "4",
	"May":       "5",
	"June":      "6",
	"July":      "7",
	"August":    "8",
	"September": "9",
	"October":   "10",
	"November":  "11",
	"December":  "12"}

func FormatTime(format_str string, unix_timestamp int64) string {
	t := time.Unix(unix_timestamp, 0)

	format_field := &map[string]string{
		"YY": strconv.Itoa(t.Year()),
		"MM": StrPad(month[t.Month().String()], "0", 2, STR_PAD_LEFT),
		"DD": StrPad(strconv.Itoa(t.Day()), "0", 2, STR_PAD_LEFT),
		"yy": strconv.Itoa(t.Year())[2:],
		"mm": month[t.Month().String()],
		"dd": strconv.Itoa(t.Day()),
		"HH": StrPad(strconv.Itoa(t.Hour()), "0", 2, STR_PAD_LEFT),
		"II": StrPad(strconv.Itoa(t.Minute()), "0", 2, STR_PAD_LEFT),
		"SS": StrPad(strconv.Itoa(t.Second()), "0", 2, STR_PAD_LEFT)}

	for k, v := range *format_field {
		format_str = strings.Replace(format_str, k, v, -1)
	}
	return format_str
}

const (
	STR_PAD_LEFT = 1 + iota
	STR_PAD_RIGHT
)

func StrPad(str string, pad_str string, number int, t int) string {
	if str_leng := len(str); str_leng < number {
		diff := number - str_leng
		for i := 0; i < diff; i++ {
			if t == STR_PAD_LEFT {
				str = pad_str + str
			} else if t == STR_PAD_RIGHT {
				str += pad_str
			}
		}
	}

	return str
}

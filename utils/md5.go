package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodeMD5(s string) string {
	h := md5.New()
	h.Write([]byte("uu<8"+s+"end*u^3"))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

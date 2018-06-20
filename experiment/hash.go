package experiment

import "unicode"

const MAXMPQHASHTABLELEN = 8192

type MPQHashTable struct {
	nHashA  int64
	nHashB  int64
	bExists uint
}

var cryptTable [0x500]uint64

func InitCryptTable() {
	var seed, index1, index2 uint64 = 0x00100001, 0, 0
	i := 0

	for index1 = 0; index1 < 0x100; index1++ {
		for index2, i = index1, 0; i < 5; i, index2 = i+1, index2+0x100 {
			var tmp1, tmp2 uint64
			seed = (seed*125 + 3) % 0x2AAAAB
			tmp1 = (seed & 0xFFFF) << 0x10
			seed = (seed*125 + 3) % 0x2AAAAB
			tmp2 = seed & 0xFFFF
			cryptTable[index2] = tmp1 | tmp2
		}
	}
}

func HashString(lpszString string,dwHashType int) uint64 {
	var key uint8
	var seed1,seed2 uint64 = 0x7FED7FED,0xEEEEEEEE
	strLen := len(lpszString)
	i,ch := 0,0
	for i< strLen {
		key = lpszString[i]
		ch = int(unicode.ToUpper(rune(key)))
		seed1 = cryptTable[(dwHashType<<8)+ch] ^ (seed1 + seed2)
		seed2 = uint64(ch) + seed1 + seed2 + (seed2 << 5) + 3
		i++
	}

	return seed1
}

func MPQHashTableInit(ppHashTable []*MPQHashTable, nTableLength int64) {
	InitCryptTable()
}



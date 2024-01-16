package components

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"runtime"
	"sort"
	"testing"
	"unsafe"

	"github.com/clakeboy/golib/utils"
)

func TestNewBTreePlus(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ex := utils.NewExecTime()
	b := NewBTreePlus(128)
	ex.Start()
	for i := 0; i < 10000000; i++ {
		b.Insert(i, fmt.Sprintf("data:%d", i+1))
	}
	ex.End(true)
	// b.Print()
	ex.Start()
	for i := 0; i < 100000; i++ {
		// fmt.Println(b.Search(5))
		b.Search(i)
	}
	ex.End(true)
	//var arr []int
	//for i:= 0;i<100000000;i++ {
	//	arr = append(arr,i+1)
	//}
	//size := unsafe.Sizeof(0)
	//fmt.Println("pass 10 second",size*100000000/1024/1024)
	//time.Sleep(time.Second * 10)
	//b.Print()
}

func TestBTreePlus_Insert(t *testing.T) {
	arr := []int{2, 5, 4, 3, 7, 1, 9}
	//arr := []int{1,2,3,4,5,6,7,8,9}
	b := NewBTreePlus(2)
	for _, v := range arr {
		b.Insert(v, fmt.Sprintf("data:%d", v))
		b.Print()
	}

	fmt.Println(b.Search(6))
	fmt.Println(b.Search(3))
	fmt.Println(b.Search(5))
	fmt.Println(b.Search(8))
}

func TestBTreePlus_Print(t *testing.T) {
	//index := 2
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	index := sort.Search(len(arr), func(i int) bool {
		return arr[i] == 5
	})
	fmt.Println(index)
	//start := arr[:index]
	//end := arr[index+1:]
	//start = append(start,10)
	//arr = append(start,end...)
	//fmt.Println(start,end,arr)
	ss := 0xff
	fmt.Println(unsafe.Sizeof(arr[1]))
	fmt.Println(unsafe.Sizeof(ss))
	fmt.Println(unsafe.Sizeof(TreeNode{}))
	fmt.Println(arr[:5:7])
}

func TestBTreePlus_Search(t *testing.T) {
	index := 1
	keys := KeysInt{1, 2, 4, 5, 6, 7, 8, 9, 0}
	arr := []int{2, 5}
	fmt.Println(arr[:index], arr[index:])
	fmt.Println(arr[len(arr):])
	fmt.Println(keys.Pop(), keys)

	var kk [][]byte
	var ki [][]byte
	kk = append(kk, []byte("calkesssbbbcccksidke"), []byte("dfeeeee"))
	ki = append(ki, utils.IntToBytes(5, 64), utils.IntToBytes(10, 64))

	fmt.Println(bytes.Compare(kk[0], kk[1]))
	fmt.Println(bytes.Compare(ki[1], ki[0]))
}

func TestBTreePlus_Remove(t *testing.T) {
	a := []string{"hello", "", "world", "yes", "hello", "nihao", "shijie", "hello", "yes", "nihao", "good"}
	sort.Strings(a)
	fmt.Println(a)
	fmt.Println(RemoveDuplicatesAndEmpty(a))
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func TestEnBytes(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode("clakebfdds33332")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%X,%d", buf.Bytes(), buf.Len())
}

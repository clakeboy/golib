package snowflake

import (
	"testing"
	"fmt"
	"time"
)

func TestNewShowFlake(t *testing.T) {
	fmt.Println(^(-1<<41))

	fmt.Println(time.Now().UnixNano()/1e6)

	fmt.Println()

	//90117296928067856 90117351416270848 90121840378515943
	//222408272404287488 222408382316023808
	//6443400298957639680
	//9223372036854775807
}

func TestSnowFlake_NextId(t *testing.T) {
	snow,err := NewShowFlake(1514739661000,1,2)
	if err != nil {
		panic(err)
	}
	for i:=0;i<1e3;i++ {
		id,err := snow.NextId()
		if err != nil {
			fmt.Println("generate id error:",err)
			break
		}
		fmt.Println(id.RawId)
	}
}
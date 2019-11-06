package components

import (
	"fmt"
	"testing"
)

func TestNewPoll(t *testing.T) {
	out := make(chan int, 4)
	fmt.Printf("%v\n", out)
	close(out)
	fmt.Printf("%v", out)
}

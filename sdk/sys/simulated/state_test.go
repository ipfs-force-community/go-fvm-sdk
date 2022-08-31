package simulated


import (
	"fmt"
	"testing"
)

func Test_state_block_read(t *testing.T) {
	s1 := make([]int, 1, 10)
	sss((&s1))
	fmt.Printf("%v\n", s1)
}

func sss(si *[]int) {
	*si = append(*si, 77)
}

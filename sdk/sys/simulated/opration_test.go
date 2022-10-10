package simulated

import (
	"fmt"
	"testing"
)

func TestFsm_NewActorAddress(t *testing.T) {
	s := Fsm{}
	got, _ := s.NewActorAddress()
	fmt.Printf("%s\n", got.String())

}

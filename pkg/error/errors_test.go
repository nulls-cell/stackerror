package stackerror

import (
	"fmt"
	"testing"
)

func TestRuntimeError(t *testing.T) {
	SetStopKeyWordsByString("local")
	err := NewStackError("a test error")
	err2 := WrapStackError(err)
	err3 := WrapStackError(err2)
	err4 := NewStackError("normal error")
	err4 = NewStackErrorf("normal error, args: %s", "arg1")
	fmt.Println(err3.Error(), err3.GetStack())
	fmt.Println(err4.Error(), err4.GetStack())
}

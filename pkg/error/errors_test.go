package stackerror

import (
	"fmt"
	"testing"
)

func TestRuntimeError(t *testing.T) {
	SetStopKeyWordsByString("local")
	err := f2()
	err2 := Wrap(err)
	err3 := Wrap(err2)
	err4 := NewError("normal error")
	err4 = NewErrorf("normal error, args: %s", "arg1")
	fmt.Println(err3.Error(), err3.GetDepthStackErrMsg(1))
	fmt.Println(err3.Error(), err3.GetDepthStackErrMsg(2))
	fmt.Println(err3.Error(), err3.GetDepthStackErrMsg(3))
	fmt.Println(err4.Error(), err4.GetStack())
}

func f1() IStackError {
	err := NewError("a test error")
	return err
}

func f2() IStackError {
	err := f1()
	return err
}

# stackerror

## example
```go
package main

import "nulls-cell/stackerror/stackerror"

func TestRuntimeError(t *testing.T) {
	stackerror.SetStopKeyWordsByString("local")
	err := f2()
	err2 := stackerror.Wrap(err)
	err3 := stackerror.Wrap(err2)
	err4 := stackerror.NewError("normal error")
	err4 = stackerror.NewErrorf("normal error, args: %s", "arg1")
	fmt.Println(err3.Error(), err3.GetDepthStackErrMsg(1))
	fmt.Println(err3.Error(), err3.GetDepthStackErrMsg(2))
	fmt.Println(err3.Error(), err3.GetDepthStackErrMsg(3))
	fmt.Println(err4.Error(), err4.GetStack())
}

func f1() stackerror.IStackError {
	err := NewError("a test error")
	return err
}

func f2() stackerror.IStackError {
	err := f1()
	return err
}

```
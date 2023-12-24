package stackerror

import (
	"bytes"
	"fmt"
	"runtime/debug"
)

type StackError struct {
	msg          string
	runtimeStack string
	OrgError     error
}

var stopKeyWords []byte
var rTrimCutSet = string([]byte{byte(0), '\n'})

func SetStopKeyWords(bts []byte) {
	stopKeyWords = bts
}

func SetStopKeyWordsByString(s string) {
	stopKeyWords = []byte(s)
}

func (r *StackError) Error() string {
	return r.msg
}

func (r *StackError) GetStack() string {
	return r.runtimeStack
}

func (r *StackError) GetStackErrMsg() string {
	return r.msg + "\n" + r.runtimeStack
}

func NewStackError(msg string) *StackError {
	stack := getStackWithoutString(debug.Stack(), stopKeyWords)
	err := &StackError{
		msg:          msg,
		runtimeStack: stack,
		OrgError:     nil,
	}
	return err
}

func NewStackErrorf(format string, args ...interface{}) *StackError {
	msg := fmt.Sprintf(format, args...)
	stack := getStackWithoutString(debug.Stack(), stopKeyWords)
	err := &StackError{
		msg:          msg,
		runtimeStack: stack,
		OrgError:     nil,
	}
	return err
}

func getStackWithoutString(stack []byte, stop []byte) string {
	var firstLineEndIndex, startIndex, endIndex int
	lines := bytes.Split(stack, []byte{'\n'})
	for i, line := range lines {
		lineLen := len(line)
		if i == 0 {
			firstLineEndIndex = lineLen + 1
		}
		if i <= 4 {
			startIndex += lineLen + 1
		} else if len(stop) > 0 && bytes.Contains(line, stop) {
			break
		}
		endIndex += lineLen + 1
	}
	firstLine := string(stack[:firstLineEndIndex])
	body := bytes.TrimRight(stack[startIndex:endIndex], rTrimCutSet)
	return firstLine + string(body)
}

func WrapStackError(unknownErr error) *StackError {
	if unknownErr == nil {
		return nil
	}
	errObj, ok := unknownErr.(*StackError)
	if ok {
		return errObj
	}
	stack := getStackWithoutString(debug.Stack(), stopKeyWords)

	newErr := &StackError{
		msg:          unknownErr.Error(),
		runtimeStack: stack,
		OrgError:     unknownErr,
	}
	return newErr
}

type IStackError interface {
	Error() string
	GetStack() string
	GetStackErrMsg() string
}

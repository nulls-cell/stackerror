package stackerror

import (
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
)

type StackError struct {
	msg               string
	runtimeStackBytes []byte
	OrgError          error
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
	return getStackByBytes(r.runtimeStackBytes, stopKeyWords)
}

func (r *StackError) GetStackErrMsg() string {
	return r.msg + "\n" + getStackByBytes(r.runtimeStackBytes, stopKeyWords)
}

func (r *StackError) GetDepthStackErrMsg(depth int32) string {
	return r.msg + "\n" + getDepthStackByBytes(r.runtimeStackBytes, stopKeyWords, depth)
}

func Error(msg string) *StackError {
	stack := debug.Stack()
	err := &StackError{
		msg:               msg,
		runtimeStackBytes: stack,
		OrgError:          nil,
	}
	return err
}

func Errorf(format string, args ...interface{}) *StackError {
	msg := fmt.Sprintf(format, args...)
	stack := debug.Stack()
	err := &StackError{
		msg:               msg,
		runtimeStackBytes: stack,
		OrgError:          nil,
	}
	return err
}

func getStackByBytes(stack []byte, stop []byte) string {
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

func getDepthStackByBytes(stack []byte, stop []byte, depth int32) string {
	var firstLineEndIndex, startIndex, endIndex int
	if depth <= 0 {
		depth = 1 << 30
	} else {
		depth *= 2
		depth += 5
	}
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
		depth -= 1
		if depth == 0 {
			break
		}
	}
	firstLine := string(stack[:firstLineEndIndex])
	body := bytes.TrimRight(stack[startIndex:endIndex], rTrimCutSet)
	return firstLine + string(body)
}

func Wrap(unknownErr error) IError {
	if unknownErr == nil {
		return nil
	}
	var errObj *StackError
	ok := errors.As(unknownErr, &errObj)
	if ok {
		return errObj
	}

	newErr := &StackError{
		msg:               unknownErr.Error(),
		runtimeStackBytes: debug.Stack(),
		OrgError:          unknownErr,
	}
	return newErr
}

type IError interface {
	Error() string
	GetStack() string
	GetStackErrMsg() string
	GetDepthStackErrMsg(depth int32) string
}

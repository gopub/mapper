package goparam

import "fmt"

type Error struct {
	paramName string
	msg       string
}

func newError(paramName, msg string) *Error  {
	return &Error{
		paramName:paramName,
		msg:msg,
	}
}

func (e *Error) ParamName() string  {
	return e.paramName
}

func (e *Error) Message() string  {
	return e.msg
}

func (e *Error) Error() string {
	return fmt.Sprintf("param:%s, msg:%s", e.paramName, e.msg)
}


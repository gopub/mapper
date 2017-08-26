package param

import "fmt"

type Error interface {
	error
	ParamName() string
	Message() string
}

type paramError struct {
	paramName string
	msg       string
}

var _ Error = (*paramError)(nil)

func newError(paramName, msg string) *paramError {
	return &paramError{
		paramName: paramName,
		msg:       msg,
	}
}

func (e *paramError) ParamName() string {
	return e.paramName
}

func (e *paramError) Message() string {
	return e.msg
}

func (e *paramError) Error() string {
	return fmt.Sprintf("param:%s, msg:%s", e.paramName, e.msg)
}

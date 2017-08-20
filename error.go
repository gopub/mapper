package goparam

import "fmt"

type Error struct {
	ParamName string
	Message   string
}

func (e *Error) Error() string {
	return fmt.Sprintf("param:%s, msg:%s", e.ParamName, e.Message)
}

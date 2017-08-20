package goparam

import (
	"regexp"
)

var _nameToRegexp = map[string]*regexp.Regexp{
	"version":  regexp.MustCompile(RegexpVersion),
	"mobile":   regexp.MustCompile(RegexpMobile_ZHCN),
	"email":    regexp.MustCompile(RegexpEmail),
	"variable": regexp.MustCompile(RegexpVariable),
}

const (
	RegexpVersion     = "^[1-9]\\d*(\\.\\d*)*([\\w-]+)?$"
	RegexpMobile_ZHCN = "^1[34578][0-9]{9}$"
	RegexpEmail       = "^[_a-zA-Z0-9-]+(\\.[_a-zA-Z0-9-]+)*@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*(\\.[a-zA-Z]{2,3})$"
	RegexpVariable    = "^[_a-zA-Z][_a-zA-Z0-9]*$"
)

const (
	TypeVersion  = "version"
	TypeMobile   = "mobile"
	TypeEmail    = "email"
	TypeVariable = "variable"
)

func MatchPattern(v string, patternName string) bool {
	return false
}

package goparam

import (
	"regexp"
)

type PatternMatcher interface {
	Match(i interface{}) bool
}

type PatternMatchFunc func(i interface{}) bool

func (p PatternMatchFunc) Match(i interface{}) bool {
	return p(i)
}

type Regexp regexp.Regexp

func (r *Regexp) Match(i interface{}) bool {
	s, ok := i.(string)
	if !ok {
		return false
	}

	return (*regexp.Regexp)(r).MatchString(s)
}

var _patternToMatcher = map[string]PatternMatcher{
	"version":  (*Regexp)(regexp.MustCompile(RegexpVersion)),
	"mobile":   (*Regexp)(regexp.MustCompile(RegexpMobile_ZHCN)),
	"email":    (*Regexp)(regexp.MustCompile(RegexpEmail)),
	"variable": (*Regexp)(regexp.MustCompile(RegexpVariable)),
}

const (
	RegexpVersion     = "^[1-9]\\d*(\\.\\d*)*([\\w-]+)?$"
	RegexpMobile_ZHCN = "^1[34578][0-9]{9}$"
	RegexpEmail       = "^[_a-zA-Z0-9-]+(\\.[_a-zA-Z0-9-]+)*@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*(\\.[a-zA-Z]{2,3})$"
	RegexpVariable    = "^[_a-zA-Z][_a-zA-Z0-9]*$"
)

func MatchPattern(name string, v interface{}) bool {
	return _patternToMatcher[name].Match(v)
}

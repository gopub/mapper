package param

import (
	"net/url"
	"regexp"
	"time"
)

const (
	RegexpVersion     = "^[1-9]\\d*(\\.\\d*)*([\\w-]+)?$"
	RegexpMobile_ZHCN = "^1[34578][0-9]{9}$"
	RegexpEmail       = "^[_a-zA-Z0-9-]+(\\.[_a-zA-Z0-9-]+)*@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*(\\.[a-zA-Z]{2,3})$"
	RegexpVariable    = "^[_a-zA-Z][_a-zA-Z0-9]*$"
)

const (
	PatternVersion   = "version"
	PatternURL       = "url"
	PatternEmail     = "email"
	PatternVariable  = "variable"
	PatternMobile    = "mobile"
	PatternBirthDate = "birth_date"
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
	PatternVersion:   (*Regexp)(regexp.MustCompile(RegexpVersion)),
	PatternMobile:    (*Regexp)(regexp.MustCompile(RegexpMobile_ZHCN)),
	PatternEmail:     (*Regexp)(regexp.MustCompile(RegexpEmail)),
	PatternVariable:  (*Regexp)(regexp.MustCompile(RegexpVariable)),
	PatternBirthDate: PatternMatchFunc(MatchBirthDate),
	PatternURL:       PatternMatchFunc(MatchURL),
}

func MatchPattern(name string, v interface{}) bool {
	return _patternToMatcher[name].Match(v)
}

func MatchBirthDate(i interface{}) bool {
	s, ok := i.(string)
	if !ok {
		return false
	}

	t, e := time.Parse("2006-01-02", s)
	if e != nil {
		return false
	}

	if t.After(time.Now()) {
		return false
	}

	return true
}

func MatchURL(i interface{}) bool {
	s, ok := i.(string)
	if !ok {
		return false
	}

	if u, err := url.Parse(s); err != nil {
		return false
	} else if len(u.Scheme) == 0 {
		return false
	} else if len(u.Host) == 0 {
		return false
	}

	return true
}

package goparam

import "regexp"

type Validator struct {
	tagName      string
	nameToRegexp map[string]*regexp.Regexp
}

func NewValidator(tagName string) *Validator {
	v := &Validator{}
	v.nameToRegexp = make(map[string]*regexp.Regexp, len(_nameToRegexp))
	for t, r := range _nameToRegexp {
		v.nameToRegexp[t] = r
	}
	return v
}

func (v *Validator) RegisterPattern(name, pattern string) error {
	return nil
}

func (v *Validator) Validate() error {
	return nil
}

var _defaultValidator = NewValidator("param")

func RegisterPattern(name, pattern string) error {
	return _defaultValidator.RegisterPattern(name, pattern)
}

func Validate(model interface{}) error {
	return _defaultValidator.Validate(model)
}

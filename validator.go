package goparam

type Validator struct {
	tagName           string
	patternToMatcher  map[string]PatternMatcher
	nameToTransformer map[string]Transformer
}

func NewValidator(tagName string) *Validator {
	v := &Validator{}
	v.patternToMatcher = make(map[string]PatternMatcher, len(_patternToMatcher))
	for t, r := range _patternToMatcher {
		v.patternToMatcher[t] = r
	}
	v.nameToTransformer = make(map[string]Transformer, len(_nameToTransformer))
	for n, t := range _nameToTransformer {
		v.nameToTransformer[n] = t
	}
	return v
}

func (v *Validator) RegisterPatternMatcher(name string, matcher PatternMatcher) error {
	return nil
}

func (v *Validator) RegisterPatternMatchFunc(name string, matchFunc PatternMatchFunc) error {
	return nil
}

func (v *Validator) Validate(model interface{}) error {
	return nil
}

var _defaultValidator = NewValidator("param")

func RegisterPatternMatcher(name string, matcher PatternMatcher) error {
	return _defaultValidator.RegisterPatternMatcher(name, matcher)
}

func RegisterPatternMatchFunc(name string, matchFunc PatternMatchFunc) error {
	return _defaultValidator.RegisterPatternMatchFunc(name, matchFunc)
}

func Validate(model interface{}) error {
	return _defaultValidator.Validate(model)
}

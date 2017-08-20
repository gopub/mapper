package goparam

import (
	"fmt"
	"reflect"
	"sync"
)

type Validator struct {
	tagName           string
	patternToMatcher  map[string]PatternMatcher
	nameToTransformer map[string]Transformer
	typeToModelInfo   sync.Map
}

func NewValidator(tagName string) *Validator {
	v := &Validator{
		tagName: tagName,
	}
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

func (v *Validator) RegisterRegexpPattern(name string, regularExpression string) error {
	return nil
}

func (v *Validator) Validate(model interface{}) error {
	val := reflect.ValueOf(model)
	if val.IsValid() == false {
		panic("not valid")
	}

	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic("not struct")
	}

	info := v.getModelInfo(val.Type())

	for i, pi := range info {
		fmt.Printf("%d min=%v max=%v pattern=%s transformer=%s optional=%v\n", i, pi.minVal, pi.maxVal, pi.patternName, pi.transformName, pi.optional)
	}

	for i := 0; i < val.NumField(); i++ {
		fv := val.Field(i)
		ft := val.Type().Field(i)
		if ft.Anonymous {
			err := v.Validate(fv.Interface())
			if err != nil {
				return err
			}
			continue
		}

		pi, ok := info[i]
		if !ok {
			continue
		}

		if fv.Interface() == reflect.Zero(fv.Type()).Interface() && pi.optional {
			continue
		}

		if pi.minVal != nil {
			switch fv.Kind() {
			case reflect.Float32, reflect.Float64:
				if fv.Float() < pi.minVal.(float64) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("larger than %v", pi.minVal),
					}
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if fv.Int() < pi.minVal.(int64) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("larger than %v", pi.minVal),
					}
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if fv.Uint() < pi.minVal.(uint64) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("larger than %v", pi.minVal),
					}
				}
			case reflect.String:
				i := pi.minVal.(int64)
				if len(fv.String()) < int(i) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("longer than %v", pi.minVal),
					}
				}
			default:
				panic("invalid kind: " + fv.Kind().String())
			}
		}

		if pi.maxVal != nil {
			switch fv.Kind() {
			case reflect.Float32, reflect.Float64:
				if fv.Float() > pi.maxVal.(float64) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("less than %v", pi.maxVal),
					}
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if fv.Int() > pi.maxVal.(int64) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("less than %v", pi.maxVal),
					}
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if fv.Uint() > pi.maxVal.(uint64) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("less than %v", pi.maxVal),
					}
				}
			case reflect.String:
				i := pi.maxVal.(int64)
				if len(fv.String()) > int(i) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("shorter than %v", pi.maxVal),
					}
				}
			default:
				panic("invalid kind: " + fv.Kind().String())
			}
		}

		if len(pi.patternName) > 0 {
			if matcher, ok := v.patternToMatcher[pi.patternName]; ok {
				if !matcher.Match(fv.Interface()) {
					return &Error{
						ParamName: pi.name,
						Message:   fmt.Sprintf("not match pattern: %s", pi.patternName),
					}
				}
			} else {
				panic("invalid pattern: " + pi.patternName)
			}
		}
	}
	return nil
}

func (v *Validator) getModelInfo(modelType reflect.Type) modelInfo {
	if infoVal, ok := v.typeToModelInfo.Load(modelType); ok {
		return infoVal.(modelInfo)
	}

	info, err := parseModelInfo(modelType, v.tagName)
	if err != nil {
		panic(err.Error())
	}
	v.typeToModelInfo.Store(modelType, info)
	return info
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

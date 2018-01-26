package mapper

import (
	"fmt"
	"github.com/gopub/log"
	"github.com/gopub/types"
	"reflect"
)

// Assign assigns params to model with DefaultValidator
func Assign(model interface{}, params interface{}) Error {
	return AssignWithValidator(model, params, _defaultValidator)
}

func AssignWithTag(model interface{}, params interface{}, tagName string) Error {
	var v *Validator
	if tagName == _defaultValidator.tagName {
		v = _defaultValidator
	} else {
		v = NewValidator(tagName)
	}
	return AssignWithValidator(model, params, v)
}

func AssignWithValidator(model interface{}, params interface{}, validator *Validator) Error {
	v := reflect.ValueOf(model)
	if v.IsValid() == false {
		panic("not valid")
	}

	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	//v must be a nil pointer or a valid value
	err := assignValue(v, reflect.ValueOf(params), validator)
	if err != nil {
		return err
	}
	return validator.Validate(model)
}

// dstVal is valid value or pointer to value
func assignValue(dstVal reflect.Value, srcVal reflect.Value, validator *Validator) *paramError {
	if !srcVal.IsValid() {
		return newError("", "no source value")
	}

	if !dstVal.IsValid() {
		log.Panicf("invalid values:dstVal=%v,srcVal=%v", dstVal, srcVal)
	}

	v := dstVal
	if v.Kind() == reflect.Ptr {
		if v.IsNil() && v.CanSet() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	if !v.CanSet() {
		panic("can't set")
	}

	switch v.Kind() {
	case reflect.Bool:
		b, err := types.ParseBool(srcVal.Interface())
		if err != nil {
			return newError("", err.Error())
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := types.ParseInt(srcVal.Interface())
		if err != nil {
			return newError("", err.Error())
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := types.ParseInt(srcVal.Interface())
		if err != nil {
			return newError("", err.Error())
		}
		v.SetUint(uint64(i))
	case reflect.Float32, reflect.Float64:
		i, err := types.ParseFloat(srcVal.Interface())
		if err != nil {
			return newError("", err.Error())
		}
		v.SetFloat(i)
	case reflect.String:
		if srcVal.Kind() != reflect.String {
			return newError("", "source value is not string")
		}
		v.SetString(srcVal.String())
	case reflect.Slice:
		if srcVal.Kind() != reflect.Slice {
			return newError("", "source value is slice")
		}
		v.Set(reflect.MakeSlice(v.Type(), srcVal.Len(), srcVal.Cap()))
		for i := 0; i < srcVal.Len(); i++ {
			assignValue(v.Index(i), srcVal.Index(i), validator)
		}
	case reflect.Map:
		err := assignMap(v, srcVal, validator)
		if err != nil {
			return err
		}
	case reflect.Struct:
		err := assignStruct(v, srcVal, validator)
		if err != nil {
			return err
		}
	default:
		panic("unknown kind: " + v.Kind().String())
	}

	if dstVal.Kind() == reflect.Ptr && dstVal.IsNil() {
		dstVal.Set(v.Addr())
	}
	return nil
}

// dstVal is map
// srcVal is map
func assignMap(dstVal reflect.Value, srcVal reflect.Value, validator *Validator) *paramError {
	if dstVal.Kind() != reflect.Map {
		panic("not map")
	}

	if srcVal.Kind() != reflect.Map {
		return newError("", "srcVal is not map")
	}

	if !srcVal.Type().Key().AssignableTo(dstVal.Type().Key()) {
		msg := fmt.Sprintf("%s can't be assigned to %s", srcVal.Type().Key().String(), dstVal.Type().Key().String())
		return newError("", msg)
	}

	if dstVal.IsNil() {
		dstVal.Set(reflect.MakeMap(dstVal.Type()))
	}

	de := dstVal.Type().Elem()
	canAssign := srcVal.Type().Elem().AssignableTo(de)
	for _, k := range srcVal.MapKeys() {
		switch {
		case canAssign:
			dstVal.SetMapIndex(k, srcVal.MapIndex(k))
		case de.Kind() == reflect.Ptr:
			kv := reflect.New(de.Elem())
			err := assignValue(kv, srcVal.MapIndex(k), validator)
			if err != nil {
				return err
			}
			dstVal.SetMapIndex(k, kv)
		default:
			kv := reflect.New(de)
			err := assignValue(kv, srcVal.MapIndex(k), validator)
			if err != nil {
				return err
			}
			dstVal.SetMapIndex(k, kv.Elem())
		}
	}

	return nil
}

// dstVal is struct
// srcVal is map
func assignStruct(dstVal reflect.Value, srcVal reflect.Value, validator *Validator) *paramError {
	if dstVal.Kind() != reflect.Struct {
		panic("not struct")
	}

	if srcVal.Kind() == reflect.Interface {
		srcVal = srcVal.Elem()
	}

	if srcVal.Kind() != reflect.Map {
		return &paramError{
			msg: "srcVal is not map",
		}
	}

	if srcVal.Type().Key().Kind() != reflect.String {
		return &paramError{
			msg: "key type must be string",
		}
	}

	info := validator.getModelInfo(dstVal.Type())
	for i := 0; i < dstVal.NumField(); i++ {
		fv := dstVal.Field(i)
		if fv.IsValid() == false || fv.CanSet() == false {
			continue
		}

		ft := dstVal.Type().Field(i)
		if ft.Anonymous {
			err := assignValue(fv, srcVal, validator)
			if err != nil {
				return err
			}
			continue
		}

		pi, ok := info[i]
		if !ok {
			continue
		}

		fsv := srcVal.MapIndex(reflect.ValueOf(pi.name))
		if fsv.IsValid() {
			err := assignValue(fv, reflect.ValueOf(fsv.Interface()), validator)
			if err != nil {
				if len(err.paramName) > 0 {
					err.paramName = pi.name + "." + err.paramName
				} else {
					err.paramName = pi.name
				}
				return err
			}
		} else if !pi.optional {
			return &paramError{
				paramName: pi.name,
				msg:       "no value",
			}
		}
	}
	return nil
}

package mapper

import (
	"errors"
	"fmt"
	"github.com/gopub/log"
	"github.com/gopub/types"
	"reflect"
)

func Assign(dst interface{}, src interface{}) error {
	return AssignWithNameMapper(dst, src, _defaultNameMapper)
}

// Assign assigns params to model with DefaultValidator
func AssignWithNameMapper(dst interface{}, src interface{}, nameMapper NameMapper) error {
	dstVal := reflect.ValueOf(dst)
	if dstVal.IsValid() == false {
		log.Panic("dst value is invalid")
	}

	for dstVal.Kind() == reflect.Ptr && !dstVal.IsNil() {
		dstVal = dstVal.Elem()
	}

	if nameMapper == nil {
		log.Panic("nameMapper is nil")
	}

	// dstVal must be a nil pointer or a valid value
	err := assignValue(dstVal, reflect.ValueOf(src), nameMapper)
	if err != nil {
		log.Error(err)
	}
	return err
}

// dstVal is valid value or pointer to value
func assignValue(dstVal reflect.Value, srcVal reflect.Value, nameMapper NameMapper) error {
	if !srcVal.IsValid() {
		return errors.New("srcVal is invalid")
	}

	if !dstVal.IsValid() {
		log.Panicf("invalid values:dstVal=%#v,srcVal=%#v", dstVal, srcVal)
	}

	v := dstVal
	if v.Kind() == reflect.Ptr {
		if v.IsNil() && v.CanSet() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	for (srcVal.Kind() == reflect.Ptr || srcVal.Kind() == reflect.Interface) && !srcVal.IsNil() {
		srcVal = srcVal.Elem()
	}

	if !v.CanSet() {
		log.Panic("can't set")
	}

	switch v.Kind() {
	case reflect.Bool:
		b, err := types.ParseBool(srcVal.Interface())
		if err != nil {
			log.Error(err)
			return err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := types.ParseInt(srcVal.Interface())
		if err != nil {
			log.Error(err)
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := types.ParseInt(srcVal.Interface())
		if err != nil {
			log.Error(err)
			return err
		}
		v.SetUint(uint64(i))
	case reflect.Float32, reflect.Float64:
		i, err := types.ParseFloat(srcVal.Interface())
		if err != nil {
			log.Error(err)
			return err
		}
		v.SetFloat(i)
	case reflect.String:
		if srcVal.Kind() != reflect.String {
			err := errors.New("source value is not string")
			log.Error(err)
			return err
		}
		v.SetString(srcVal.String())
	case reflect.Slice:
		if srcVal.Kind() != reflect.Slice {
			err := errors.New("source value is slice")
			log.Error(err)
			return err
		}
		v.Set(reflect.MakeSlice(v.Type(), srcVal.Len(), srcVal.Cap()))
		for i := 0; i < srcVal.Len(); i++ {
			err := assignValue(v.Index(i), srcVal.Index(i), nameMapper)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	case reflect.Map:
		err := mapToMap(v, srcVal, nameMapper)
		if err != nil {
			return err
		}
	case reflect.Struct:
		var err error
		if srcVal.Kind() == reflect.Map {
			err = mapToStruct(v, srcVal, nameMapper)
		} else if srcVal.Kind() == reflect.Struct {
			err = structToStruct(v, srcVal, nameMapper)
		} else {
			err = errors.New("srcVal is not struct or map")
		}

		if err != nil {
			log.Errorf("err:%s src:%s dst:%s", err, srcVal.Kind(), v.Kind())
			return err
		}
	default:
		log.Panicf("unknown kind=%s", v.Kind().String())
	}

	if dstVal.Kind() == reflect.Ptr && dstVal.IsNil() {
		dstVal.Set(v.Addr())
	}
	return nil
}

// dstVal is map
// srcVal is map
func mapToMap(dstVal reflect.Value, srcVal reflect.Value, nameMapper NameMapper) error {
	if dstVal.Kind() != reflect.Map {
		panic("not map")
	}

	if srcVal.Kind() != reflect.Map {
		err := errors.New("srcVal is not map")
		log.Error(err)
		return err
	}

	if !srcVal.Type().Key().AssignableTo(dstVal.Type().Key()) {
		msg := fmt.Sprintf("%s can't be assigned to %s", srcVal.Type().Key().String(), dstVal.Type().Key().String())
		err := errors.New(msg)
		log.Error(err)
		return err
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
			err := assignValue(kv, srcVal.MapIndex(k), nameMapper)
			if err != nil {
				return err
			}
			dstVal.SetMapIndex(k, kv)
		default:
			kv := reflect.New(de)
			err := assignValue(kv, srcVal.MapIndex(k), nameMapper)
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
func mapToStruct(dstVal reflect.Value, srcVal reflect.Value, nameMapper NameMapper) error {
	if dstVal.Kind() != reflect.Struct {
		panic("not struct")
	}

	if srcVal.Kind() == reflect.Interface {
		srcVal = srcVal.Elem()
	}

	if srcVal.Kind() != reflect.Map {
		err := errors.New("srcVal is not map")
		log.Error(err)
		return err
	}

	if srcVal.Type().Key().Kind() != reflect.String {
		err := errors.New("key type must be string")
		log.Error(err)
		return err
	}

	for i := 0; i < dstVal.NumField(); i++ {
		fv := dstVal.Field(i)
		if fv.IsValid() == false || fv.CanSet() == false {
			continue
		}

		ft := dstVal.Type().Field(i)
		if ft.Anonymous {
			err := assignValue(fv, srcVal, nameMapper)
			if err != nil {
				return err
			}
			continue
		}

		for _, key := range srcVal.MapKeys() {
			if !nameMapper.MapName(key.String(), ft.Name) {
				continue
			}

			fsv := srcVal.MapIndex(key)
			if !fsv.IsValid() {
				log.Warnf("invalid value for key:%s", ft.Name)
				continue
			}

			err := assignValue(fv, reflect.ValueOf(fsv.Interface()), nameMapper)
			if err != nil {
				log.Error(err, ft.Name)
				return err
			}
			break
		}
	}
	return nil
}

// dstVal is struct
// srcVal is struct
func structToStruct(dstVal reflect.Value, srcVal reflect.Value, nameMapper NameMapper) error {
	if dstVal.Kind() != reflect.Struct {
		panic("not struct")
	}

	if srcVal.Kind() == reflect.Interface {
		srcVal = srcVal.Elem()
	}

	if srcVal.Kind() != reflect.Struct {
		err := errors.New("srcVal is not struct")
		log.Error(err)
		return err
	}

	for i := 0; i < dstVal.NumField(); i++ {
		fv := dstVal.Field(i)
		if fv.IsValid() == false || fv.CanSet() == false {
			continue
		}

		ft := dstVal.Type().Field(i)
		if ft.Anonymous {
			err := assignValue(fv, srcVal, nameMapper)
			if err != nil {
				return err
			}
			continue
		}

		for i := 0; i < srcVal.NumField(); i++ {
			sfv := srcVal.Field(i)
			sfName := srcVal.Type().Field(i).Name
			if sfv.IsValid() == false || sfName[0] < 'A' || sfName[0] > 'Z' {
				continue
			}

			if !nameMapper.MapName(sfName, ft.Name) {
				continue
			}

			err := assignValue(fv, reflect.ValueOf(sfv.Interface()), nameMapper)
			if err != nil {
				log.Error(err, ft.Name)
				return err
			}
			break
		}
	}
	return nil
}

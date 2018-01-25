package mapper

import (
	"reflect"
	"strings"
)

// KindNameMapper represents an incasesensitive mapper which is very kind
func KindNameMapper(srcName, dstName string) bool {
	return strings.ToLower(srcName) == strings.ToLower(dstName)
}

//Copy copies src's fields to ptrDst's fields by name
func Copy(ptrDst interface{}, src interface{}, nameMapper func(srcFieldName, dstFieldName string) bool) {
	dv := reflect.ValueOf(ptrDst)
	sv := reflect.ValueOf(src)
	if dv.IsValid() == false || sv.IsValid() == false {
		return
	}

	if dv.Kind() != reflect.Ptr {
		return
	}
	dv = dv.Elem()
	if dv.Kind() != reflect.Struct {
		return
	}

	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}

	if sv.Kind() != reflect.Struct {
		return
	}

	if nameMapper == nil {
		for i := 0; i < sv.NumField(); i++ {
			sfv := sv.Field(i)
			if sfv.IsValid() == false {
				continue
			}

			dfv := dv.FieldByName(sv.Type().Field(i).Name)
			if dfv.IsValid() && dfv.CanSet() {
				if sfv.Type().AssignableTo(dfv.Type()) {
					dfv.Set(sfv)
				} else if sfv.Type().ConvertibleTo(dfv.Type()) {
					dfv.Set(sfv.Convert(dfv.Type()))
				}
			}
		}

		return
	}

	for i := 0; i < sv.NumField(); i++ {
		sfv := sv.Field(i)
		sfName := sv.Type().Field(i).Name
		if sfv.IsValid() == false || sfName[0] < 'A' || sfName[0] > 'Z' {
			continue
		}

		for j := 0; j < dv.NumField(); j++ {
			dfv := dv.Field(j)
			if nameMapper(sfName, dv.Type().Field(j).Name) {
				if dfv.IsValid() && dfv.CanSet() {
					if sfv.Type().AssignableTo(dfv.Type()) {
						dfv.Set(sfv)
					} else if sfv.Type().ConvertibleTo(dfv.Type()) {
						dfv.Set(sfv.Convert(dfv.Type()))
					}
				}
				break
			}
		}
	}

}

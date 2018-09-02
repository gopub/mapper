package mapper

import "strings"

type NameMapper interface {
	MapName(srcName, dstName string) bool
}

type NameMapFunc func(string, string) bool

func (f NameMapFunc) MapName(srcName, dstName string) bool {
	return f(srcName, dstName)
}

// DefaultNameMapper represents an incasesensitive mapper which is very kind
type DefaultNameMapper struct {
}

func (m *DefaultNameMapper) MapName(srcName, dstName string) bool {
	return strings.ToLower(srcName) == strings.ToLower(dstName)
}

var _defaultNameMapper = &DefaultNameMapper{}

func NameMapperWithMap(srcToDst map[string]string) NameMapper {
	return NameMapFunc(func(srcName, dstName string) bool {
		if v, ok := srcToDst[srcName]; ok {
			return v == dstName
		}

		return srcName == dstName
	})
}

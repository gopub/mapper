package mapper

import (
	"strings"

	"github.com/gopub/conv"
)

type NameMapper interface {
	MapName(srcName, dstName string) bool
}

type NameMapFunc func(string, string) bool

func (f NameMapFunc) MapName(srcName, dstName string) bool {
	return f(srcName, dstName)
}

// CaseInsensitiveNameMapper represents an incasesensitive mapper which is very kind
var CaseInsensitiveNameMapper NameMapFunc = func(srcName, dstName string) bool {
	return strings.ToLower(srcName) == strings.ToLower(dstName)
}

var defaultNameMapper = CaseInsensitiveNameMapper

func NameMapperWithMap(srcToDst map[string]string) NameMapper {
	return NameMapFunc(func(srcName, dstName string) bool {
		if v, ok := srcToDst[srcName]; ok {
			return v == dstName
		}

		return srcName == dstName
	})
}

var SnakeToCamelNameMapper NameMapFunc = func(snakeSrcName string, camelDstName string) bool {
	srcName := strings.ToLower(conv.ToCamel(snakeSrcName))
	dstName := strings.ToLower(camelDstName)
	return srcName == dstName
}

var CamelToSnakeNameMapper NameMapFunc = func(camelSrcName string, snakeDstName string) bool {
	srcName := strings.ToLower(camelSrcName)
	dstName := strings.ToLower(conv.ToCamel(snakeDstName))
	return srcName == dstName
}

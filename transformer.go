package param

type Transformer interface {
	Transform(v interface{}) (result interface{}, err error)
}

type TransformFunc func(i interface{}) (result interface{}, err error)

func (p TransformFunc) Transform(i interface{}) (result interface{}, err error) {
	return p(i)
}

var _nameToTransformer = map[string]Transformer{}

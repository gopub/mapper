package goparam

type Foo struct {
	Link string `param:"name,min=1,max=300,type=url,transformer=toURL"`
}

package main

import (
	"fmt"
	"github.com/natande/goparam"
)

type User struct {
	Name string `param:"min=6"`
}

type Foo struct {
	User  *User
	Nick  string  `param:"min=2,max=10"`
	Email string  `param:"pattern=email,optional"`
	Price float32 `param:"p,min=2.2,max=3.5"`
}

func main() {
	//f := &Foo{
	//	Nick:  "aa",
	//	Email: "bb@c.com",
	//}
	//
	//err := goparam.Validate(f)
	//fmt.Println(err)

	params := map[string]interface{}{
		"name":  "lisi",
		"nick":  "zhangsan",
		"email": "k",
		"user":  map[string]interface{}{"name": "lisi"},
		"p":     2.0,
	}
	var f *Foo
	err := goparam.Assign(params, &f)
	fmt.Println(err)
}

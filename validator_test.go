package goparam

import "testing"

type Foo struct {
	Nick  string `param:"min=2,max=10"`
	Email string `param:"pattern=email"`
}

func TestValidate(t *testing.T) {
	f := &Foo{
		Nick:  "a",
		Email: "aa",
	}

	err := Validate(f)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

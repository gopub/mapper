package mapper_test

import (
	"github.com/gopub/mapper"
	"testing"
)

type Person struct {
	Age       int
	Birthdate string
}

func (p *Person) Error() string {
	return ""
}

type Teacher struct {
	Name      string
	Age       int
	BirthDate string
}

func TestCopy(t *testing.T) {
	p := &Person{}
	te := &Teacher{}
	te.Age = 20
	te.Name = "h"
	te.BirthDate = "2009-01-02"
	mapper.Copy(p, te, nil)
	if p.Age != te.Age {
		t.FailNow()
	}

	mapper.Copy(p, te, mapper.KindNameMapper)
	if p.Birthdate != te.BirthDate {
		t.FailNow()
	}
}

package gofmcon

import (
	"testing"
)


func TestFieldsCount(t *testing.T) {

	a := []FMQueryFieldGroup{}
	for i := 0; i < 5; i++ {
		f := FMQueryField{}
		g := FMQueryFieldGroup{Fields: []FMQueryField{f, f, f}}
		a = append(a, g)
	}

	q := FMQuery{ QueryFields: a}
	if q.fieldsCount() != 15 {
		t.Error(`FMquery fieldsCount is not correct`)
	}
}

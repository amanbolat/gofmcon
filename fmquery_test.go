package gofmcon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFieldsCount(t *testing.T) {

	var a []FMQueryFieldGroup
	for i := 0; i < 5; i++ {
		f := FMQueryField{}
		g := FMQueryFieldGroup{Fields: []FMQueryField{f, f, f}}
		a = append(a, g)
	}

	q := FMQuery{QueryFields: a}
	assert.Equal(t, 15, q.fieldsCount(), "FMQuery fieldsCount is not correct")
}

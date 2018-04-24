package gofmcon

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestFieldsCount(t *testing.T) {

	a := []FMQueryFieldGroup{}
	for i := 0; i < 5; i++ {
		f := FMQueryField{}
		g := FMQueryFieldGroup{Fields: []FMQueryField{f, f, f}}
		a = append(a, g)
	}

	q := FMQuery{ QueryFields: a}
	assert.Equal(t, 15, q.fieldsCount(), "FMquery fieldsCount is not correct")
}

func TestFMQuery_QueryString(t *testing.T) {
	q := NewFMQuery("db", "layout", Edit)

	q.WithScriptParams("|", "hello", "world")
	q.WithPostFindScripts("audit_log")
	t.Log(q.QueryString())
}
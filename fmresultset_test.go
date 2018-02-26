package gofmcon_test

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"encoding/xml"
	"bitbucket.org/amanbolat/gofmcon"
	"fmt"
)

func TestFMResultsetXMLUnmarshal(t *testing.T) {
	testFile, err := os.Open("simple_test_1.xml")
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(testFile)
	assert.NoError(t, err)

	fmResultSet := &gofmcon.FMResultset{}
	err = xml.Unmarshal(b, fmResultSet)
	assert.NoError(t, err)

	fmt.Printf("%+v", fmResultSet)
}

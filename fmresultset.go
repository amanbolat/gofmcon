package gofmcon

import "strings"

type FMResultset struct {
	Resultset Resultset `xml:"resultset"`
	Version   string    `xml:"version,attr"`
	FMError   FMError
}

type FMError struct {
	Code int `xml:"code,attr"`
}

type Resultset struct {
	Count   int      `xml:"count,attr" json:"count"`
	Fetched int      `xml:"fetch-size,attr" json:"fetched"`
	Records []Record `xml:"record"`
}

type Record struct {
	ID         int          `xml:"record-id,attr"`
	Fields     []Field      `xml:"field"`
	RelatedSet []RelatedSet `xml:"relatedset"`
}

type RelatedSet struct {
	Count   int      `xml:"count,attr"`
	Table   string   `xml:"table,attr"`
	Records []Record `xml:"record"`
}

type Field struct {
	FieldName string `xml:"name,attr" json:"fieldName"`
	FieldData string `xml:"data" json:"fieldData"`
}

func (r *Record) RelatedSetFromTable(t string) RelatedSet {
	rset := RelatedSet{}
	for _, elem := range r.RelatedSet {
		if elem.Table == t {
			rset = elem
		}
	}
	return rset
}

func (r *Record) DataFromFieldInex(fname string) string {
	var s string
	for _, elem := range r.Fields {
		if strings.Contains(elem.FieldName, fname) {
			s = elem.FieldData
		}
	}
	return s
}

func (r *Record) ToMap() map[string]string {
	fMap := make(map[string]string)
	for _, f := range r.Fields {
		fMap[f.FieldName] = f.FieldData
	}
	return fMap
}

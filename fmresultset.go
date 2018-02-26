package gofmcon

import (
	"strings"
	"fmt"
)

type FMResultset struct {
	Resultset Resultset `xml:"resultset"`
	DataSource DataSource `xml:"datasource"`
	Version   string    `xml:"version,attr"`
	FMError   FMError	`xml:"error"`
}

func (rs *FMResultset) HasError() bool {
	return rs.FMError.Code != 0
}

type DataSource struct {
	Database string `xml:"database"`
	DateFormat string `xml:"date_format"`
	Layout string `xml:"layout"`
	Table string `xml:"table"`
	TimeFormat string `xml:"time-format"`
	TimestampFormat string `xml:"timestamp-format"`
	TotalCount int `xml:"total-count"`
}

type FMError struct {
	Code int `xml:"code,attr"`
}

func (e *FMError) String() string {
	return fmt.Sprintf("filemaker_error: %s", FileMakerErrorCodes[e.Code])
}

func (e *FMError) Error() string {
	return fmt.Sprintf("filemaker_error: %s", FileMakerErrorCodes[e.Code])
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
	rSet := RelatedSet{}
	for _, elem := range r.RelatedSet {
		if elem.Table == t {
			rSet = elem
		}
	}
	return rSet
}

// Field returns field data if exists
func (r *Record) Field(name string) string {
	var s string
	for _, elem := range r.Fields {
		if strings.Contains(elem.FieldName, name) {
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

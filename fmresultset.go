package gofmcon

import (
	"fmt"
)

type FMResultset struct {
	Resultset  *Resultset  `xml:"resultset"`
	DataSource *DataSource `xml:"datasource"`
	Version    string     `xml:"version,attr"`
	FMError    FMError    `xml:"error"`
}

func (rs *FMResultset) HasError() bool {
	return rs.FMError.Code != 0
}

type DataSource struct {
	Database        string `xml:"database"`
	DateFormat      string `xml:"date_format"`
	Layout          string `xml:"layout"`
	Table           string `xml:"table"`
	TimeFormat      string `xml:"time-format"`
	TimestampFormat string `xml:"timestamp-format"`
	TotalCount      int    `xml:"total-count"`
}

type FMError struct {
	Code int `xml:"code,attr"`
}

func (e *FMError) String() string {
	return fmt.Sprintf("filemaker_error: %s", fileMakerErrorCodes[e.Code])
}

func (e *FMError) Error() string {
	return fmt.Sprintf("filemaker_error: %s", fileMakerErrorCodes[e.Code])
}

type Resultset struct {
	Count   int      `xml:"count,attr" json:"count"`
	Fetched int      `xml:"fetch-size,attr" json:"fetched"`
	Records []*Record `xml:"record"`
}

type Record struct {
	ID         int          `xml:"record-id,attr"`
	Fields     []*Field      `xml:"field"`
	fieldsMap  map[string][]string
	RelatedSet []*RelatedSet `xml:"relatedset"`
}

type RelatedSet struct {
	Count   int      `xml:"count,attr"`
	Table   string   `xml:"table,attr"`
	Records []*Record `xml:"record"`
}

func (rs *Resultset) prepareRecords() {
	for _, r := range rs.Records {
		r.makeFieldsMap()
	}
}

func (r *Record) makeFieldsMap() {
	if r.fieldsMap == nil {
		r.fieldsMap = map[string][]string{}
	}
	for _, f := range r.Fields {
		r.fieldsMap[f.FieldName] = f.FieldData
	}

	for _, rs := range r.RelatedSet {
		for _, rr := range rs.Records {
			rr.makeFieldsMap()
		}
	}
}

type Field struct {
	FieldName string `xml:"name,attr" json:"fieldName"`
	FieldData []string `xml:"data" json:"fieldData"`
}

func (r *Record) RelatedSetFromTable(t string) *RelatedSet {
	rSet := &RelatedSet{}
	for _, elem := range r.RelatedSet {
		if elem.Table == t {
			rSet = elem
		}
	}
	return rSet
}

// Field returns fields first data if exists
func (r *Record) Field(name string) string {
	f := r.fieldsMap[name]
	if len(f) < 1{
		return ""
	}
	return f[0]
}

// ContainerField returns all data of container field
func (r *Record) ContainerField(name string) []string {
	return r.fieldsMap[name]
}

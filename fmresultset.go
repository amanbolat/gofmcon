package gofmcon

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DATE_FORMAT = "01/02/2006"
	TIME_FORMAT = "15:04:05"
	TIMESTAMP_FORMAT = "01/02/2006 15:04:05"
)

type FMResultset struct {
	Resultset  *Resultset  `xml:"resultset"`
	DataSource *DataSource `xml:"datasource"`
	MetaData   *MetaData   `xml:"metadata"`
	Version    string      `xml:"version,attr"`
	FMError    FMError     `xml:"error"`
}

func (rs *FMResultset) prepareRecords() {
	for _, r := range rs.Resultset.Records {
		r.makeFieldsMap(false, rs.MetaData.getAllFieldDefinitions())
	}
}

func (rs *FMResultset) HasError() bool {
	return rs.FMError.Code != 0
}

type DataSource struct {
	Database        string `xml:"database,attr"`
	DateFormat      string `xml:"date_format,attr"`
	Layout          string `xml:"layout,attr"`
	Table           string `xml:"table,attr"`
	TimeFormat      string `xml:"time-format,attr"`
	TimestampFormat string `xml:"timestamp-format,attr"`
	TotalCount      int    `xml:"total-count,attr"`
}

type MetaData struct {
	FieldDefinitions     []*FieldDefinition    `xml:"field-definition"`
	RelatedSetDefinition *RelatedSetDefinition `xml:"relatedset-definition"`
}

func (md MetaData) getAllFieldDefinitions() []FieldDefinition {
	var definitions []FieldDefinition
	for _, def := range md.FieldDefinitions {
		definitions = append(definitions, *def)
	}

	for _, def := range md.RelatedSetDefinition.FieldDefinitions {
		definitions = append(definitions, *def)
	}

	return definitions
}

type RelatedSetDefinition struct {
	Table            string             `xml:"table,attr"`
	FieldDefinitions []*FieldDefinition `xml:"field-definition"`
}

type fieldDefinition struct {
	Name          string `xml:"name,attr"`
	AutoEnter     string `xml:"auto-enter,attr"`
	FourDigitYear string `xml:"four-digit-year,attr"`
	Global        string `xml:"global,attr"`
	MaxRepeat     string `xml:"max-repeat,attr"`
	NotEmpty      string `xml:"not-empty,attr"`
	NumericOnly   string `xml:"numeric-only,attr"`
	Result        string `xml:"result,attr"`
	TimeOfDay     string `xml:"time-of-day,attr"`
	Type          string `xml:"type,attr"`
}

type FieldType string

const (
	TypeText      FieldType = "text"
	TypeNumber    FieldType = "number"
	TypeDate      FieldType = "date"
	TypeTime      FieldType = "time"
	TypeTimestamp FieldType = "timestamp"
	TypeContainer FieldType = "container"
)

type FieldDefinition struct {
	Name          string `xml:"name,attr"`
	AutoEnter     bool
	FourDigitYear bool
	Global        bool
	MaxRepeat     int
	NotEmpty      bool
	NumericOnly   bool
	TimeOfDay     bool
	Type          FieldType
}

type FieldsDefinitions []FieldDefinition

func (fds FieldsDefinitions) getType(name string) FieldType {
	for _, fd := range fds {
		if fd.Name == name {
			return fd.Type
		}
	}

	return ""
}

func (fds FieldsDefinitions) getMaxRepeat(name string) int {
	for _, fd := range fds {
		if fd.Name == name {
			return fd.MaxRepeat
		}
	}

	return 1
}

func (f *FieldDefinition) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var fd fieldDefinition
	err := d.DecodeElement(&fd, &start)
	if err != nil {
		return err
	}

	f.AutoEnter = getBoolFromString(fd.AutoEnter)
	f.FourDigitYear = getBoolFromString(fd.FourDigitYear)
	f.Global = getBoolFromString(fd.Global)
	f.NotEmpty = getBoolFromString(fd.NotEmpty)
	f.NumericOnly = getBoolFromString(fd.NumericOnly)
	f.TimeOfDay = getBoolFromString(fd.TimeOfDay)
	f.Name = fd.Name
	f.Type = FieldType(fd.Result)
	mr, _ := strconv.Atoi(fd.MaxRepeat)
	f.MaxRepeat = mr

	return nil
}

func getBoolFromString(str string) bool {
	if str == "yes" {
		return true
	}
	return false
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
	Count   int       `xml:"count,attr" json:"count"`
	Fetched int       `xml:"fetch-size,attr" json:"fetched"`
	Records []*Record `xml:"record"`
}

type Record struct {
	ID         int      `xml:"record-id,attr"`
	Fields     []*Field `xml:"field"`
	fieldsMap  map[string]interface{}
	RelatedSet []*RelatedSet `xml:"relatedset"`
}

type RelatedSet struct {
	Count   int       `xml:"count,attr"`
	Table   string    `xml:"table,attr"`
	Records []*Record `xml:"record"`
}

func (r *Record) makeFieldsMap(isNested bool, fieldsDefinitions FieldsDefinitions) {
	if r.fieldsMap == nil {
		r.fieldsMap = map[string]interface{}{}
	}

	for _, f := range r.Fields {
		var dataArr []interface{}

		for _, val := range f.Data {
			switch fieldsDefinitions.getType(f.Name) {
			case TypeNumber:
				number, err := strconv.ParseFloat(val, 64)
				if err != nil {
					dataArr = append(dataArr, val)
				} else {
					dataArr = append(dataArr, number)
				}
			case TypeDate:
				t, _ := time.Parse(DATE_FORMAT, val)
				dataArr = append(dataArr, t)
			case TypeTime:
				t, _ := time.Parse(TIME_FORMAT, val)
				dataArr = append(dataArr, t)
			case TypeTimestamp:
				t, _ := time.Parse(TIMESTAMP_FORMAT, val)
				dataArr = append(dataArr, t)
			default:
				dataArr = append(dataArr, val)
			}
		}

		var fieldData interface{}

		maxRepeat := fieldsDefinitions.getMaxRepeat(f.Name)

		if maxRepeat > 1 || len(dataArr) > 1 {
			fieldData = dataArr
		} else if len(dataArr) == 1{
			fieldData = dataArr[0]
		}

		fieldName := f.Name
		if isNested {
			idx := strings.Index(f.Name, "::")
			if idx > 0 {
				fieldName = f.Name[idx+2:]
			}
		}

		r.fieldsMap[fieldName] = fieldData
	}

	for _, rs := range r.RelatedSet {
		var relatedRecordsFieldMaps []interface{}
		for _, rr := range rs.Records {
			rr.makeFieldsMap(true, fieldsDefinitions)
			relatedRecordsFieldMaps = append(relatedRecordsFieldMaps, rr.fieldsMap)
		}

		r.fieldsMap[rs.Table] = relatedRecordsFieldMaps
	}
}

type Field struct {
	Name string   `xml:"name,attr" json:"field"`
	Data []string `xml:"data" json:"data"`
	Type FieldType
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

// Field returns fields data for given field name
func (r *Record) Field(name string) interface{} {
	return r.fieldsMap[name]
}

func (r *Record) JsonFields() ([]byte, error) {
	return  json.MarshalIndent(r.fieldsMap, "", "	")
}


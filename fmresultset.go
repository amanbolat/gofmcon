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
	// DateFormat is a format of date on a particular layout
	DateFormat = "01/02/2006"
	// TimeFormat is a format of time on a particular layout
	TimeFormat = "15:04:05"
	// TimestampFormat is a format of timestamp on a particular layout
	TimestampFormat = "01/02/2006 15:04:05"
)

// FMResultset is a collection of ResultSets
type FMResultset struct {
	Resultset  *Resultset  `xml:"resultset"`
	DataSource *DataSource `xml:"datasource"`
	MetaData   *MetaData   `xml:"metadata"`
	Version    string      `xml:"version,attr"`
	FMError    FMError     `xml:"error"`
}

func (rs *FMResultset) prepareRecords() {
	var fd FieldsDefinitions
	if rs.MetaData != nil {
		fd = rs.MetaData.getAllFieldDefinitions()
	}
	for _, r := range rs.Resultset.Records {
		r.makeFieldsMap(false, fd)
	}
}

// HasError checks if FMResultset was fetched with an error
func (rs *FMResultset) HasError() bool {
	return rs.FMError.Code != 0
}

// DataSource store database name, layout name and time formats
type DataSource struct {
	Database        string `xml:"database,attr"`
	DateFormat      string `xml:"date_format,attr"`
	Layout          string `xml:"layout,attr"`
	Table           string `xml:"table,attr"`
	TimeFormat      string `xml:"time-format,attr"`
	TimestampFormat string `xml:"timestamp-format,attr"`
	TotalCount      int    `xml:"total-count,attr"`
}

// MetaData store fields' and related sets' meta information
type MetaData struct {
	FieldDefinitions     []*FieldDefinition    `xml:"field-definition"`
	RelatedSetDefinition *RelatedSetDefinition `xml:"relatedset-definition"`
}

func (md MetaData) getAllFieldDefinitions() []FieldDefinition {
	var definitions []FieldDefinition
	for _, def := range md.FieldDefinitions {
		definitions = append(definitions, *def)
	}

	if md.RelatedSetDefinition == nil {
		return definitions
	}

	for _, def := range md.RelatedSetDefinition.FieldDefinitions {
		definitions = append(definitions, *def)
	}

	return definitions
}

// RelatedSetDefinition is a meta information of related set
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

// FieldType represents a type of the field
type FieldType string

const (
	// TypeText is a text type of the field
	TypeText FieldType = "text"
	// TypeNumber is a number type of the field
	TypeNumber FieldType = "number"
	// TypeDate is a date type of the field
	TypeDate FieldType = "date"
	// TypeTime is a time type of the field
	TypeTime FieldType = "time"
	// TypeTimestamp is a timestamp type of the field
	TypeTimestamp FieldType = "timestamp"
	// TypeContainer is a container type of the field
	TypeContainer FieldType = "container"
)

// FieldDefinition store information about a field in given layout
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

// FieldsDefinitions is type of []FieldDefinition
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

// UnmarshalXML serializes xml data of FieldDefinition into the object
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

// FMError represents a FileMaker error
type FMError struct {
	Code int `xml:"code,attr"`
}

func (e *FMError) String() string {
	return fmt.Sprintf("filemaker_error: %s", FileMakerErrorCodes[e.Code])
}

func (e *FMError) Error() string {
	return fmt.Sprintf("filemaker_error: %s", FileMakerErrorCodes[e.Code])
}

// Resultset is a set of records with meta information
type Resultset struct {
	Count   int       `xml:"count,attr" json:"count"`
	Fetched int       `xml:"fetch-size,attr" json:"fetched"`
	Records []*Record `xml:"record"`
}

// Record is FileMaker record
type Record struct {
	ID         int      `xml:"record-id,attr"`
	Fields     []*Field `xml:"field"`
	fieldsMap  map[string]interface{}
	RelatedSet []*RelatedSet `xml:"relatedset"`
}

// RelatedSet is a set of records returned from FileMaker database
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
					dataArr = append(dataArr, nil)
				} else {
					dataArr = append(dataArr, number)
				}
			case TypeDate:
				t, _ := time.Parse(DateFormat, val)
				dataArr = append(dataArr, t)
			case TypeTime:
				t, _ := time.Parse(TimeFormat, val)
				dataArr = append(dataArr, t)
			case TypeTimestamp:
				t, _ := time.Parse(TimestampFormat, val)
				dataArr = append(dataArr, t)
			default:
				dataArr = append(dataArr, val)
			}
		}

		var fieldData interface{}

		maxRepeat := fieldsDefinitions.getMaxRepeat(f.Name)

		if maxRepeat > 1 || len(dataArr) > 1 {
			fieldData = dataArr
		} else if len(dataArr) == 1 {
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

// Field stands for field in a record
type Field struct {
	Name string   `xml:"name,attr" json:"field"`
	Data []string `xml:"data" json:"data"`
	Type FieldType
}

// RelatedSetFromTable returns the set of related records from given
// related table
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

// JSONFields return JSON representation of Record
func (r *Record) JSONFields() ([]byte, error) {
	return json.MarshalIndent(r.fieldsMap, "", "	")
}

package gofmcon

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	fmNoRecordID = -1
	fmAllRecords = -1
)

// FMAction is a type iof action can be done to the record
type FMAction string

const (
	// Find -findquery
	Find FMAction = "-findquery"
	// FindAll -findall
	FindAll FMAction = "-findall"
	// FindAny findany
	FindAny FMAction = "-findany"
	// New -new
	New FMAction = "-new"
	// Edit -edit
	Edit FMAction = "-edit"
	// Delete -delete
	Delete FMAction = "-delete"
	// Duplicate -dup
	Duplicate FMAction = "-dup"
)

func (a FMAction) String() string {
	return string(a)
}

// FMSortOrder is a type of order
type FMSortOrder string

const (
	// Ascending -ascend
	Ascending FMSortOrder = "ascend"
	// Descending -descend
	Descending FMSortOrder = "descend"
	// Custom -custom
	Custom FMSortOrder = "custom"
)

func (so FMSortOrder) String() string {
	return string(so)
}

// FMSortField is a field that should be sorted during the query
type FMSortField struct {
	Name  string
	Order FMSortOrder
}

// FMFieldOp is type of operator for a FMField
type FMFieldOp string

const (
	// Equal -eq
	Equal FMFieldOp = "eq"
	// Contains -cn
	Contains FMFieldOp = "cn"
	// BeginsWith -bw
	BeginsWith FMFieldOp = "bw"
	// EndsWith -ew
	EndsWith FMFieldOp = "ew"
	// GreaterThan -gt
	GreaterThan FMFieldOp = "gt"
	// GreaterThanEqual -gte
	GreaterThanEqual FMFieldOp = "gte"
	// LessThan -lt
	LessThan FMFieldOp = "lt"
	// LessThanEqual -lte
	LessThanEqual FMFieldOp = "lte"
)

// FMQueryField is a field used in FMQuery
type FMQueryField struct {
	Name  string
	Value string
	Op    FMFieldOp
}

func (qf *FMQueryField) valueWithOp() string {
	switch qf.Op {
	case Equal:
		return "==" + qf.Value
	case Contains:
		return "==*" + qf.Value + "*"
	case BeginsWith:
		return "==" + qf.Value + "*"
	case EndsWith:
		return "==*" + qf.Value
	case GreaterThan:
		return ">" + qf.Value
	case GreaterThanEqual:
		return ">=" + qf.Value
	case LessThan:
		return "<" + qf.Value
	case LessThanEqual:
		return "<=" + qf.Value
	default:
		return qf.Value
	}
}

// FMLogicalOp is a type for logical operators
type FMLogicalOp string

const (
	// And operator
	And FMLogicalOp = "and"
	// Or operator
	Or FMLogicalOp = "or"
	// Not operator
	Not FMLogicalOp = "not"
)

// FMQueryFieldGroup groups of fields used for the find request
type FMQueryFieldGroup struct {
	Op     FMLogicalOp
	Fields []FMQueryField
}

func (fg *FMQueryFieldGroup) simpleFieldsString() string {
	var strArray []string
	for _, f := range fg.Fields {
		strArray = append(strArray, url.QueryEscape(f.Name)+"="+url.QueryEscape(f.Value))
	}
	return strings.Join(strArray, "&")
}

// FMQuery represent the query you are sending to the server
type FMQuery struct {
	Database            string
	Layout              string
	Action              FMAction
	QueryFields         []FMQueryFieldGroup
	SortFields          []FMSortField
	RecordID            int // default should be -1
	PreSortScript       string
	PreFindScript       string
	PostFindScript      string
	PreSortScriptParam  string
	PreFindScriptParam  string
	PostFindScriptParam string
	ResponseLayout      string
	ResponseFields      []string
	MaxRecords          int // default should be -1
	SkipRecords         int // default should be 0
	Query               map[string]string
}

// NewFMQuery creates new FMQuery object
func NewFMQuery(database string, layout string, action FMAction) *FMQuery {
	return &FMQuery{
		Database:    database,
		Layout:      layout,
		Action:      action,
		RecordID:    fmNoRecordID,
		MaxRecords:  fmAllRecords,
		SkipRecords: 0,
	}
}

// WithRecordID sets RecordID for the query.
// MUST have if you want to edit/delete record
func (q *FMQuery) WithRecordID(id int) *FMQuery {
	q.RecordID = id
	return q
}

// WithFieldGroups sets groups of fields for find request
func (q *FMQuery) WithFieldGroups(fieldGroups ...FMQueryFieldGroup) *FMQuery {
	q.QueryFields = append(q.QueryFields, fieldGroups...)
	return q
}

// WithFields sets field for the find request
func (q *FMQuery) WithFields(fields ...FMQueryField) *FMQuery {
	group := FMQueryFieldGroup{
		Fields: fields,
		Op:     And,
	}
	q.QueryFields = append(q.QueryFields, group)
	return q
}

// WithSortFields sets sort fields' name and order
func (q *FMQuery) WithSortFields(sortFields ...FMSortField) *FMQuery {
	q.SortFields = append(q.SortFields, sortFields...)
	return q
}

// WithPreSortScript sets PreSortScript and params
func (q *FMQuery) WithPreSortScript(script, param string) *FMQuery {
	q.PreSortScript = script
	q.PreSortScriptParam = param
	return q
}

// WithPreFindScript sets PreFindScript and params
func (q *FMQuery) WithPreFindScript(script, param string) *FMQuery {
	q.PreFindScript = script
	q.PreFindScriptParam = param
	return q
}

// WithPostFindScript sets PostFindScript script and params
func (q *FMQuery) WithPostFindScript(script, param string) *FMQuery {
	q.PostFindScript = script
	q.PostFindScriptParam = param
	return q
}

// WithResponseLayout sets layout name you want to fetch records from
func (q *FMQuery) WithResponseLayout(lay string) *FMQuery {
	q.ResponseLayout = lay
	return q
}

// WithResponseFields adds field names that FileMaker server should return
func (q *FMQuery) WithResponseFields(fields ...string) *FMQuery {
	q.ResponseFields = append(q.ResponseFields, fields...)
	return q
}

// Max sets maximum amount of records to fetch
func (q *FMQuery) Max(n int) *FMQuery {
	q.MaxRecords = n
	return q
}

// Skip skips n amount of recrods
func (q *FMQuery) Skip(n int) *FMQuery {
	q.SkipRecords = n
	return q
}

func withAmp(s string) string {
	if s == "" {
		return ""
	}
	return s + "&"
}

func (q *FMQuery) fieldsCount() int {
	var count int
	for _, group := range q.QueryFields {
		for range group.Fields {
			count++
		}
	}
	return count
}

func (q *FMQuery) dbLayString() string {
	return "-db=" + url.QueryEscape(q.Database) + "&-lay=" + url.QueryEscape(q.Layout) + "&"
}

func (q *FMQuery) sortFieldsString() string {
	var strArray []string
	colNum := 1
	for _, f := range q.SortFields {
		i := strconv.Itoa(colNum)
		str := "-sortfield." + i + "=" + url.QueryEscape(f.Name) + "&-sortorder." + i + "=" + f.Order.String()
		colNum++
		strArray = append(strArray, str)
	}
	return strings.Join(strArray, "&")
}

func (q *FMQuery) responseFieldsString() string {
	var strArray []string
	for _, f := range q.ResponseFields {
		str := "-field=" + url.QueryEscape(f)
		strArray = append(strArray, str)
	}
	return strings.Join(strArray, "&")
}

func (q *FMQuery) scriptsString() string {
	var preSort string
	if q.PreSortScript != "" {
		preSort = "-script.presort=" + url.QueryEscape(q.PreSortScript)
	}
	var preFind string
	if q.PreFindScript != "" {
		preFind = "-script.prefind=" + url.QueryEscape(q.PreFindScript)
	}
	var postFind string
	if q.PostFindScript != "" {
		postFind = "-script=" + url.QueryEscape(q.PostFindScript)
	}

	return withAmp(preSort) + withAmp(preFind) + postFind
}

func (q *FMQuery) scriptParamsString() string {
	var preSort string
	var preFind string
	var postFind string

	if len(q.PreSortScriptParam) > 0 {
		preSort = fmt.Sprintf("-script.presort.param=%s", url.QueryEscape(q.PreSortScriptParam))
	}

	if len(q.PreFindScriptParam) > 0 {
		preFind = fmt.Sprintf("-script.prefind.param=%s", url.QueryEscape(q.PreFindScriptParam))
	}

	if len(q.PostFindScriptParam) > 0 {
		postFind = fmt.Sprintf("-script.param=%s", url.QueryEscape(q.PostFindScriptParam))
	}

	return withAmp(preSort) + withAmp(preFind) + postFind
}

func (q *FMQuery) maxSkipString() string {
	var maxString string
	if q.MaxRecords == fmAllRecords {
		maxString = "-max=all"
	} else {
		maxString = "-max=" + strconv.Itoa(q.MaxRecords)
	}

	return "-skip=" + strconv.Itoa(q.SkipRecords) + "&" + maxString
}

func (q *FMQuery) recordIDString() string {
	if q.RecordID != fmNoRecordID && q.Action != FindAny {
		return "-recid=" + strconv.Itoa(q.RecordID)
	}
	return ""
}

func (q *FMQuery) responseLayoutString() string {
	if q.ResponseLayout == "" {
		return ""
	}
	return "-lay.response=" + url.QueryEscape(q.ResponseLayout)
}

func (q *FMQuery) simpleFieldsString() string {
	var strArray []string
	for _, f := range q.QueryFields {
		strArray = append(strArray, f.simpleFieldsString())
	}
	return strings.Join(strArray, "&")
}

func (q *FMQuery) compoundQueryString() string {
	var segments []string
	var i int
	for _, g := range q.QueryFields {
		switch g.Op {
		case And:
			var fieldsArray []string
			for range g.Fields {
				i++
				str := fmt.Sprintf("q%d", i)
				fieldsArray = append(fieldsArray, str)
			}
			str := "(" + strings.Join(fieldsArray, ",") + ")"
			segments = append(segments, str)
		case Or:
			var fieldsArray []string
			for range g.Fields {
				i++
				str := fmt.Sprintf("(q%d)", i)
				fieldsArray = append(fieldsArray, str)
			}
			str := strings.Join(fieldsArray, ";")
			segments = append(segments, str)
		case Not:
			var fieldsArray []string
			for range g.Fields {
				i++
				str := fmt.Sprintf("q%d", i)
				fieldsArray = append(fieldsArray, str)
			}
			str := "!(" + strings.Join(fieldsArray, ",") + ")"
			segments = append(segments, str)
		}
	}
	return "-query=" + url.QueryEscape(strings.Join(segments, ";"))
}

func (q *FMQuery) compoundFieldsString() string {
	var segments []string
	var i int
	for _, g := range q.QueryFields {
		var strArray []string
		for _, f := range g.Fields {
			i++
			str := "-q" + strconv.Itoa(i) + "=" + url.QueryEscape(f.Name) + "&-q" + strconv.Itoa(i) + ".value=" + url.QueryEscape(f.valueWithOp())
			strArray = append(strArray, str)
		}
		segments = append(segments, strings.Join(strArray, "&"))
	}
	return strings.Join(segments, "&")
}

// QueryString creates query string based on FMQuery
func (q *FMQuery) QueryString() string {
	var startString = q.dbLayString() + withAmp(q.responseLayoutString()) + withAmp(q.scriptsString()) + withAmp(q.scriptParamsString())
	switch q.Action {
	case Delete, Duplicate:
		return startString +
			withAmp(q.recordIDString()) +
			q.Action.String()
	case Edit:
		return startString +
			withAmp(q.recordIDString()) +
			withAmp(q.simpleFieldsString()) +
			q.Action.String()
	case New:
		return startString +
			withAmp(q.simpleFieldsString()) +
			q.Action.String()
	case FindAny:
		return startString +
			q.Action.String()
	case FindAll:
		return startString +
			withAmp(q.sortFieldsString()) +
			withAmp(q.maxSkipString()) +
			q.Action.String()
	case Find:
		if q.RecordID != fmNoRecordID {
			return startString +
				withAmp(q.recordIDString()) + "-find"
		}
		if q.compoundQueryString() == "" || q.compoundFieldsString() == "" {
			return startString +
				withAmp(q.sortFieldsString()) +
				withAmp(q.maxSkipString()) +
				"-find"
		}
		return startString +
			withAmp(q.sortFieldsString()) +
			withAmp(q.maxSkipString()) +
			withAmp(q.compoundQueryString()) +
			withAmp(q.compoundFieldsString()) +
			q.Action.String()
	default:
		return ""
	}
}

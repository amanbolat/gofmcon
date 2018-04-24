package gofmcon

import (
	"net/url"
	"strconv"
	"strings"
	"fmt"
)

const (
	fmNoRecordId = -1
	fmAllRecords = -1
)

type FMAction string

const (
	Find      FMAction = "-findquery"
	FindAll   FMAction = "-findall"
	FindAny   FMAction = "-findany"
	New       FMAction = "-new"
	Edit      FMAction = "-edit"
	Delete    FMAction = "-delete"
	Duplicate FMAction = "-dup"
)

func (a FMAction) String() string {
	return string(a)
}

type FMSortOrder string

const (
	Ascending  FMSortOrder = "ascend"
	Descending FMSortOrder = "descend"
	Custom     FMSortOrder = "custom"
)

func (so FMSortOrder) String() string {
	return string(so)
}

type FMSortField struct {
	Name  string
	Order FMSortOrder
}

type FMFieldOp string

const (
	Equal            FMFieldOp = "eq"
	Contains         FMFieldOp = "cn"
	BeginsWith       FMFieldOp = "bw"
	EndsWith         FMFieldOp = "ew"
	GreaterThan      FMFieldOp = "gt"
	GreaterThanEqual FMFieldOp = "gte"
	LessThan         FMFieldOp = "lt"
	LessThanEqual    FMFieldOp = "lte"
)

type FMQueryField struct {
	Name  string
	Value string
	Op    FMFieldOp
}

func (qf *FMQueryField) ValueWithOp() string {
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

type FMLogicalOp string

const (
	And FMLogicalOp = "and"
	Or  FMLogicalOp = "or"
	Not FMLogicalOp = "not"
)

type FMQueryFieldGroup struct {
	Op     FMLogicalOp
	Fields []FMQueryField
}

func (fg *FMQueryFieldGroup) SimpleFieldsString() string {
	var strArray []string
	for _, f := range fg.Fields {
		strArray = append(strArray, url.QueryEscape(f.Name)+"="+url.QueryEscape(f.Value))
	}
	return strings.Join(strArray, "&")
}

type FMQuery struct {
	Database        string
	Layout          string
	Action          FMAction
	QueryFields     []FMQueryFieldGroup
	SortFields      []FMSortField
	RecordId        int // default should be -1
	PreSortScripts  []string
	PreFindScripts  []string
	PostFindScripts []string
	ScriptParams []string
	ScriptParamsDelimiter string
	ResponseLayout  string
	ResponseFields  []string
	MaxRecords      int // default should be -1
	SkipRecords     int // default should be 0
	Query           map[string]string
}

func NewFMQuery(database string, layout string, action FMAction) *FMQuery {
	return &FMQuery{
		Database:    database,
		Layout:      layout,
		Action:      action,
		RecordId:    fmNoRecordId,
		MaxRecords:  fmAllRecords,
		SkipRecords: 0,
	}
}

func (q *FMQuery) WithRecordId(id int) *FMQuery {
	q.RecordId = id
	return q
}

func (q *FMQuery) WithFieldGroups(fieldGroups ...FMQueryFieldGroup) *FMQuery {
	q.QueryFields = append(q.QueryFields, fieldGroups...)
	return q
}

func (q *FMQuery) WithFields(fields ...FMQueryField) *FMQuery {
	group := FMQueryFieldGroup{
		Fields: fields,
		Op:     And,
	}
	q.QueryFields = append(q.QueryFields, group)
	return q
}

func (q *FMQuery) WithSortFields(sortFields ...FMSortField) *FMQuery {
	q.SortFields = append(q.SortFields, sortFields...)
	return q
}

func (q *FMQuery) WithPreSortScripts(scripts ...string) *FMQuery {
	q.PreSortScripts = append(q.PreSortScripts, scripts...)
	return q
}

func (q *FMQuery) WithPreFindScripts(scripts ...string) *FMQuery {
	q.PreFindScripts = append(q.PreFindScripts, scripts...)
	return q
}

func (q *FMQuery) WithPostFindScripts(scripts ...string) *FMQuery {
	q.PostFindScripts = append(q.PostFindScripts, scripts...)
	return q
}

func (q *FMQuery) WithScriptParams(delimiter string, params ...string) *FMQuery {
	q.ScriptParams = append(q.ScriptParams, params...)
	q.ScriptParamsDelimiter = delimiter
	return q
}

func (q *FMQuery) WithResponseLayout(lay string) *FMQuery {
	q.ResponseLayout = lay
	return q
}

func (q *FMQuery) WithResponseFields(fields ...string) *FMQuery {
	q.ResponseFields = append(q.ResponseFields, fields...)
	return q
}

func (q *FMQuery) WithMaxRecords(max int) *FMQuery {
	q.MaxRecords = max
	return q
}

func (q *FMQuery) WithSkipRecords(skip int) *FMQuery {
	q.SkipRecords = skip
	return q
}

func WithAmp(s string) string {
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
	var preSortsArray []string
	for _, s := range q.PreSortScripts {
		preSortsArray = append(preSortsArray, "-script.presort="+url.QueryEscape(s))
	}
	var preFindsArray []string
	for _, s := range q.PreFindScripts {
		preFindsArray = append(preFindsArray, "-script.prefind="+url.QueryEscape(s))
	}
	var postFinds []string
	for _, s := range q.PostFindScripts {
		postFinds = append(postFinds, "-script="+url.QueryEscape(s))
	}
	return WithAmp(strings.Join(preSortsArray, "&")) +
		WithAmp(strings.Join(preFindsArray, "&")) +
		strings.Join(postFinds, "&")
}

func (q *FMQuery) scriptParamsString() string {
	if len(q.ScriptParams) < 1 {
		return ""
	}
	baseStr := "-script.param="
	params := strings.Join(q.ScriptParams, q.ScriptParamsDelimiter)
	return baseStr + url.QueryEscape(params)
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

func (q *FMQuery) recidString() string {
	if q.RecordId != fmNoRecordId && q.Action != FindAny {
		return "-recid=" + strconv.Itoa(q.RecordId)
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
		strArray = append(strArray, f.SimpleFieldsString())
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
			str := "-q" + strconv.Itoa(i) + "=" + url.QueryEscape(f.Name) + "&-q" + strconv.Itoa(i) + ".value=" + url.QueryEscape(f.ValueWithOp())
			strArray = append(strArray, str)
		}
		segments = append(segments, strings.Join(strArray, "&"))
	}
	return strings.Join(segments, "&")
}

func (q *FMQuery) QueryString() string {
	var startString = q.dbLayString() + WithAmp(q.scriptsString()) + WithAmp(q.scriptParamsString()) + WithAmp(q.ResponseLayout)
	switch q.Action {
	case Delete, Duplicate:
		return startString +
			WithAmp(q.recidString()) +
			q.Action.String()
	case Edit:
		return startString +
			WithAmp(q.recidString()) +
			WithAmp(q.simpleFieldsString()) +
			q.Action.String()
	case New:
		return startString +
			WithAmp(q.simpleFieldsString()) +
			q.Action.String()
	case FindAny:
		return startString +
			q.Action.String()
	case FindAll:
		return startString +
			WithAmp(q.sortFieldsString()) +
			WithAmp(q.maxSkipString()) +
			q.Action.String()
	case Find:
		if q.RecordId != fmNoRecordId {
			return startString +
				WithAmp(q.recidString()) + "-find"
		}
		return startString +
			WithAmp(q.sortFieldsString()) +
			WithAmp(q.maxSkipString()) +
			WithAmp(q.compoundQueryString()) +
			WithAmp(q.compoundFieldsString()) +
			q.Action.String()
	default:
		return ""
	}
}

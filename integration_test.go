package gofmcon

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type Table struct {
	Text string `json:"text"`
	Number float32 `json:"number"`
	Date time.Time `json:"date"`
	Time time.Time `json:"time"`
	Timestamp time.Time `json:"timestamp"`
	Container string `json:"container"`
	RepeatedContainer []string `json:"repeated_container"`
	NotEmptyNumber float32 `json:"not_empty_number"`
	SummaryNumber float32 `json:"summary_total_of_number"`
	NestedRecords []*NestedTable `json:"nested_table"`
}

type NestedTable struct {
	NestedText string `json:"nested_text"`
	RepeatedNumber []int `json:"nested_repeated_number"`
}

func Test1(t *testing.T) {
	fmHost := os.Getenv("FM_HOST")
	fmPort := os.Getenv("FM_PORT")
	fmUser := os.Getenv("FM_USER")
	fmPass := os.Getenv("FM_PASS")
	conn := NewFMConnector(fmHost, fmPort,fmUser,fmPass)
	q := NewFMQuery("test", "table", FindAll)
	q.WithResponseLayout("table")
	res, err := conn.Query(q)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	var tableRecs []Table
	for _, record := range res.Resultset.Records {
		var tableRec Table
		b, err := record.JsonFields()
		assert.NoError(t, err)
		err = json.Unmarshal(b, &tableRec)
		assert.NoError(t, err)
		tableRecs = append(tableRecs, tableRec)
	}
}

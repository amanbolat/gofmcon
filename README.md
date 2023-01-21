[![Go Report Card](https://goreportcard.com/badge/github.com/amanbolat/gofmcon)](https://goreportcard.com/report/github.com/amanbolat/gofmcon)

## FileMaker Server connector for golang
This library provide access to FileMaker Server using XML Web publishing

This library is a port of https://github.com/PerfectlySoft/Perfect-FileMaker on golang.

## In Production
We use this lib in our company in production for some public APIs and feature migration to other DBs like Postgres.
It's some kind of bridge.

API **could change** and sometimes might not be documented well. So look for commits and updates. 

## Tests
FileMaker is not Postgres or MySQL, so we cannot run docker and test automatically. Maybe we could run EC2 with Windows and install FileMaker Server on it.
If you have any ideas, open new issue or contact me.

## Installation

```
go get github.com/amanbolat/gofmcon
```
Then add the line below in your code 
```go
import "github.com/amanbolat/gofmcon"
```

## Example

In main.go
```go
package main

import (
    "encoding/json"
fm "github.com/amanbolat/gofmcon"
    "log"
    "github.com/kelseyhightower/envconfig"
    "fmt"
    "errors"
)

// config represents all the configuration we need in order to
// create a new FMConnector and establish the connection with 
// FileMaker database 
type config struct {
    FmHost          string `split_words:"true" required:"true"`
    FmUser          string `split_words:"true" required:"true"`
    FmPort          string `split_words:"true" required:"true"`
    FmDatabaseName  string `split_words:"true" required:"true"`
    FmPass          string `split_words:"true" required:"true"`
}

type postStore struct {
    fmConn *fm.FMConnector
    dbName string
}

type Post struct {
    Author string `json:"Author"`
    Title string `json:"Title"`
    Content string `json:"Content"`
}

func (p *Post) Populate(record *fm.Record) {
    p.Author = record.Field("author")
    p.Title = record.Field("title")
    p.Content = record.Field("content")
}

func main() {
    var conf = &config{}
    err := envconfig.Process("", conf)
    if err != nil {
        log.Fatal(err)
    }
    
    fmConn := fm.NewFMConnector(conf.FmHost, conf.FmPort, conf.FmUser, conf.FmPass)
    store := postStore{fmConn: fmConn, dbName: conf.FmDatabaseName}
    
    posts, err := store.GetAllPosts(fmConn)
    if err != nil {                                    
        log.Fatal(err)
    }
    
    fmt.Print(posts)
}

func (ps *postStore) GetAllPosts() ([]Post, error) {
	var posts []Post

	q := fm.NewFMQuery(ps.dbName, "posts_list_layout", fm.FindAll)
	fmset, err := ps.fmConn.Query(q)                        
	if err != nil {                                         
		return posts, errors.New("failed to get posts")
	}

    // Populate it with record
	for _, r := range fmset.Resultset.Records { 
		p := Post{}
        
		b, _ := r.JsonFields()
        _ = json.Unmarshal(b, &p)
		posts = append(posts, p)
	}

	return posts, nil
}
```


### Get a single record
```go
    q := fm.NewFMQuery(databaseName, layout_name, fm.Find)
    q.WithFields(
        fm.FMQueryField{Name: "field_name", Value: "001", Op: fm.Equal},
    ).Max(1)
```

### Check for FileMaker internal error
```go
    fmSet, err := fmConn.Query(q)
    if err != nil {
        if err.Error() == fmt.Sprintf("FileMaker_error: %s", fm.FileMakerErrorCodes[401]) {
            // your code
        }
    
        // else do something
    }
```


### Create record
```go
    q := fm.NewFMQuery(databaseName, layout_name, fm.New)
    q.WithFields(
        fm.FMQueryField{Name: "field_name", Value: "some_value"},
    )
    
    fmSet, err := fmConn.Query(q)
```


### Sort by some field
```go
    q.WithSortFields(fm.FMSortField{Name: "some_field", Order: fm.Descending})
```

### Update record
Your object should have FileMaker record id to update record in database. Please see more in FileMaker documentation.
```go
    q := fm.NewFMQuery(databaseName, layout_name, fm.Edit)
    q.WithFields(
        fm.FMQueryField{Name: "field_name", Value: "some_new_value"},
    )
    q.WithRecordId(updated_object.FMRecordID)
```


### Use some script
```go
    // SCRIPT_DELIMITER can be '|', '_' or any other symbol that will be
    // parsed on FileMaker side to get all the parameters from the string
    q.WithPostFindScripts(SCRIPT_NAME, strings.Join([]string{param_1, param_2, param_3}, SCRIPT_DELIMITER))
```

## TODO

- [x] Add tests
- [x] Add methods to get information about database and layouts

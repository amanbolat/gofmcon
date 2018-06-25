##Filemaker Server connector for golang
This library provide access to Filemaker Server using XML Web publishing

**It's mirror of repo on Bitbucket!** 

This library is a port of https://github.com/PerfectlySoft/Perfect-FileMaker, but in golang.

## In Production
We use this lib in our company in production for some public APIs and feature migration to other DBs like Postgres.
It's some kind of bridge.

API **could change** and sometimes might not be documented well. So look for commits and updates. 

## Tests
FileMaker is not Postgres or MySQL so we cannot run docker and test it everywhere. May be we could run EC2 with Windows and installed FileMaker server on it.
If you have any ideas, open new issue or contact me.

## Install

```
go get bitbucket.org/amanbolat/gofmcon
```
then in your code 
```go
import "bitbucket.org/amanbolat/gofmcon"
```

## Example


In main.go
```go
package main

import (
    fm "bitbucket.org/amanbolat/gofmcon"
    "log"
    "github.com/kelseyhightower/envconfig"
    "fmt"
    "errors"
)

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
    Author string
    Title string
    Content string
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

func (ps *postStore) GetAllPosts(conn *fm.FMConnector) (*[]Post, error) {
	var posts []Post

	q := fm.NewFMQuery(ps.dbName, "Posts_list", fm.FindAll)  // Create query
	fmset, err := ps.fmConn.Query(q)                        // Make request with query
	if err != nil {                                         // Check for errors
		return &posts, errors.New("Failed get posts: " + err.Error())
	}

	for _, r := range fmset.Resultset.Records[0:] {         // Iterate through records
		p := Post{}
		p.Populate(&r)                                      // Populate it with record
		posts = append(posts, p)
	}

	return &posts, nil                                      // Return posts
}
```


### Get single record
```go
    q := fm.NewFMQuery(databaseName, layout_name, fm.Find)
    q.WithFields(
        fm.FMQueryField{Name: "field_name", Value: "001", Op: fm.Equal},
    ).WithMaxRecords(1)
```

### Get filemaker error
```go
    fmSet, err := r.conn.Query(q)
    if err != nil {
        if err.Error() == fmt.Sprintf("filemaker_error: %s", fm.FileMakerErrorCodes[401]) {
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
    // SCRIPT_DELIMITER can be '|', '_' or any other symbol for
    // On FileMaker side input should be parsed to get every parameter
    // that is why we need delimiter
    q.WithPostFindScripts("script_name")
    q.WithScriptParams(SCRIPT_DELIMITER, param_1, param_2, param_3)
```

### Container fields
You should use `ContainerField` method to get slice of container field urls(image, file...) or it could be base64 encoded content.
```go
    containerURLs := append(containerURLs, fmRecord.ContainerField("container_field_name")...)
```


#TODO

- [ ] Add tests
- [ ] Add methods to get information about database and layouts
- [ ] Add documentation 

##Filemaker Server connector for golang
This library provide access to Filemaker Server using XML Web publishing

## Install

```
go get github.com/amanbolat/gofmcon
```
then in your code 
```go
import "github.com/amanbolat/gofmcon"
```

## Example


In main.go
```go
package main

import (
  fm "github.com/amanbolat/gofmcon"
)

const (
  DBName = "Posts"
  DBUser = "User123"
  DBPassword = "Pass123"
  DBHost = "myfmserver.com"
  Port = ""
)

var FMdb *fm.FMConnector

func main() {
  posts, err := GetAllPosts(FMdb)
  if err != nil {                                    
    log.Fatal(err)
	}
}

func InitFMdb(host string, port string, username string, password string) (e error) {
  FMdb, e = fm.NewFMConnector(host, port, username, password)
	return
}

func GetAllPosts(connector *fm.FMConnector) (*[]models.Post, error) {
	var posts []models.Post

	q := fm.NewFMQuery(dbName, "Posts_list", fm.FindAll)  // Create query
	fmset, err := connector.Query(q)                      // Make request with query
	if err != nil {                                       // Check for errors
		return &posts, errors.New("Failed get posts: " + err.Error())
	}

	for _, r := range fmset.Resultset.Records[0:] {       // Iterate through records
		s := models.Post{}
		s.Populate(&r)                                      // Populate it with record
		shipments = append(shipments, s)
	}

	return &posts, nil                                    // Retrun posts
}
```

in models.go
```go
package models

import (
  fm "github.com/amanbolat/gofmcon"
)

type Post struct {
  ID string
  Date string
}

func (p *Post) Populate(fmrecord *fm.Record) {
	p.ID = fmrecord.ToMap()["ID"]
	p.Date = fmrecord.ToMap()[""]
}
```

#TODO

- [ ] Add tests
- [ ] Add methods to get information about database and layouts
- [ ] Add documentation 

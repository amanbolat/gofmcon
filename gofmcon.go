package gofmcon

import (
	_ "encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	fmiPath = "fmi/xml/fmresultset.xml"
)

type FMConnector struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewFMConnector(host string, port string, username string, password string) (*FMConnector, error) {
	newF := &FMConnector{
		host,
		port,
		username,
		password,
	}

	var newURL = &url.URL{}
	newURL.Scheme = "http"
	newURL.Host = host
	if port != "" {
		newURL.Host += ":" + port
	}
	newURL.Path = fmiPath
	newURL.RawQuery = newURL.Query().Encode() + "&" + FMDBnames
	request, err := http.NewRequest("GET", newURL.String(), nil)
	if err != nil {
		return newF, err
	}
	request.SetBasicAuth(username, password)
	request.Header.Set("User-Agent", "Gopher crossasia api v1")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return newF, errors.New("Failed create new FMConnector: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return newF, errors.New("Failed create new FMConnector with status: %v\n" + res.Status)
	}

	return newF, err
}

const (
	FMDBnames = "-dbnames"
)

// func newFMURL(f *FMConnector, q *FMQuery) *url.URL {
// 	var port string
// 	if f.Port != "" {
// 		port = ":" + f.Port
// 	}
// 	var fmURL = &url.URL{}
// 	fmURL.Scheme = "http"
// 	fmURL.Host = f.Host + port
// 	fmURL.Path = fmiPath

// 	query := fmURL.Query()
// 	query.Set("-db", f.Database)
// 	query.Set("-lay", q.Layout)
// 	for k, v := range q.Query {
// 		query.Set(k, v)
// 	}
// 	fmURL.RawQuery = query.Encode() + "&" + q.Action.String()

// 	return fmURL
// }

// func (fmc *FMConnector) Request(q *FMQuery) (FMResultset, error) {
// 	queryResultset := FMResultset{}
// 	fmURL := newFMURL(fmc, q)

// 	request, err := http.NewRequest("GET", fmURL.String(), nil)
// 	if err != nil {
// 		fmt.Printf("Error on creating request: %v\n", err)
// 	}
// 	request.Header.Set("User-Agent", "Gopher crossasia api v1")
// 	request.SetBasicAuth(fmc.Username, fmc.Password)
// 	//fmt.Println(request)
// 	client := &http.Client{}
// 	res, err := client.Do(request)
// 	if err != nil {
// 		return queryResultset, errors.New("Failed get response: " + err.Error())
// 	}
// 	defer res.Body.Close()

// 	b, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return queryResultset, errors.New("Failed parse response: " + err.Error())
// 	}

// 	uError := xml.Unmarshal(b, &queryResultset)
// 	if uError != nil {
// 		return queryResultset, errors.New("Unmarshal parse error: " + uError.Error())
// 	}
// 	return queryResultset, nil
// }

func (fmc *FMConnector) Query(q *FMQuery) (FMResultset, error) {
	resultSet := FMResultset{}
	queryURL := fmc.makeURL(q)

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		fmt.Printf("Error on creating request: %v\n", err)
	}
	request.Header.Set("User-Agent", "Gopher crossasia api v1")
	request.SetBasicAuth(fmc.Username, fmc.Password)

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return resultSet, errors.New("Failed get response: " + err.Error())
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resultSet, errors.New("Failed parse response: " + err.Error())
	}

	err = xml.Unmarshal(b, &resultSet)
	if err != nil {
		return resultSet, errors.New("Unmarshal parse error: " + err.Error())
	}
	return resultSet, nil
}

func (fmc *FMConnector) makeURL(q *FMQuery) string {
	var newURL = &url.URL{}
	newURL.Scheme = "http"
	newURL.Host = fmc.Host
	if fmc.Port != "" {
		newURL.Host += ":" + fmc.Port
	}
	newURL.Path = fmiPath
	fmt.Printf("%v\n", newURL.String() + "?" + q.QueryString())
	return newURL.String() + "?" + q.QueryString()
}

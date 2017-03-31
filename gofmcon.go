package gofmcon

import (
	"encoding/xml"
	"errors"
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
		return newF, errors.New("gofmcon: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return newF, errors.New("gofmcon: failed connect fm server")
	}

	return newF, err
}

const (
	FMDBnames = "-dbnames"
)

func (fmc *FMConnector) Query(q *FMQuery) (FMResultset, error) {
	resultSet := FMResultset{}
	queryURL := fmc.makeURL(q)

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return resultSet, errors.New("gofmcon: " + err.Error())
	}
	request.Header.Set("User-Agent", "Gopher crossasia api v1")
	request.SetBasicAuth(fmc.Username, fmc.Password)

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return resultSet, errors.New("gofmcon: " + err.Error())
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resultSet, errors.New("gofmcon: " + err.Error())
	}

	err = xml.Unmarshal(b, &resultSet)
	if err != nil {
		return resultSet, errors.New("gofmcon: " + err.Error())
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
	return newURL.String() + "?" + q.QueryString()
}

package conn

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
)

type ConnectionString struct {
	LocalAddr      string
	RemmoteAddr    string
	Url            string
	LocalInterface string
	Type           ConnType
	Id             string
}

func NewURL(c ConnectionString) string {
	urlA, err := url.Parse(c.Url)
	if err != nil {
		log.Fatal(err)
	}
	values := urlA.Query()

	s := reflect.ValueOf(&c).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		values.Add(typeOfT.Field(i).Name, fmt.Sprintf("%v", f.Interface()))
	}
	urlA.RawQuery = values.Encode()

	return urlA.String()
}

func GetConnFromParameter(values url.Values) ConnectionString {

	conn := ConnectionString{}
	s := reflect.ValueOf(&conn).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		//values.Add(typeOfT.Field(i).Name, fmt.Sprintf("%v", f.Interface()))
		f := s.Field(i)
		ftype := f.Type().String()
		value := values.Get(typeOfT.Field(i).Name)

		if ftype == "string" || ftype == "conn.ConnType" {
			reflect.ValueOf(&conn).Elem().Field(i).SetString(value)
		}

	}

	return conn
}

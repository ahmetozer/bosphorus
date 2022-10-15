package client

import (
	"errors"
	"strings"
)

type arrFlag []string

func (i *arrFlag) String() string {
	return ""
}

func (i *arrFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var tcpFlag arrFlag

type connType string

type connectionString struct {
	localAddr   string
	remmoteAddr string
	host        string
}

func parseConnectionString(config string) (connectionString, error) {
	s := strings.Split(config, ";")

	if len(s) != 3 {
		return connectionString{}, errors.New("property does not meet struct")
	}

	return connectionString{
		localAddr:   s[0],
		host:        s[1],
		remmoteAddr: s[2],
	}, nil

}

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
	url         string
}

func parseConnectionString(config string) (connectionString, error) {

	// Example connection string per argument
	// 1 localAddr					2 host						3 remoteAddr
	// localhost:8119;http://nyc-sv-nss.ahmet.engineer/ws/tcp;localhost:8118

	s := strings.Split(config, ";")

	if len(s) != 3 {
		return connectionString{}, errors.New("property does not meet struct")
	}

	return connectionString{
		localAddr:   s[0],
		url:         s[1],
		remmoteAddr: s[2],
	}, nil

}

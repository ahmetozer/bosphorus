package client

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

const (
	tcpKeepAliveInterval = 10 * time.Second
	defaultDialTimeout   = 15 * time.Second
)

var wsPortMap = map[string]string{"ws": "80", "wss": "443", "http": "80", "https": "443"}

func wsDialAddress(location *url.URL) string {
	if _, ok := wsPortMap[location.Scheme]; ok {
		if _, _, err := net.SplitHostPort(location.Host); err != nil {
			return net.JoinHostPort(location.Host, wsPortMap[location.Scheme])
		}
	}
	return location.Host
}

func dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	d := &net.Dialer{KeepAlive: tcpKeepAliveInterval}
	return d.DialContext(ctx, network, addr)
}

func contextDialer(ctx context.Context) *net.Dialer {
	dialer := &net.Dialer{Cancel: ctx.Done(), KeepAlive: tcpKeepAliveInterval}
	if deadline, ok := ctx.Deadline(); ok {
		dialer.Deadline = deadline
	} else {
		dialer.Deadline = time.Now().Add(defaultDialTimeout)
	}
	return dialer
}

type quicStatus uint8

const (
	noInfo uint8 = 0
	noQuicSupport
	quicSupported
)

func newWSSocket(wsConfig *websocket.Config) (net.Conn, error) {
	var conn net.Conn
	var err error
	ctx := context.TODO()
	switch wsConfig.Location.Scheme {
	case "ws", "http":
		conn, err = dialContext(ctx, "tcp", wsDialAddress(wsConfig.Location))
	case "wss", "https":
		dialer := contextDialer(ctx)
		conn, err = tls.DialWithDialer(dialer, "tcp", wsDialAddress(wsConfig.Location), wsConfig.TlsConfig)
	default:
		err = websocket.ErrBadScheme
	}
	if err != nil {
		return nil, err
	}
	return conn, nil
}

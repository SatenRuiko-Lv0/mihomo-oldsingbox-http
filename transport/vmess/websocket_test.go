package vmess

import (
	"net"
	"net/url"
	"strings"
	"testing"

	"github.com/metacubex/http"
)

type captureConn struct {
	net.Conn
	written strings.Builder
}

func (c *captureConn) Write(b []byte) (int, error) {
	return c.written.Write(b)
}

func TestWriteWebsocketClientRequestLegacyOrder(t *testing.T) {
	request := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "ws",
			Host:   "192.0.2.1:80",
			Path:   "/",
		},
		Host:   "example.com",
		Header: make(http.Header),
	}
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Sec-WebSocket-Version", "13")
	request.Header.Set("Sec-WebSocket-Key", "test-key")

	conn := &captureConn{}
	if err := writeWebsocketClientRequest(conn, request); err != nil {
		t.Fatal(err)
	}

	expected := strings.Join([]string{
		"GET / HTTP/1.1",
		"Host: example.com",
		"Upgrade: websocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Version: 13",
		"Sec-WebSocket-Key: test-key",
		"User-Agent: Go-http-client/1.1",
		"",
		"",
	}, "\r\n")

	if actual := conn.written.String(); actual != expected {
		t.Fatalf("unexpected request:\n%s", actual)
	}
}

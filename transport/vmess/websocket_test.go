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

func TestWriteWebsocketClientRequestMinimalOrder(t *testing.T) {
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
		"",
		"",
	}, "\r\n")

	actual := conn.written.String()
	if actual != expected {
		t.Fatalf("unexpected request:\n%s", actual)
	}
	if strings.Contains(actual, "User-Agent:") {
		t.Fatalf("request should not include default User-Agent:\n%s", actual)
	}
}

func TestWriteWebsocketClientRequestPreservesConfiguredHeaders(t *testing.T) {
	request := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "ws",
			Host:   "192.0.2.1:80",
			Path:   "/ws",
		},
		Host:   "example.com",
		Header: make(http.Header),
	}
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Sec-WebSocket-Version", "13")
	request.Header.Set("Sec-WebSocket-Key", "test-key")
	request.Header.Set("User-Agent", "configured-agent")
	request.Header.Set("X-Test", "configured-value")

	conn := &captureConn{}
	if err := writeWebsocketClientRequest(conn, request); err != nil {
		t.Fatal(err)
	}

	actual := conn.written.String()
	if !strings.Contains(actual, "\r\nUser-Agent: configured-agent\r\n") {
		t.Fatalf("missing configured User-Agent:\n%s", actual)
	}
	if !strings.Contains(actual, "\r\nX-Test: configured-value\r\n") {
		t.Fatalf("missing configured header:\n%s", actual)
	}
}

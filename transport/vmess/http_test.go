package vmess

import (
	"strings"
	"testing"
)

func TestWriteHTTPClientRequestLegacyDefault(t *testing.T) {
	conn := &captureConn{}
	payload := []byte("vmess-payload")

	err := writeHTTPClientRequest(conn, "PUT", "/", "example.com", nil, payload)
	if err != nil {
		t.Fatal(err)
	}

	expected := strings.Join([]string{
		"PUT / HTTP/1.1",
		"Host: example.com",
		"",
		"vmess-payload",
	}, "\r\n")

	actual := conn.written.String()
	if actual != expected {
		t.Fatalf("unexpected request:\n%s", actual)
	}
	if strings.Contains(actual, "Content-Length:") {
		t.Fatal("request should not include Content-Length")
	}
	if strings.Contains(actual, "User-Agent:") {
		t.Fatal("request should not include default User-Agent")
	}
	if strings.Contains(actual, "Connection:") {
		t.Fatal("request should not include default Connection")
	}
}

func TestWriteHTTPClientRequestPreservesConfiguredHeaders(t *testing.T) {
	conn := &captureConn{}
	payload := []byte("body")

	err := writeHTTPClientRequest(conn, "GET", "/path", "example.com", map[string][]string{
		"Host":       {"example.com"},
		"User-Agent": {"Go-http-client/1.1"},
		"Connection": {"keep-alive"},
		"X-Empty":    {},
	}, payload)
	if err != nil {
		t.Fatal(err)
	}

	actual := conn.written.String()
	if !strings.HasPrefix(actual, "GET /path HTTP/1.1\r\nHost: example.com\r\n") {
		t.Fatalf("unexpected prefix:\n%s", actual)
	}
	if !strings.Contains(actual, "\r\nUser-Agent: Go-http-client/1.1\r\n") {
		t.Fatalf("missing user-agent:\n%s", actual)
	}
	if !strings.Contains(actual, "\r\nConnection: keep-alive\r\n") {
		t.Fatalf("missing connection:\n%s", actual)
	}
	if strings.Contains(actual, "\r\nX-Empty:") {
		t.Fatalf("empty header list should be skipped:\n%s", actual)
	}
	if !strings.HasSuffix(actual, "\r\n\r\nbody") {
		t.Fatalf("payload is not appended after headers:\n%s", actual)
	}
	if strings.Contains(actual, "Content-Length:") {
		t.Fatal("request should not include Content-Length")
	}
}

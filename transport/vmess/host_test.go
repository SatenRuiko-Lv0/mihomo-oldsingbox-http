package vmess

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/metacubex/http"
)

type recordingConn struct {
	bytes.Buffer
	readTimeout time.Duration
}

func (c *recordingConn) Read(p []byte) (int, error) {
	if c.readTimeout > 0 {
		time.Sleep(c.readTimeout)
	}
	return 0, io.EOF
}

func (c *recordingConn) Write(p []byte) (int, error) {
	return c.Buffer.Write(p)
}

func (c *recordingConn) Close() error {
	return nil
}

func (c *recordingConn) LocalAddr() net.Addr {
	return dummyAddr("local")
}

func (c *recordingConn) RemoteAddr() net.Addr {
	return dummyAddr("remote")
}

func (c *recordingConn) SetDeadline(time.Time) error {
	return nil
}

func (c *recordingConn) SetReadDeadline(time.Time) error {
	return nil
}

func (c *recordingConn) SetWriteDeadline(time.Time) error {
	return nil
}

type dummyAddr string

func (a dummyAddr) Network() string {
	return string(a)
}

func (a dummyAddr) String() string {
	return string(a)
}

func TestHTTPConnWritesSingleStableHost(t *testing.T) {
	conn := &recordingConn{}
	stream := StreamHTTPConn(conn, &HTTPConfig{
		Method: "GET",
		Host:   "server.example",
		Path:   []string{"/test"},
		Headers: map[string][]string{
			"host":       {"front.example"},
			"Host":       {"ignored.example"},
			"User-Agent": {"mihomo-test"},
		},
	})

	if _, err := stream.Write([]byte("payload")); err != nil {
		t.Fatal(err)
	}

	request := conn.String()
	if got := strings.Count(strings.ToLower(request), "\r\nhost: "); got != 1 {
		t.Fatalf("expected one Host header, got %d in:\n%s", got, request)
	}
	if !strings.Contains(request, "\r\nHost: ignored.example\r\n") {
		t.Fatalf("expected canonical Host from explicit Host key, got:\n%s", request)
	}
	if strings.Contains(request, "front.example") {
		t.Fatalf("lower-case host header leaked as duplicate:\n%s", request)
	}
}

func TestWebSocketConnWritesSingleStableHost(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	result := make(chan string, 1)
	go func() {
		var buffer bytes.Buffer
		reader := bufio.NewReader(server)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			buffer.WriteString(line)
			if line == "\r\n" {
				break
			}
		}
		_ = server.Close()
		result <- buffer.String()
	}()

	_, err := StreamWebsocketConn(context.Background(), client, &WebsocketConfig{
		Host:    "server.example",
		Port:    "80",
		Path:    "/ws",
		Headers: http.Header{"host": {"front.example"}, "User-Agent": {"mihomo-test"}},
	})
	if err == nil {
		t.Fatal("expected handshake failure without server response")
	}
	client.Close()

	request := <-result
	if got := strings.Count(strings.ToLower(request), "\r\nhost: "); got != 1 {
		t.Fatalf("expected one Host header, got %d in:\n%s", got, request)
	}
	if !strings.Contains(request, "\r\nHost: front.example\r\n") {
		t.Fatalf("expected configured Host, got:\n%s", request)
	}
}

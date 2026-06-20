package vmess

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/textproto"
	"strings"

	"github.com/metacubex/randv2"
)

type httpConn struct {
	net.Conn
	cfg        *HTTPConfig
	reader     *bufio.Reader
	whandshake bool
}

type HTTPConfig struct {
	Method  string
	Host    string
	Path    []string
	Headers map[string][]string
}

// Read implements net.Conn.Read()
func (hc *httpConn) Read(b []byte) (int, error) {
	if hc.reader != nil {
		n, err := hc.reader.Read(b)
		return n, err
	}

	reader := textproto.NewConn(hc.Conn)
	// First line: GET /index.html HTTP/1.0
	if _, err := reader.ReadLine(); err != nil {
		return 0, err
	}

	if _, err := reader.ReadMIMEHeader(); err != nil {
		return 0, err
	}

	hc.reader = reader.R
	return reader.R.Read(b)
}

// Write implements io.Writer.
func (hc *httpConn) Write(b []byte) (int, error) {
	if hc.whandshake {
		return hc.Conn.Write(b)
	}

	path := "/"
	if len(hc.cfg.Path) > 0 {
		path = hc.cfg.Path[randv2.IntN(len(hc.cfg.Path))]
	}

	method := hc.cfg.Method
	if method == "" {
		method = "PUT"
	}
	host := pickHTTPHeader(hc.cfg.Host, hc.cfg.Headers, "Host")
	if err := writeHTTPClientRequest(hc.Conn, method, path, host, hc.cfg.Headers, b); err != nil {
		return 0, err
	}
	hc.whandshake = true
	return len(b), nil
}

func (hc *httpConn) Close() error {
	return hc.Conn.Close()
}

func StreamHTTPConn(conn net.Conn, cfg *HTTPConfig) net.Conn {
	return &httpConn{
		Conn: conn,
		cfg:  cfg,
	}
}

func writeHTTPClientRequest(conn net.Conn, method string, path string, host string, headers map[string][]string, payload []byte) error {
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	var buffer bytes.Buffer
	buffer.WriteString(method)
	buffer.WriteByte(' ')
	buffer.WriteString(path)
	buffer.WriteString(" HTTP/1.1\r\n")
	buffer.WriteString("Host: ")
	buffer.WriteString(host)
	buffer.WriteString("\r\n")

	for key, values := range headers {
		if strings.EqualFold(key, "Host") || len(values) == 0 {
			continue
		}
		buffer.WriteString(key)
		buffer.WriteString(": ")
		buffer.WriteString(values[randv2.IntN(len(values))])
		buffer.WriteString("\r\n")
	}

	buffer.WriteString("\r\n")
	buffer.Write(payload)
	n, err := conn.Write(buffer.Bytes())
	if err == nil && n != buffer.Len() {
		err = io.ErrShortWrite
	}
	return err
}

func pickHTTPHeader(fallback string, headers map[string][]string, name string) string {
	for key, values := range headers {
		if strings.EqualFold(key, name) && len(values) > 0 {
			return values[randv2.IntN(len(values))]
		}
	}
	return fallback
}

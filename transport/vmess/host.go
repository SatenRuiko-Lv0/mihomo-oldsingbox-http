package vmess

import (
	"net/textproto"
	"sort"
	"strings"

	"github.com/metacubex/http"
)

func isHostHeader(key string) bool {
	return strings.EqualFold(key, "Host")
}

func firstHeaderValue(values []string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func requestHostFromHeaders(defaultHost string, headers map[string][]string) string {
	if value := firstHeaderValue(headers["Host"]); value != "" {
		return value
	}

	canonicalHost := textproto.CanonicalMIMEHeaderKey("Host")
	if canonicalHost != "Host" {
		if value := firstHeaderValue(headers[canonicalHost]); value != "" {
			return value
		}
	}

	keys := make([]string, 0, len(headers))
	for key := range headers {
		if isHostHeader(key) {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		if value := firstHeaderValue(headers[key]); value != "" {
			return value
		}
	}

	return strings.TrimSpace(defaultHost)
}

func deleteHostHeaders(headers http.Header) {
	for key := range headers {
		if isHostHeader(key) {
			delete(headers, key)
		}
	}
}

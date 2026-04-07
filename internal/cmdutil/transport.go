// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package cmdutil

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/larksuite/cli/internal/util"
)

// RetryTransport is an http.RoundTripper that retries on 5xx responses
// and network errors. MaxRetries defaults to 0 (no retries).
type RetryTransport struct {
	Base       http.RoundTripper
	MaxRetries int
	Delay      time.Duration // base delay for exponential backoff; defaults to 500ms
}

func (t *RetryTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return util.FallbackTransport()
}

func (t *RetryTransport) delay() time.Duration {
	if t.Delay > 0 {
		return t.Delay
	}
	return 500 * time.Millisecond
}

// RoundTrip implements http.RoundTripper.
func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base().RoundTrip(req)
	if t.MaxRetries <= 0 {
		return resp, err
	}

	for attempt := 0; attempt < t.MaxRetries; attempt++ {
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		// Clone request for retry
		cloned := req.Clone(req.Context())
		if req.Body != nil && req.GetBody != nil {
			cloned.Body, _ = req.GetBody()
		}
		delay := t.delay() * (1 << uint(attempt))
		time.Sleep(delay)
		resp, err = t.base().RoundTrip(cloned)
	}
	return resp, err
}

// UserAgentTransport is an http.RoundTripper that sets the User-Agent header.
// Used in the SDK transport chain to override the SDK's default User-Agent.
type UserAgentTransport struct {
	Base http.RoundTripper
}

func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set(HeaderUserAgent, UserAgentValue())
	for headerKey, headerValue := range extraHeadersFromEnv() {
		req.Header.Set(headerKey, headerValue)
	}
	if t.Base != nil {
		return t.Base.RoundTrip(req)
	}
	return util.FallbackTransport().RoundTrip(req)
}

func extraHeadersFromEnv() map[string]string {
	headers := map[string]string{}
	if key, value, ok := readExtraHeaderEnv("LARKSUITE_CLI_EXTRA_HEADER_KEY", "LARKSUITE_CLI_EXTRA_HEADER_VALUE"); ok {
		headers[key] = value
	}
	for i := 2; ; i++ {
		keyName := fmt.Sprintf("LARKSUITE_CLI_EXTRA_HEADER_KEY_%d", i)
		valueName := fmt.Sprintf("LARKSUITE_CLI_EXTRA_HEADER_VALUE_%d", i)
		key, value, ok := readExtraHeaderEnv(keyName, valueName)
		if !ok {
			if os.Getenv(keyName) == "" && os.Getenv(valueName) == "" {
				break
			}
			continue
		}
		headers[key] = value
	}
	return headers
}

func readExtraHeaderEnv(keyName, valueName string) (key, value string, ok bool) {
	key = os.Getenv(keyName)
	value = os.Getenv(valueName)
	if key == "" || value == "" {
		return "", "", false
	}
	return key, value, true
}

// SecurityHeaderTransport is an http.RoundTripper that injects CLI security
// headers into every request. Shortcut headers are read from the request context.
type SecurityHeaderTransport struct {
	Base http.RoundTripper
}

func (t *SecurityHeaderTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return util.FallbackTransport()
}

// RoundTrip implements http.RoundTripper.
func (t *SecurityHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	for k, vs := range BaseSecurityHeaders() {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}
	// Shortcut headers are propagated via context (see section 5.6 of the design doc).
	if name, ok := ShortcutNameFromContext(req.Context()); ok {
		req.Header.Set(HeaderShortcut, name)
	}
	if eid, ok := ExecutionIdFromContext(req.Context()); ok {
		req.Header.Set(HeaderExecutionId, eid)
	}
	return t.base().RoundTrip(req)
}

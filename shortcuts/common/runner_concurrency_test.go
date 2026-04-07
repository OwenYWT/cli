// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package common

import (
	"net/http"
	"testing"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"

	"github.com/larksuite/cli/internal/cmdutil"
	"github.com/larksuite/cli/internal/core"
)

func TestRuntimeContextGetAPIClient_InitializesOnceConcurrently(t *testing.T) {
	cfg := &core.CliConfig{AppID: "test-app", AppSecret: "test-secret", Brand: core.BrandFeishu}
	entered := make(chan struct{}, 8)
	release := make(chan struct{})

	rt := &RuntimeContext{
		Config: cfg,
		Factory: &cmdutil.Factory{
			Config: func() (*core.CliConfig, error) {
				entered <- struct{}{}
				<-release
				return cfg, nil
			},
			LarkClient: func() (*lark.Client, error) {
				return nil, nil
			},
			HttpClient: func() (*http.Client, error) {
				return &http.Client{}, nil
			},
		},
	}

	errCh := make(chan error, 8)
	done := make(chan struct{}, 8)
	start := make(chan struct{})
	for range 8 {
		go func() {
			<-start
			_, err := rt.getAPIClient()
			errCh <- err
			done <- struct{}{}
		}()
	}

	close(start)

	select {
	case <-entered:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for API client initialization")
	}

	select {
	case <-entered:
		t.Fatal("getAPIClient initialized the API client more than once under concurrent access")
	case <-time.After(50 * time.Millisecond):
	}

	close(release)

	for range 8 {
		<-done
		if err := <-errCh; err != nil {
			t.Fatalf("getAPIClient returned error: %v", err)
		}
	}
}

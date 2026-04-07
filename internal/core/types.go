// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package core

import "os"

// LarkBrand represents the Lark platform brand.
// "feishu" targets China-mainland, "lark" targets international.
// Any other string is treated as a custom base URL.
type LarkBrand string

const (
	BrandFeishu LarkBrand = "feishu"
	BrandLark   LarkBrand = "lark"
)

// Endpoints holds resolved endpoint URLs for different Lark services.
type Endpoints struct {
	Open     string // e.g. "https://open.feishu.cn"
	Accounts string // e.g. "https://accounts.feishu.cn"
	MCP      string // e.g. "https://mcp.feishu.cn"
}

// ResolveEndpoints resolves endpoint URLs based on brand.
func ResolveEndpoints(brand LarkBrand) Endpoints {
	if endpoints, ok := resolveEndpointOverrides(brand); ok {
		return endpoints
	}
	return ResolveOnlineEndpoints(brand)
}

func resolveEndpointOverrides(brand LarkBrand) (Endpoints, bool) {
	switch brand {
	case BrandLark:
		return Endpoints{}, false
	default:
		openBaseURL := os.Getenv("LARKSUITE_CLI_FEISHU_OPEN_BASE_URL")
		accountsBaseURL := os.Getenv("LARKSUITE_CLI_FEISHU_ACCOUNTS_BASE_URL")
		mcpBaseURL := os.Getenv("LARKSUITE_CLI_FEISHU_MCP_BASE_URL")
		if openBaseURL == "" && accountsBaseURL == "" && mcpBaseURL == "" {
			return Endpoints{}, false
		}
		endpoints := ResolveOnlineEndpoints(brand)
		if openBaseURL != "" {
			endpoints.Open = openBaseURL
		}
		if accountsBaseURL != "" {
			endpoints.Accounts = accountsBaseURL
		}
		if mcpBaseURL != "" {
			endpoints.MCP = mcpBaseURL
		}
		return endpoints, true
	}
}

// ResolveOnlineEndpoints resolves the production endpoint URLs based on brand.
func ResolveOnlineEndpoints(brand LarkBrand) Endpoints {
	switch brand {
	case BrandLark:
		return Endpoints{
			Open:     "https://open.larksuite.com",
			Accounts: "https://accounts.larksuite.com",
			MCP:      "https://mcp.larksuite.com",
		}
	default:
		return Endpoints{
			Open:     "https://open.feishu.cn",
			Accounts: "https://accounts.feishu.cn",
			MCP:      "https://mcp.feishu.cn",
		}
	}
}

// ResolveOpenBaseURL returns the Open API base URL for the given brand.
func ResolveOpenBaseURL(brand LarkBrand) string {
	return ResolveEndpoints(brand).Open
}

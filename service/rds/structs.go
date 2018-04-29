package rds

import "time"

// Response ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route_config.html?highlight=virtual_hosts
type Response struct {
	ValidateClusters        bool          `json:"validate_clusters,omitempty"`
	VirtualHosts            []VirtualHost `json:"virtual_hosts"`
	InternalOnlyHeaders     []string      `json:"internal_only_headers,omitempty"`
	ResponseHeadersToRemove []string      `json:"response_headers_to_remove,omitempty"`
	// response_headers_to_add
	// request_headers_to_add
}

// VirtualHost ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/vhost
type VirtualHost struct {
	Name            string           `json:"name"`
	Domains         []string         `json:"domains"`
	Routes          []Route          `json:"routes"`
	RequireSSL      string           `json:"require_ssl,omitempty"`
	VirtualClusters []VirtualCluster `json:"virtual_clusters,omitempty"`
	RateLimits      []RateLimit      `json:"rate_limits,omitempty"`
	// request_headers_to_add
}

// VirtualCluster ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/vcluster#config-http-conn-man-route-table-vcluster
type VirtualCluster struct {
	Pattern string `json:"pattern"`
	Name    string `json:"name"`
	Method  string `json:"method"`
}

// RateLimit ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/rate_limits#config-http-conn-man-route-table-rate-limit-config
type RateLimit struct {
	Stage      int      `json:"stage,omitempty"`
	DisableKey string   `json:"disable_key,omitempty"`
	Actions    []Action `json:"actions"`
}

// Action ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/rate_limits#actions
type Action struct {
	Type string `json:"type,omitempty"`
}

// Route ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route#config-http-conn-man-route-table-route
type Route struct {
	Prefix              string        `json:"prefix,omitempty"`
	Path                string        `json:"path,omitempty"`
	Regex               string        `json:"regex,omitempty"`
	Cluster             string        `json:"cluster"`
	HostRedirect        string        `json:"host_redirect,omitempty"`
	PathRedirect        string        `json:"path_redirect,omitempty"`
	PrefixRewrite       string        `json:"prefix_rewrite,omitempty"`
	HostRewrite         string        `json:"host_rewrite,omitempty"`
	AutoHostRewrite     bool          `json:"auto_host_rewrite,omitempty"`
	CaseSensitive       bool          `json:"case_sensitive,omitempty"`
	UseWebsocket        bool          `json:"use_websocket,omitempty"`
	TimeoutMS           time.Duration `json:"timeout_ms,omitempty"`
	RetryPolicy         *RetryPolicy  `json:"retry_policy,omitempty"`
	Shadow              *Shadow       `json:"shadow,omitempty"`
	Priority            string        `json:"priority,omitempty"`
	Headers             []Header      `json:"headers,omitempty"`
	RateLimits          []RateLimit   `json:"rate_limits,omitempty"`
	IncludeVhRateLimits bool          `json:"include_vh_rate_limits,omitempty"`
	HashPolicy          *HashPolicy   `json:"hash_policy,omitempty"`
	Decorator           *Decorator    `json:"decorator,omitempty"`
	// cors
	// cluster_header
	// weighted_clusters
	// runtime
	// request_headers_to_add
	// opaque_config
}

// RetryPolicy ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route#config-http-conn-man-route-table-route-retry
type RetryPolicy struct {
	RetryOn         string        `json:"retry_on"`
	NumRetries      int           `json:"num_retries,omitempty"`
	PerTryTimeoutMS time.Duration `json:"per_try_timeout_ms,omitempty"`
}

// Shadow ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route#config-http-conn-man-route-table-route-shadow
type Shadow struct {
	Cluster    string `json:"cluster"`
	RuntimeKey string `json:"runtime_key,omitempty"`
}

// Header ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route#config-http-conn-man-route-table-route-headers
type Header struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
	Regex string `json:"regex,omitempty"`
}

// HashPolicy ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route#config-http-conn-man-route-table-hash-policy
type HashPolicy struct {
	HeaderName string `json:"header_name"`
}

// Decorator ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route#config-http-conn-man-route-table-decorator
type Decorator struct {
	Operation string `json:"operation"`
}

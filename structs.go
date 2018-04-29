package main

import (
	"time"
)

// ServiceDiscoveryResponse ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds.html?highlight=hosts#get--v1-registration-(string-%20service_name)
type ServiceDiscoveryResponse struct {
	Hosts []Host `json:"hosts"`
}

// Host ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds.html?highlight=hosts#host-json
type Host struct {
	IP   string    `json:"ip_address"`
	Port int       `json:"port"`
	Tags *HostTags `json:"tags,omitempty"`
}

// HostTags ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds.html?highlight=hosts#host-json
type HostTags struct {
	AZ                  string `json:"az,omitempty"`
	Canary              bool   `json:"canary,omitempty"`
	LoadBalancingWeight int    `json:"load_balancing_weight,omitempty"`
}

// ClusterDiscoveryResponse ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/cds#config-cluster-manager-cds-v1
type ClusterDiscoveryResponse struct {
	Clusters []Cluster `json:"clusters"`
}

// Cluster response ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/cluster#config-cluster-manager-cluster
type Cluster struct {
	Name                          string            `json:"name"`
	Type                          string            `json:"type"`
	ConnectTimeoutMS              time.Duration     `json:"connect_timeout_ms,omitempty"`
	PerConnectionBufferLimitBytes int               `json:"per_connection_buffer_limit_bytes,omitempty"`
	LBtype                        string            `json:"lb_type"`
	Hosts                         []Host            `json:"hosts,omitempty"`
	ServiceName                   string            `json:"service_name"`
	HealthCheck                   *HealthCheck      `json:"health_check,omitempty"`
	MaxRequestsPerConnection      int               `json:"max_requests_per_connection,omitempty"`
	CleanupIntervalMS             time.Duration     `json:"cleanup_interval_ms,omitempty"`
	DNSRefreshRateMS              time.Duration     `json:"dns_refresh_rate_ms,omitempty"`
	OutlierDetection              *OutlierDetection `json:"outlier_detection,omitempty"`
	// ring_hash_lb_config
	// circuit_breakers
	// ssl_context
	// features
	// http2_settings
	// dns_lookup_family
	// dns_resolvers
}

// OutlierDetection ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/cluster_outlier_detection#config-cluster-manager-cluster-outlier-detection
type OutlierDetection struct {
	Consecutive5xx                     int           `json:"consecutive_5xx,omitempty"`
	ConsecutiveGatewayFailure          int           `json:"consecutive_gateway_failure,omitempty"`
	IntervalMS                         time.Duration `json:"interval_ms,omitempty"`
	BaseJjectionTimeMS                 time.Duration `json:"base_ejection_time_ms,omitempty"`
	MaxEjectionPercent                 int           `json:"max_ejection_percent,omitempty"`
	EnforcingConsecutive5xx            int           `json:"enforcing_consecutive_5xx,omitempty"`
	EnforcingConsecutiveGatewayFailure int           `json:"enforcing_consecutive_gateway_failure,omitempty"`
	EnforcingSuccessRate               int           `json:"enforcing_success_rate,omitempty"`
	SuccessRateMinimumHosts            int           `json:"success_rate_minimum_hosts,omitempty"`
	SuccessRateRequestVolume           int           `json:"success_rate_request_volume,omitempty"`
	SuccessRateStdevFactor             int           `json:"success_rate_stdev_factor,omitempty"`
}

// HealthCheck ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/cluster_hc#config-cluster-manager-cluster-hc-v1
type HealthCheck struct {
	Type               string              `json:"type"`
	TimeoutMS          time.Duration       `json:"timeout_ms"`
	IntervalMS         time.Duration       `json:"interval_ms"`
	UnhealthyThreshold int                 `json:"unhealthy_threshold"`
	HealthyThreshold   int                 `json:"healthy_threshold"`
	Path               string              `json:"path,omitempty"`
	IntervalJitterMS   time.Duration       `json:"interval_jitter_ms,omitempty"`
	ServiceName        string              `json:"service_name,omitempty"`
	Send               []map[string]string `json:"send"`
	Receive            []map[string]string `json:"receive"`
	// redis_key
}

// RouteDiscoveryResponse ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route_config.html?highlight=virtual_hosts
type RouteDiscoveryResponse struct {
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

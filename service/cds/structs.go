package cds

import "time"

// Response ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/cds#config-cluster-manager-cds-v1
type Response struct {
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

package main

// ServiceDiscoveryResponse ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds.html?highlight=hosts#get--v1-registration-(string-%20service_name)
type ServiceDiscoveryResponse struct {
	Hosts []Host `json:"hosts"`
}

// Host ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds.html?highlight=hosts#host-json
type Host struct {
	IP   string   `json:"ip_address"`
	Port int      `json:"port"`
	Tags HostTags `json:"tags,omitempty"`
}

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
	Name                     string `json:"name"`
	Type                     string `json:"type"`
	ServiceName              string `json:"service_name"`
	ConnectTimeoutMS         int    `json:"connect_timeout_ms,omitempty"`
	MaxRequestsPerConnection int    `json:"max_requests_per_connection,omitempty"`
	LBtype                   string `json:"lb_type"`
}

// RouteDiscoveryResponse ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route_config.html?highlight=virtual_hosts
type RouteDiscoveryResponse struct {
	VirtualHosts []VirtualHost `json:"virtual_hosts"`
}

// VirtualHost ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/vhost
type VirtualHost struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
	Routes  []Route  `json:"routes"`
}

// Route ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/route#config-http-conn-man-route-table-route
type Route struct {
	Prefix       string                 `json:"prefix,omitempty"`
	Path         string                 `json:"path,omitempty"`
	Cluster      string                 `json:"cluster"`
	UseWebsocket bool                   `json:"use_websocket"`
	RetryPolicy  map[string]interface{} `json:"retry_policy,omitempty"`
}

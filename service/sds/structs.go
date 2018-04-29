package sds

import (
	"github.com/jippi/consul-envoy/service/cds"
)

// Response ...
// https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds.html?highlight=hosts#get--v1-registration-(string-%20service_name)
type Response struct {
	Hosts []cds.Host `json:"hosts"`
}

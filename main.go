package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

type ServiceDiscoveryResponse struct {
	Hosts []ServiceHost `json:"hosts"`
}

type ServiceHost struct {
	IP   string            `json:"ip_address"`
	Port int               `json:"port"`
	Tags map[string]string `json:"tags,omitempty"`
}

type ClusterDiscoveryResponse struct {
	Clusters Clusters `json:"clusters"`
}

type Clusters []Cluster

type Cluster struct {
	Name                     string `json:"name"`
	Type                     string `json:"type"`
	ServiceName              string `json:"service_name"`
	ConnectTimeoutMS         int    `json:"connect_timeout_ms,omitempty"`
	MaxRequestsPerConnection int    `json:"max_requests_per_connection,omitempty"`
	LBtype                   string `json:"lb_type"`
}

type RouteDiscoveryResponse struct {
	VirtualHosts []VirtualHost `json:"virtual_hosts"`
}

type VirtualHost struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
	Routes  []Route  `json:"routes"`
}

type Route struct {
	Prefix       string                 `json:"prefix,omitempty"`
	Path         string                 `json:"path,omitempty"`
	Cluster      string                 `json:"cluster"`
	UseWebsocket bool                   `json:"use_websocket"`
	RetryPolicy  map[string]interface{} `json:"retry_policy,omitempty"`
}

var (
	clusterResponse ClusterDiscoveryResponse
	routeResponse   RouteDiscoveryResponse
	serviceResponse sync.Map
	consulDomain    string
)

// https://github.com/lyft/discovery
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing PORT to listen on")
	}

	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Could not create consul client: %s", err)
	}

	node, err := consul.Agent().Self()
	if err != nil {
		log.Fatalf("Could not find 'self' from consul catalog: %s", err)
	}

	watchers(consul)

	var ok bool
	consulDomain, ok = node["DebugConfig"]["DNSDomain"].(string)
	if !ok {
		log.Fatal("Could not find consul domain")
	}

	router := mux.NewRouter()

	// CDS - Cluster discovery service - https://www.envoyproxy.io/docs/envoy/v1.5.0/api-v1/cluster_manager/cds#config-cluster-manager-cds-v1
	router.HandleFunc("/v1/clusters/{service_cluster}/{service_node}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Infof("/v1/clusters/%s/%s", params["service_cluster"], params["service_node"])
		json.NewEncoder(w).Encode(clusterResponse)
	})

	// RDS - Route discovery service - https://www.envoyproxy.io/docs/envoy/v1.5.0/configuration/http_conn_man/rds
	router.HandleFunc("/v1/routes/{route_config_name}/{service_cluster}/{service_node}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Infof("/v1/routes/%s/%s/%s", params["route_config_name"], params["service_cluster"], params["service_node"])
		json.NewEncoder(w).Encode(routeResponse)
	})

	// SDS - Service discovery service - https://www.envoyproxy.io/docs/envoy/v1.5.0/api-v1/cluster_manager/sds#config-cluster-manager-sds-api
	router.HandleFunc("/v1/registration/{service_name}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Infof("/v1/registration/%s", params["service_name"])
		payload, _ := serviceResponse.Load(params["service_name"])
		json.NewEncoder(w).Encode(payload)
	})

	if err := http.ListenAndServe("0.0.0.0:"+port, router); err != nil {
		log.Fatal(err)
	}
}

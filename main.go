package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

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

	var ok bool
	consulDomain, ok = node["DebugConfig"]["DNSDomain"].(string)
	if !ok {
		log.Fatal("Could not find consul domain")
	}

	servicesCh := make(chan map[string][]string, 0)
	go servicesReader(consul, servicesCh)
	go clusterAndRouteBuilder(consul, servicesCh)

	router := mux.NewRouter()

	// CDS - Cluster discovery service - https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/cds#config-cluster-manager-cds-v1
	router.HandleFunc("/v1/clusters/{service_cluster}/{service_node}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Infof("/v1/clusters/%s/%s", params["service_cluster"], params["service_node"])
		json.NewEncoder(w).Encode(clusterResponse)
	})

	// RDS - Route discovery service - https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/rds#config-http-conn-man-rds-v1
	router.HandleFunc("/v1/routes/{route_config_name}/{service_cluster}/{service_node}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Infof("/v1/routes/%s/%s/%s", params["route_config_name"], params["service_cluster"], params["service_node"])
		json.NewEncoder(w).Encode(routeResponse)
	})

	// SDS - Service discovery service - https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds#config-cluster-manager-sds-api
	router.HandleFunc("/v1/registration/{service_name}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Debugf("/v1/registration/%s", params["service_name"])
		payload, ok := serviceResponse.Load(params["service_name"])
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(payload)
	})

	// wait for loaders to complete
	time.Sleep(1 * time.Second)

	// Listen on HTTP
	if err := http.ListenAndServe("0.0.0.0:"+port, router); err != nil {
		log.Fatal(err)
	}
}

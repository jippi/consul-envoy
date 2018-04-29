package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"github.com/jippi/consul-envoy/service/cds"
	"github.com/jippi/consul-envoy/service/rds"
	"github.com/jippi/consul-envoy/service/sds"
	log "github.com/sirupsen/logrus"
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

	consulDomain, ok := node["DebugConfig"]["DNSDomain"].(string)
	if !ok {
		log.Fatal("Could not find consul domain")
	}

	cdsCh := make(chan map[string][]string, 10)
	rdsCh := make(chan map[string][]string, 10)
	sdsCh := make(chan map[string][]string, 10)
	go servicesReader(consul, cdsCh, rdsCh, sdsCh)

	cdsWorker := cds.NewWorker(consul, cdsCh)
	go cdsWorker.Start()

	rdsWorker := rds.NewWorker(consul, consulDomain, rdsCh)
	go rdsWorker.Start()

	sdsWorker := sds.NewWorker(consul, sdsCh)
	go sdsWorker.Start()

	router := mux.NewRouter()

	// CDS - Cluster discovery service - https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/cds#config-cluster-manager-cds-v1
	router.HandleFunc("/v1/clusters/{service_cluster}/{service_node}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Infof("/v1/clusters/%s/%s", params["service_cluster"], params["service_node"])
		json.NewEncoder(w).Encode(cdsWorker.Response())
	})

	// RDS - Route discovery service - https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/route_config/rds#config-http-conn-man-rds-v1
	router.HandleFunc("/v1/routes/{route_config_name}/{service_cluster}/{service_node}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Infof("/v1/routes/%s/%s/%s", params["route_config_name"], params["service_cluster"], params["service_node"])
		json.NewEncoder(w).Encode(rdsWorker.Response())
	})

	// SDS - Service discovery service - https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds#config-cluster-manager-sds-api
	router.HandleFunc("/v1/registration/{service_name}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Debugf("/v1/registration/%s", params["service_name"])
		payload, ok := sdsWorker.Response(params["service_name"])
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(payload)
	})

	// Listen on HTTP
	if err := http.ListenAndServe("0.0.0.0:"+port, router); err != nil {
		log.Fatal(err)
	}
}

func servicesReader(client *api.Client, cdsCh, rdsCh, sdsCh chan map[string][]string) {
	query := &api.QueryOptions{
		AllowStale: true,
		WaitIndex:  0,
		WaitTime:   5 * time.Minute,
	}

	for {
		log.Info("Reading services")
		services, meta, err := client.Catalog().Services(query)
		log.Info("Read services")
		if err != nil {
			log.Error(err)
			time.Sleep(jitter(5 * time.Second))
			continue
		}

		query.WaitIndex = meta.LastIndex
		cdsCh <- services
		rdsCh <- services
		sdsCh <- services
	}
}

func jitter(d time.Duration) time.Duration {
	const jitter = 0.30
	jit := 1 + jitter*(rand.Float64()*2-1)
	return time.Duration(jit * float64(d))
}

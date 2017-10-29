package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
)

type host struct {
	URL string `json:"url"`
}

type serviceHost struct {
	IP   string `json:"ip_address"`
	Port int    `json:"port"`
}

type sslContext struct {
}

type clusters []cluster

type cluster struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	ServiceName      string `json:"service_name"`
	ConnectTimeoutMS int    `json:"connect_timeout_ms"`
	LBtype           string `json:"lb_type"`
}

type virtualHosts []virtualHost

type virtualHost struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
	Routes  []route  `json:"routes"`
}

type route struct {
	Prefix       string `json:"prefix"`
	Cluster      string `json:"cluster"`
	UseWebsocket bool   `json:"use_websocket"`
}

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

	nodeIP, ok := node["Member"]["Addr"]
	if !ok {
		log.Fatal("Could not find node IP")
	}

	domain, ok := node["DebugConfig"]["DNSDomain"]
	if !ok {
		log.Fatal("Could not find consul domain")
	}

	router := mux.NewRouter()

	// RDS - https://www.envoyproxy.io/envoy/configuration/http_conn_man/rds
	router.HandleFunc("/v1/routes/{route_config_name}/{service_cluster}/{service_node}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		fmt.Printf("/v1/routes/%s/%s/%s\n", params["route_config_name"], params["service_cluster"], params["service_node"])

		res := make([]virtualHost, 0)
		seen := make(map[string]string)

		catalog, _, _ := consul.Catalog().Node(params["service_node"], &api.QueryOptions{AllowStale: true})
		for _, service := range catalog.Services {
			serviceName := service.Service

			if _, ok := seen[serviceName]; ok {
				continue
			}

			seen[serviceName] = serviceName

			vhost := virtualHost{
				Name: serviceName,
				Domains: []string{
					fmt.Sprintf("%s.service.%s", serviceName, domain),
					fmt.Sprintf("*.%s.service.%s", serviceName, domain),
				},
				Routes: []route{
					route{
						Cluster:      serviceName,
						Prefix:       "/",
						UseWebsocket: false,
					},
				},
			}

			res = append(res, vhost)
		}

		x := struct {
			VirtualHosts []virtualHost `json:"virtual_hosts"`
		}{res}
		d, _ := json.Marshal(x)
		w.Write(d)
	})

	// SDS - Service discovery service - https://www.envoyproxy.io/envoy/configuration/cluster_manager/sds_api
	router.HandleFunc("/v1/registration/{service_name}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		fmt.Printf("/v1/registration/%s\n", params["service_name"])

		hosts := make([]serviceHost, 0)
		checks, _, _ := consul.Health().Service(params["service_name"], "", true, &api.QueryOptions{AllowStale: true})

		for _, entry := range checks {
			// check if it's an valid IP
			if ip := net.ParseIP(entry.Service.Address); ip != nil {
				hosts = append(hosts, serviceHost{
					IP:   entry.Service.Address,
					Port: entry.Service.Port,
				})
				continue
			}

			// if not an IP, resolve the hostname
			ips, err := net.LookupIP(entry.Service.Address)
			if err != nil {
				continue
			}

			// add reach resolved IP from the hostname to the registration
			for _, ip := range ips {
				hosts = append(hosts, serviceHost{
					IP:   ip.String(),
					Port: entry.Service.Port,
				})
			}
		}

		// consturct the valid response
		response := struct {
			Hosts []serviceHost `json:"hosts"`
		}{hosts}

		bytes, _ := json.Marshal(response)
		w.Write(bytes)
	})

	// CDS - https://www.envoyproxy.io/envoy/configuration/cluster_manager/cds
	router.HandleFunc("/v1/clusters/{service_cluster}/{service_node}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		fmt.Printf("/v1/clusters/%s/%s\n", params["service_cluster"], params["service_node"])

		// list of IP + Port for a given service name
		clusterHosts := make(map[string]*[]host)

		// local clusters for the consul agent we are attached to
		localClusters := make([]cluster, 0)

		nodeCatalog, _, _ := consul.Catalog().Node(params["service_node"], &api.QueryOptions{AllowStale: true})
		for _, service := range nodeCatalog.Services {
			// Always construct the service host struct
			serviceHost := host{
				URL: fmt.Sprintf("tcp://%s:%d", nodeIP, service.Port),
			}

			// if we already have a host-map for the service in question,
			// append the current host to the existing list rather than
			// adding a new cluster
			if hosts, ok := clusterHosts[service.Service]; ok {
				*hosts = append(*hosts, serviceHost)
				continue
			}

			// Create a new hostMap for the new cluster
			clusterHosts[service.Service] = &[]host{serviceHost}

			// Construct the envoy cluster config
			c := cluster{
				Name:             service.Service,
				ServiceName:      service.Service,
				Type:             "sds",
				LBtype:           "least_request",
				ConnectTimeoutMS: 1000,
			}

			// Append the cluster
			localClusters = append(localClusters, c)
		}

		// consturct the valid response
		response := struct {
			Clusters clusters `json:"clusters"`
		}{localClusters}

		bytes, _ := json.Marshal(response)
		w.Write(bytes)
	})

	if err := http.ListenAndServe("0.0.0.0:"+port, router); err != nil {
		log.Fatal(err)
	}
}

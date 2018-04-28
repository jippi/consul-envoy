package main

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

func watchers(client *api.Client) {
	servicesCh := make(chan map[string][]string, 0)
	go servicesReader(client, servicesCh)
	go clusterAndRouteBuilder(client, servicesCh)
}

func servicesReader(client *api.Client, servicesCh chan map[string][]string) {
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
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		query.WaitIndex = meta.LastIndex
		servicesCh <- services
	}
}

type serviceBuilder struct {
	lastSeen time.Time
	closeCh  chan interface{}
	client   *api.Client
	service  string
}

func (b *serviceBuilder) work() {
	q := &api.QueryOptions{
		AllowStale: true,
		WaitIndex:  0,
		WaitTime:   5 * time.Minute,
	}

	defer serviceResponse.Delete(b.service)

	for {
		select {
		case <-b.closeCh:
			return

		default:
			log.Infof("Reading service health %s", b.service)
			backends, meta, err := b.client.Health().Service(b.service, "", true, q)
			if err != nil {
				fmt.Println(err)
				time.Sleep(1)
				continue
			}
			if q.WaitIndex == meta.LastIndex {
				log.Infof("Read service health %s (but no changes)", b.service)
				continue
			}
			log.Infof("Read service health %s (with changes)", b.service)

			q.WaitIndex = meta.LastIndex

			hosts := make([]ServiceHost, 0)
			for _, entry := range backends {
				if ip := net.ParseIP(entry.Service.Address); ip != nil {
					hosts = append(hosts, ServiceHost{
						IP:   entry.Service.Address,
						Port: entry.Service.Port,
					})
					continue
				}

				ips, err := net.LookupIP(entry.Service.Address)
				if err != nil {
					continue
				}

				for _, ip := range ips {
					hosts = append(hosts, ServiceHost{
						IP:   ip.String(),
						Port: entry.Service.Port,
					})
				}
			}

			serviceResponse.Store(b.service, ServiceDiscoveryResponse{Hosts: hosts})
		}
	}
}

func clusterAndRouteBuilder(client *api.Client, servicesCh chan map[string][]string) {
	running := make(map[string]*serviceBuilder)
	cleanup := time.NewTicker(10 * time.Minute)

	for {
		select {
		case <-cleanup.C:
			log.Warn("Starting cleanup")
			timeout := time.Now().Add(-1 * time.Hour)

			for name, builder := range running {
				if builder.lastSeen.After(timeout) {
					continue
				}

				log.Infof("Deleting service %s due to timeout", name)
				close(builder.closeCh)
				delete(running, name)
			}

		default:
			log.Info("Waiting for services")
			services := <-servicesCh
			log.Info("Got services")

			clusters := make([]Cluster, 0)
			vhosts := make([]VirtualHost, 0)

			for name := range services {
				if _, ok := running[name]; !ok {
					running[name] = &serviceBuilder{
						lastSeen: time.Now(),
						closeCh:  make(chan interface{}),
						client:   client,
						service:  name,
					}
					go running[name].work()
				}

				running[name].lastSeen = time.Now()

				clusters = append(clusters, Cluster{
					Name:             name,
					ServiceName:      name,
					Type:             "sds",
					LBtype:           "least_request",
					ConnectTimeoutMS: 180000,
				})

				vhosts = append(vhosts, VirtualHost{
					Name: name,
					Domains: []string{
						fmt.Sprintf("%s.service.%s", name, consulDomain),
					},
					Routes: []Route{
						Route{
							Cluster: name,
							Prefix:  "/",
							RetryPolicy: map[string]interface{}{
								"retry_on":    "5xx,connect-failure",
								"num_retries": 1,
							},
						},
					},
				})
			}

			clusterResponse = ClusterDiscoveryResponse{Clusters: clusters}
			routeResponse = RouteDiscoveryResponse{VirtualHosts: vhosts}
		}
	}
}

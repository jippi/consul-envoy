package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

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
			log.Error(err)
			time.Sleep(jitter(5 * time.Second))
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

func jitter(d time.Duration) time.Duration {
	const jitter = 0.30
	jit := 1 + jitter*(rand.Float64()*2-1)
	return time.Duration(jit * float64(d))
}

func (b *serviceBuilder) work() {
	q := &api.QueryOptions{
		AllowStale: true,
		WaitIndex:  0,
		WaitTime:   jitter(5 * time.Minute),
	}

	defer serviceResponse.Delete(b.service)
	logger := log.WithField("service", b.service)

	for {
		select {
		case <-b.closeCh:
			logger.Info("Shutting down builder")
			return

		default:
			logger.Info("Reading service health")
			backends, meta, err := b.client.Health().Service(b.service, "", true, q)
			if err != nil {
				logger.Error(err)
				time.Sleep(jitter(5 * time.Second))
				continue
			}

			if q.WaitIndex == meta.LastIndex {
				logger.Infof("Read service health (but no changes)")
				continue
			}
			logger.Infof("Read service health (with changes)")

			q.WaitIndex = meta.LastIndex

			hosts := make([]Host, 0)
			for _, entry := range backends {
				if ip := net.ParseIP(entry.Service.Address); ip != nil {
					hosts = append(hosts, Host{
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
					hosts = append(hosts, Host{
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

		case services := <-servicesCh:
			log.Info("Got services")

			clusters := make([]Cluster, 0)
			vhosts := make([]VirtualHost, 0)

			for name := range services {
				if _, ok := running[name]; !ok {
					log.Infof("Discovered new service %s", name)

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
					HealthCheck: &HealthCheck{
						Type:               "tcp",
						TimeoutMS:          3 * time.Millisecond,
						IntervalMS:         5 * time.Millisecond,
						UnhealthyThreshold: 1,
						HealthyThreshold:   1,
						Send:               []map[string]string{},
						Receive:            []map[string]string{},
					},
					OutlierDetection: &OutlierDetection{},
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
							RetryPolicy: &RetryPolicy{
								RetryOn:    "5xx,connect-failure",
								NumRetries: 1,
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

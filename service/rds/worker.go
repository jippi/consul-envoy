package rds

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

// Worker for RDS (Route Discovery Service)
type Worker struct {
	consul       *api.Client              // Consul API client
	consulDomain string                   // Consul domain
	serviceCh    chan map[string][]string // Consul services channel (with tags)
	stopCh       chan interface{}         // Stop channel
	response     Response                 // Pre-computed response for HTTP server
}

// NewWorker will return the struct for a RDS worker
func NewWorker(consul *api.Client, consulDomain string, serviceCh chan map[string][]string) *Worker {
	return &Worker{
		consul:       consul,
		consulDomain: consulDomain,
		serviceCh:    serviceCh,
		stopCh:       make(chan interface{}),
	}
}

// Start will start the RDS worker, listening for service channel changes
// and pre-build RDS HTTP response
func (w *Worker) Start() {
	for {
		select {
		case <-w.stopCh:
			return

		case services := <-w.serviceCh:
			log.Info("Got services")

			vhosts := make([]VirtualHost, 0)

			for name := range services {
				vhosts = append(vhosts, VirtualHost{
					Name: name,
					Domains: []string{
						fmt.Sprintf("%s.service.%s", name, w.consulDomain),
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

			w.response = Response{VirtualHosts: vhosts}
		}
	}
}

// Stop the RDS worker
func (w *Worker) Stop() {
	close(w.stopCh)
}

// Response will return the pre-computed RDS response
func (w *Worker) Response() Response {
	return w.response
}

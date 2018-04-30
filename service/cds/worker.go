package cds

import (
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

// Worker for CDS (Cluster Discovery Service)
type Worker struct {
	consul    *api.Client              // Consul API Client
	response  Response                 // Pre-computed response for HTTP server
	serviceCh chan map[string][]string // Consul services channel (with tags)
	stopCh    chan interface{}         // Stop channel
}

// NewWorker will return the struct for a CDS worker
func NewWorker(consul *api.Client, serviceCh chan map[string][]string) *Worker {
	return &Worker{
		consul:    consul,
		serviceCh: serviceCh,
	}
}

// Start will start the CDS worker, listening for service channel changes
// and pre-build CDS HTTP response
func (w *Worker) Start() {
	w.stopCh = make(chan interface{})

	for {
		select {
		case <-w.stopCh:
			return

		case services := <-w.serviceCh:
			log.Info("Got services")

			clusters := make([]Cluster, 0)

			for name := range services {
				clusters = append(clusters, Cluster{
					Name:             name,
					ServiceName:      name,
					Type:             "sds",
					LBtype:           "least_request",
					ConnectTimeoutMS: 3 * time.Minute,
					OutlierDetection: &OutlierDetection{},
				})
			}

			w.response = Response{Clusters: clusters}
		}
	}
}

// Stop the CDS worker
func (w *Worker) Stop() {
	close(w.stopCh)
}

// Response will return the pre-computed CDS response
func (w *Worker) Response() Response {
	return w.response
}

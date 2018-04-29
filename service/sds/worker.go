package sds

import (
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

// Worker for SDS (Service Discovery Service)
type Worker struct {
	consul    *api.Client              // Consul API Client
	response  sync.Map                 // Map of pre-computed SDS responses, one per cluster
	serviceCh chan map[string][]string // Consul services channel (with tags)
	stopCh    chan interface{}         // Stop channel
}

// NewWorker will return the struct for a SDS worker
func NewWorker(client *api.Client, serviceCh chan map[string][]string) *Worker {
	return &Worker{
		consul:    client,
		serviceCh: serviceCh,
		stopCh:    make(chan interface{}),
	}
}

// Start will start the SDS worker, listening for service channel changes
// and pre-build SDS HTTP responses
func (w *Worker) Start() {
	running := make(map[string]*serviceBuilder)
	cleanup := time.NewTicker(10 * time.Minute)

	for {
		select {
		case <-w.stopCh:
			log.Info("Shutting down worker")
			return

		case <-cleanup.C:
			log.Warn("Starting cleanup")
			timeout := time.Now().Add(-1 * time.Hour)

			for name, checker := range running {
				if checker.lastSeen.After(timeout) {
					continue
				}

				log.Infof("Deleting service %s due to timeout", name)
				close(checker.closeCh)
				delete(running, name)
			}

		case services := <-w.serviceCh:
			for name := range services {
				if _, ok := running[name]; !ok {
					log.Infof("Discovered new service %s", name)

					running[name] = &serviceBuilder{
						lastSeen: time.Now(),
						closeCh:  make(chan interface{}),
						client:   w.consul,
						service:  name,
						worker:   w,
					}

					go running[name].work()
				}

				running[name].lastSeen = time.Now()
			}
		}
	}
}

// Stop the CDS worker
func (w *Worker) Stop() {
	close(w.stopCh)
}

// Response will return the pre-computed SDS response for a specific service
func (w *Worker) Response(service string) (value interface{}, ok bool) {
	return w.response.Load(service)
}

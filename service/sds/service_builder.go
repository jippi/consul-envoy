package sds

import (
	"math/rand"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/jippi/consul-envoy/service/cds"
	log "github.com/sirupsen/logrus"
)

type serviceBuilder struct {
	lastSeen time.Time
	closeCh  chan interface{}
	client   *api.Client
	service  string
	worker   *Worker
}

func (c *serviceBuilder) work() {
	q := &api.QueryOptions{
		AllowStale: true,
		WaitIndex:  0,
		WaitTime:   jitter(5 * time.Minute),
	}

	defer c.worker.response.Delete(c.service)
	logger := log.WithField("service", c.service)

	for {
		select {
		case <-c.closeCh:
			logger.Info("Shutting down builder")
			return

		default:
			logger.Info("Reading service health")
			backends, meta, err := c.client.Catalog().Service(c.service, "", q)
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

			hosts := make([]cds.Host, 0)
			for _, entry := range backends {
				if ip := net.ParseIP(entry.Address); ip != nil {
					hosts = append(hosts, cds.Host{
						IP:   entry.ServiceAddress,
						Port: entry.ServicePort,
						Tags: &cds.HostTags{
							AZ: entry.NodeMeta["aws_instance_availability-zone"],
						},
					})
					continue
				}

				ips, err := net.LookupIP(entry.Address)
				if err != nil {
					continue
				}

				for _, ip := range ips {
					hosts = append(hosts, cds.Host{
						IP:   ip.String(),
						Port: entry.ServicePort,
						Tags: &cds.HostTags{
							AZ: entry.NodeMeta["aws_instance_availability-zone"],
						},
					})
				}
			}

			c.worker.response.Store(c.service, Response{Hosts: hosts})
		}
	}
}

func jitter(d time.Duration) time.Duration {
	const jitter = 0.30
	jit := 1 + jitter*(rand.Float64()*2-1)
	return time.Duration(jit * float64(d))
}

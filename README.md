## About consul-envoy

This project aim to be a quick'n'dirty way to get Envoy and Consul talking nicely, so any Consul service URL can be routed through evnoy.

The current code allows to route `*.service.consul` to the backends availble in the Consul catalog

Currently [`RDS - Route discovery service`](https://www.envoyproxy.io/envoy/configuration/http_conn_man/rds), [`SDS - Service discovery service`](https://www.envoyproxy.io/envoy/configuration/cluster_manager/sds_api) and [`CDS - Cluster discovery service`](https://www.envoyproxy.io/envoy/configuration/cluster_manager/cds) is implemented

### Project goal

Making using Envoy with Consul as easy as [fabio](https://github.com/fabiolb/fabio) and [traefik](https://github.com/containous/traefik), possible exposting additional configuration to envoy through Consul Service tags (similar to fabio `urlprefix-*` configuration).

### Configuration

- `PORT` (env) - the HTTP port to listen on (example: `8877`)
- `CONSUL_*` (env) - the default Consul environment variables is used when connecting to the Consul cluster. (e.g. `CONSUL_HTTP_ADDR`)

### Building

`make requirements` to install Go Vendor and fetch dependencies
`make install` to build the binary (`consul-envoy`) into `${GOPATH}/bin`
`make dist` to build platform specific binary into `./build/consul-enovy-${OS}-${ARCH}`

### Example envoy config

The configuration assume that this project is named `envoy-consul` in the Consul catalog and listens on port `8877`

```json
{
    "listeners": [
        {
            "address": "tcp://0.0.0.0:80",
            "filters": [
                {
                    "name": "http_connection_manager",
                    "config": {
                        "codec_type": "auto",
                        "stat_prefix": "http",
                        "use_remote_address": true,
                        "rds": {
                            "route_config_name": "default",
                            "refresh_delay_ms": 10000,
                            "cluster": "rds_http"
                        },
                        "filters": [
                            {
                                "name": "router",
                                "config": {}
                            }
                        ]
                    }
                }
            ]
        },
    ],
    "admin": {
        "access_log_path": "/dev/null",
        "address": "tcp://0.0.0.0:8001"
    },
    "cluster_manager": {
        "cds": {
            "cluster": {
                "name": "cds",
                "type": "logical_dns",
                "lb_type": "round_robin",
                "connect_timeout_ms": 1000,
                "hosts": [
                    {
                        "url": "tcp://consul-envoy.service.consul:8877"
                    }
                ]
            }
        },
        "sds": {
            "refresh_delay_ms": 5000,
            "cluster": {
                "name": "sds",
                "type": "logical_dns",
                "lb_type": "round_robin",
                "connect_timeout_ms": 1000,
                "hosts": [
                    {
                        "url": "tcp://consul-envoy.service.consul:8877"
                    }
                ]
            }
        },
        "clusters": [
            {
                "name": "rds_http",
                "type": "logical_dns",
                "lb_type": "round_robin",
                "connect_timeout_ms": 1000,
                "hosts": [
                    {
                        "url": "tcp://consul-envoy.service.consul:8877"
                    }
                ]
            }
        ]
    }
}
```

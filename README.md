### example envoy config

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
            },
            {
                "name": "rds_https",
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

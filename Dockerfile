FROM alpine

RUN apk --update upgrade && \
    apk add curl ca-certificates openssh-client && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ADD ./build/consul-envoy-linux-amd64 /consul-envoy

ENTRYPOINT ["/consul-envoy"]

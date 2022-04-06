# Auto-mtls

- injector may add the label, or user
- based on the label, Istiod should generate mTLS config
- should interoperate with Envoy sidecars. 

Istio is using PeerPolicy to determine the server-side config, 
default is 'PERMISSIVE' which doesn't work and is translated to plaintext.

In 'mixed mode' sidecar will add auto-mtls label, so thing will be broken.

Replacement for 'permissive' is to have a plaintext port and a 
strict mtls port - similar with having 2 gRPC servers, or hbone.

The label should be set only if certs are available, otherwise it 
should fail.

# Bootstrap for agent-less

- metrics under /metrics
- keep 15020, add health, etc

Istio adds:

```shell
ANNOTATIONS:
          prometheus.io/path: /stats/prometheus
          prometheus.io/port: "15020"
          prometheus.io/scrape: "true"
          sidecar.istio.io/status: '{"initContainers":["istio-init"],"containers":["istio-proxy"],"volumes":["istio-envoy","istio-data","istio-podinfo","istio-token","istiod-ca-cert"],"imagePullSecrets":null,"revision":"default"}'
          
INSTANCE_IPS: 10.48.1.61
          
        ISTIO_PROXY_SHA: af4ca235bf40fb11829f2f02cf678610641095a2
        ISTIO_VERSION: 1.14-alpha.bee1e58baac36ef75fdc460e401bf8f66db8c832
        LABELS:
          app: httpbin
          pod-template-hash: 74fb669cc6
          security.istio.io/tlsMode: istio
          service.istio.io/canonical-name: httpbin
          service.istio.io/canonical-revision: v1
          version: v1
```

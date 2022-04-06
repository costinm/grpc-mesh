# GKE implementation for gateway

- uses a proxy-only subnet - source address can be checked !

```shell
 kubectl delete crd udproutes.networking.x-k8s.io tlsroutes.networking.x-k8s.io tcproutes.networking.x-k8s.io httproutes.networking.x-k8s.io gateways.networking.x-k8s.io gatewayclasses.networking.x-k8s.io backendpolicies.networking.x-k8s.io
 
kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.3.0"  | kubectl apply -f -

gcloud compute networks subnets create gke_gw \
    --purpose=INTERNAL_HTTPS_LOAD_BALANCER \
    --role=ACTIVE \
    --region=us-central1 \
    --network=default \
    --range=10.127.254.0/24   
    
kubectl get gatewayclass
NAME          CONTROLLER                    AGE
gke-l7-gxlb   networking.gke.io/gateway     284d
gke-l7-rilb   networking.gke.io/gateway     284d
istio         istio.io/gateway-controller   283d

# For MC
gcloud alpha container hub ingress enable \
    --config-membership=/projects/costin-asm1/locations/global/memberships/big1 \
    --project=costin-asm1
```

```yaml
    annotations:
      kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"networking.x-k8s.io/v1alpha1","kind":"Gateway","metadata":{"annotations":{},"name":"internal-http","namespace":"istio-system"},"spec":{"gatewayClassName":"gke-l7-rilb","listeners":[{"port":80,"protocol":"HTTP","routes":{"kind":"HTTPRoute"}}]}}
      networking.gke.io/addresses: gkegw-lql5-istio-system-internal-http-vzhzvj4hqa1u
      networking.gke.io/backend-services: gkegw-lql5-istio-system-gw-serve404-80-0vppk4aql635
      networking.gke.io/firewalls: gkegw-l7--default
      networking.gke.io/forwarding-rules: gkegw-lql5-istio-system-internal-http-0rf972gnatfu
      networking.gke.io/health-checks: gkegw-lql5-istio-system-gw-serve404-80-0vppk4aql635
      networking.gke.io/last-reconcile-time: Wednesday, 23-Mar-22 18:06:29 UTC
      networking.gke.io/ssl-certificates: ""
      networking.gke.io/target-proxies: gkegw-lql5-istio-system-internal-http-0rf972gnatfu
      networking.gke.io/url-maps: gkegw-lql5-istio-system-internal-http-0rf972gnatfu
    creationTimestamp: "2022-03-23T17:58:58Z"

    status:
      addresses:
        - type: IPAddress
          value: 10.128.15.241
      conditions:
        - lastTransitionTime: "1970-01-01T00:00:00Z"
          message: Waiting for controller
          reason: NotReconciled
          status: "False"
          type: Scheduled


```

# TD


```shell
gcloud alpha network-services meshes import grpc-mesh \
  --source=mesh.yaml \
  --location=global
 
 ```


# Features

- discovery - clusters and endpoints resolved
- DestinationRule:
  - subsets - TODO: using services and splits
  - loadBalancer:  ROUND_ROBIN. TODO: add consistentHash
  - TLS mode DISABLE and ISTIO_MUTUAL. TODO: automatic or mesh wide
- VirtualService:
  - header and URI match
  - override destination host, subset
  - weighted shifting
- PeerAuthentication DISABLE/STRICT - TODO: auto-MTLS

Missing:
- fault, retry, timeouts, mirroring, rewrite rules

# Ports and mixed usage

For apps that use both gRPC and TCP/HTTP, we need to use a sidecar for TCP/HTTP but want to 
exclude the gRPC port from capture on both inbound and outbound directions. For inbound it is 
easy to use 'excludePorts' - but for outbound it is not, since the destination ports are unknown.

We can use a pair of dedicated port - 15081,15082 - for gRPC with plaintext and mTLS.  

# Internals

Controller.cc (xds/internal/xdsclient/controller) is set during controller.New(config, updateHandler, validator).
It creates dial options with keep alive 5 min, 20 sec to - and config.Creds.

Controller.run() handles retries with backoff, creates a 'stream', and handles recv on the stream.
The stream is a 'grpc.ClientStream' - with Sent/RecvMsg. 

NewStream creates NewAggregateDiscoveryServiceClient for each call, and StreamAggregatedResources()

recv() calls t.handleResponse(proto) - warns if the response is unsupported, and 
uses sendCh to send ack/nack.

There are 2 ways to integrate with the XDS implementation, one is to fake the client
connection the other is to add hooks and expose sendCh.Put() and a hook for handle.

## Dial() and egress

XDS is integrated using the 'xds:' prefix.

1. xds_resolver handles the xds scheme, and implements LDS and RDS - but no CDS or EDS. NewXDSResolverWithConfigForTesting
allows passing a custom config, otherwise default is used.
2. Dial options may include XdsCredentials, including 'UsesXDS' credentials
3. bootstrap config is used to create a listener name, which in turn is used for on-demand LDS.
4. watch_service handleLDSResp extracts MaxStreamDuration, HTTPFilters.
5. It does support InlineRouteConfig - which saves a RTT.
6. applyRouteConfigUpdate handles RDS - finds the matching host for the serviceName. 

Google API library has a func WithGRPCConn(conn *grpc.ClientConn) option to allow passing a custom dialed ClientConn - 
it seems useful to have a similar mechanism to intercept and custom-dial.


## Used fields

- LDS - MaxStreamDuration, HTTPFilters
- 


## LB

internal/balancer has the various options - but weighted_target_experimental depends on internal, can't be used directly.

balancer appear to be pluggable - balancer.Register() provides a named Builder, and the Balancer object
has resolver.State with the Addresses to use. 

config_test in clusterimpl has sample configs for supported clusters.

priorityBalancer UpdateClientConnState is called

resource_resolver.newEDSResolver implenents watching endpoints.

## Improvements on GRPC

### intercept the grpc.Dial() call

The typical pattern for connecting to a 'target' is to construct a set of options and 
call "grpc.Dial(target, opts...)". This requires the caller - which can be a library - to
know the details of how to authenticate and verify the identity of the target.

XDS provides an alternative mechanism, where a (dynamic) config is used.



# Troubleshooting

## grpcurl

- requires discovery - would need an aggregated discovery service for ingress.
- 

## ghz - load test

## grpcdebug

go install -v github.com/grpc-ecosystem/grpcdebug@latest

kubectl -n echo-grpc port-forward $(kubectl -n echo-grpc get pods -l app=echo-grpc,version=v1 -ojsonpath='{.items[0].metadata.name}') 17070

grpcdebug localhost:17070 xds status
2022/03/30 11:36:28 failed to fetch xds config: rpc error: code = Unimplemented desc = unknown service envoy.service.status.v3.ClientStatusDiscoveryService

grpcdebug localhost:17070 channelz channels
grpcdebug localhost:17070 channelz servers
socket, subchannel

2022/03/30 11:39:06 failed to fetch top channels: rpc error: code = Unimplemented desc = unknown service grpc.channelz.v1.Channelz

- it seems recommended to use the admin service on separate port - should be secured.
- 

May have errors: proto: google.protobuf.Any: unable to resolve "type.googleapis.com/envoy.extensions.retry.host.previous_hosts.v3.PreviousHostsPredicate": not found
(in RDS)


# Implementation



# Other info

Istio tests and support for gateway: pilot/pkg/config/kube/gateway/testdata

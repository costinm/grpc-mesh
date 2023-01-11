ROOT_DIR?=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUT?=${ROOT_DIR}/../out/grpc-mesh

GOSTATIC=CGO_ENABLED=0  GOOS=linux GOARCH=amd64 time  go build -ldflags '-s -w -extldflags "-static"' -o ${OUT}
export PATH:=$(PATH):${HOME}/go/bin

build:
	mkdir -p ${OUT}
	(cd echo && ${GOSTATIC} ./cmd/server ./cmd/client)
	(cd echo-micro && ${GOSTATIC} ./cmd/*)
	ls -l ${OUT}

DOCKER_REPO?=gcr.io/dmeshgate/grpcmesh
BASE_DISTROLESS?=gcr.io/distroless/static


_push:
		(export IMG=$(shell cd ${OUT} && tar -cf - ${PUSH_FILES} ${BIN} | \
    					  gcrane append -f - -b ${BASE_DISTROLESS} \
    						-t ${DOCKER_REPO}/${BIN}:latest \
    					   ) && \
    	gcrane mutate $${IMG} -t ${DOCKER_REPO}/${BIN}:latest --entrypoint /${BIN} \
    	)

push/uproxy:
	(cd echo-micro && ${GOSTATIC} ./cmd/uecho)
	$(MAKE) _push BIN=uecho

#gen-old:
#	protoc --go_out xds --go_opt=paths=source_relative -I xds xds/*.proto
#	protoc \
#		-I proto \
#		-I vendor/protoc-gen-validate \
#		$(find proto -name '*.proto')

proto-gen:
	cd proto && buf generate

deps:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

	# debug tool for std grpc - need http/tcp equivalent
	go install -v github.com/grpc-ecosystem/grpcdebug@latest
	# Test tool
	go install github.com/bojand/ghz/cmd/ghz@latest

#docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf lint

echo/istio:
	(cd echo; go install ./cmd/server)
	server

install-cni:
	helm repo add istio https://istio-release.storage.googleapis.com/charts
	helm repo update
	helm template istio-cni istio/cni -n kube-system --set cni.cniBinDir=/home/kubernetes/bin \
		--set cni.hub=gcr.io/gke-release/asm --set cni.tag=1.12.4-asm.2 |kubectl apply -f -

# Enable the K8S Gateway in the cluster
k8s-gw:
	kubectl get crd gateways.gateway.networking.k8s.io || { kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.4.0" | kubectl apply -f -; }

config-dump:
	curl -H"Authorization:Bearer $(k8stok -audience istio-ca)" http://istiod.svc.i.webinf.info/debug/configz

echo/port-forward:
	kubectl -n echo-grpc port-forward $(kubectl -n echo-grpc get pods -l version=v1 -ojsonpath='{.items[0].metadata.name}') 17171 &

echo/call:
	grpcurl -plaintext -d '{"url": "xds:///echo.echo-grpc.svc.cluster.local:7070", "count": 5}' :17171 proto.EchoTestService/ForwardEcho | jq -r '.output | join("")'  | grep Hostname
	grpcurl -d '{}' fortio.svc.i.webinf.info:443 proto.EchoTestService/Echo |  jq -r '.message'
	grpcurl -d '{"url": "xds:///echo.echo-grpc.svc.cluster.local:7070", "count": 1}' fortio.svc.i.webinf.info:443 proto.EchoTestService/ForwardEcho

status:
	kubectl -n istio-system get gateways.gateway.networking.k8s.io istio-ingressgateway  -o yaml
	kubectl -n echo-grpc get httproute.gateway.networking.k8s.io   -o yaml



ls/all:


td-setup: NEG_NAME=k8s1-5e434f9a-istio-system-hgate-istiod-15012-876c9370
td-setup:
	# backend created automatically via ServiceExport  or cloud.google.com/neg: '{"exposed_ports":{"8080":{}}}' annotation
	gcloud compute network-endpoint-groups list
	# kubectl -n fortio-asm get serviceimports.net.gke.io
	gcloud compute health-checks create grpc istiod-hc \
     	--use-serving-port
	gcloud compute backend-services create istiod-service    --global    --load-balancing-scheme=INTERNAL_SELF_MANAGED    --protocol=GRPC    --health-checks istiod-hc
	gcloud compute backend-services add-backend istiod-service \
       --global \
       --network-endpoint-group ${NEG_NAME} \
       --network-endpoint-group-zone us-central1-c \
       --balancing-mode RATE \
       --max-rate-per-endpoint 5
	gcloud  compute url-maps create istiod --default-service istiod-service
	gcloud compute url-maps add-path-matcher istiod  \
      --default-service istiod-service \
      --path-matcher-name istiod \
      --new-hosts istiod # This is the host that will be used in xds:// requests !
	gcloud compute target-grpc-proxies create istiod \
     --url-map istiod \
     --validate-for-proxyless
	gcloud compute forwarding-rules create istiod \
     --global \
     --load-balancing-scheme=INTERNAL_SELF_MANAGED \
     --address=0.0.0.0 \
     --target-grpc-proxy=istiod \
     --ports 15010 \
     --network default

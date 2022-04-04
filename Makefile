ROOT_DIR?=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUT?=${ROOT_DIR}/../out/grpc-mesh

GOSTATIC=CGO_ENABLED=0  GOOS=linux GOARCH=amd64 time  go build -ldflags '-s -w -extldflags "-static"' -o ${OUT}

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

gen-old:
	protoc --go_out xds --go_opt=paths=source_relative -I xds xds/*.proto
	protoc \
		-I proto \
		-I vendor/protoc-gen-validate \
		$(find proto -name '*.proto')

proto-gen:
	cd proto && buf generate

deps:
	go install -v github.com/grpc-ecosystem/grpcdebug@latest
	# Test tool
	go install github.com/bojand/ghz@latest

	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	# Used to build the protos, lint
	GO111MODULE=on go install  github.com/bufbuild/buf/cmd/buf@v1.1.0

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

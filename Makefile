ROOT_DIR?=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUT?=${ROOT_DIR}/../out/grpc-mesh

GOSTATIC=CGO_ENABLED=0  GOOS=linux GOARCH=amd64 time  go build -ldflags '-s -w -extldflags "-static"' -o ${OUT}

build:
	mkdir -p ${OUT}
	(cd echo && ${GOSTATIC} ./cmd/server ./cmd/client)
	ls -l ${OUT}

gen-old:
	protoc --go_out xds --go_opt=paths=source_relative -I xds xds/*.proto
	protoc \
		-I proto \
		-I vendor/protoc-gen-validate \
		$(find proto -name '*.proto')

gen:
	cd proto && buf build

deps:
	go install -v github.com/grpc-ecosystem/grpcdebug@latest
	# Test tool
	go install github.com/bojand/ghz@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
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

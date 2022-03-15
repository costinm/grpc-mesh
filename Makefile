ROOT_DIR?=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUT?=${ROOT_DIR}/../out/grpc-mesh

GOSTATIC=CGO_ENABLED=0  GOOS=linux GOARCH=amd64 time  go build -ldflags '-s -w -extldflags "-static"' -o ${OUT}

build:
	mkdir -p ${OUT}
	(cd echo && ${GOSTATIC} ./cmd/server ./cmd/client)
	ls -l ${OUT}

install-cni:
	helm repo add istio https://istio-release.storage.googleapis.com/charts
	helm repo update
	helm template istio-cni istio/cni -n kube-system --set cni.cniBinDir=/home/kubernetes/bin \
		--set cni.hub=gcr.io/gke-release/asm --set cni.tag=1.12.4-asm.2 |kubectl apply -f -

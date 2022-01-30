ROOT_DIR?=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUT?=${ROOT_DIR}/../out/grpc-mesh

GOSTATIC=CGO_ENABLED=0  GOOS=linux GOARCH=amd64 time  go build -ldflags '-s -w -extldflags "-static"' -o ${OUT}

build:
	mkdir -p ${OUT}
	(cd echo && ${GOSTATIC} ./cmd/server ./cmd/client)
	ls -l ${OUT}

GO ?= go
GOBUILD := GO15VENDOREXPERIMENT=1 ${GO} build

PACKAGE := gh
BINNAME := gh
BINNAME_UBUNTU := gh.ubuntu

all: build

build: ${BINNAME} ${BINNAME_UBUNTU}

run:
	./${BINNAME} \
			-owner="astronoka" \
			-repo="gh" \
			-sha="c47d8d77a6858bad827e0ac9d3be8d35f51f3dbb"

${BINNAME}: main.go
	${GOBUILD} -o ${BINNAME} ${PACKAGE}

${BINNAME_UBUNTU}: main.go
	GOOS="linux" GOARCH="amd64" ${GOBUILD} -o ${BINNAME_UBUNTU} ${PACKAGE}

clean:
	rm -rf ${BINNAME} ${BINNAME_UBUNTU}

.PHONY: all build run clean

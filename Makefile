PROG_NAME := "archive-exploder"
VERSION = 0.1.$(shell date +%Y%m%d.%H%M)
FLAGS := "-s -w -X main.version=${VERSION}"


build:
	#go mod tidy
	#go mod vendor
	CGO_ENABLED=0 go build -ldflags=${FLAGS} -o ${PROG_NAME} *.go

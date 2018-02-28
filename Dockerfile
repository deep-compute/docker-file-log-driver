FROM golang:1.8-alpine

RUN apk update && \
	apk add git

WORKDIR /go/src/github.com/deep-compute/docker-file-log-driver

RUN go get github.com/docker/docker/api/types/plugins/logdriver
RUN go get github.com/docker/docker/daemon/logger
RUN go get github.com/docker/docker/daemon/logger/loggerutils
RUN go get github.com/docker/go-plugins-helpers/sdk
RUN go get github.com/pkg/errors
RUN go get github.com/Sirupsen/logrus
RUN go get github.com/tonistiigi/fifo
RUN go get github.com/gogo/protobuf/io
RUN go get github.com/snowzach/rotatefilehook

COPY . /go/src/github.com/deep-compute/docker-file-log-driver
RUN go get
RUN go build --ldflags '-extldflags "-static"' -o /usr/bin/docker-file-log-driver

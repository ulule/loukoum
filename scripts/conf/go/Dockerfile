FROM golang:1.11.1-stretch

MAINTAINER thomas.leroux@ulule.com

ENV DEBIAN_FRONTEND noninteractive
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

RUN apt-get -y update \
    && apt-get upgrade -y \
    && apt-get -y install git \
    && mkdir -p /go/src/github.com/golang \
    && cd /go/src/github.com/golang \
    && go get -u -v github.com/stretchr/testify/require \
    && go get -u -v github.com/davecgh/go-spew/spew \
    && go get -u -v github.com/lib/pq \
    && go get -u -v github.com/jmoiron/sqlx \
    && go get -u -v github.com/alecthomas/gometalinter \
    && gometalinter --install \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY . /media/ulule/loukoum
WORKDIR /media/ulule/loukoum
RUN go mod download

CMD /bin/bash

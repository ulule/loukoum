FROM golang:1.10.0-stretch

MAINTAINER thomas.leroux@ulule.com

RUN apt-get -y update && apt-get -y install git
COPY . /go/src/github.com/ulule/loukoum
WORKDIR /go/src/github.com/ulule/loukoum

CMD /bin/bash

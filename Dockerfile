FROM golang
MAINTAINER Maxime Santerre
ENV PORT 5050

RUN go get github.com/tools/godep

ADD . /go/src/github.com/msanterre/canoe
WORKDIR /go/src/github.com/msanterre/canoe

RUN godep get github.com/msanterre/canoe
RUN go install github.com/msanterre/canoe

EXPOSE 5050

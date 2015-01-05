FROM golang
MAINTAINER Brian Morton "brian@mmm.hm"

ADD . /go/src/github.com/bmorton/deployster
RUN cd /go/src/github.com/bmorton/deployster && go get -v -d
RUN go install github.com/bmorton/deployster
ENTRYPOINT ["/go/bin/deployster"]
EXPOSE 3000

FROM golang:latest
WORKDIR /go/src/github.com/dairycart/dairyclient/v1

ADD . .

ENTRYPOINT ["go", "test", "-v", "-cover"]

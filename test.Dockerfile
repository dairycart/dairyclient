FROM golang:1.9rc2-alpine
WORKDIR /go/src/github.com/dairycart/dairyclient

ADD . .

ENTRYPOINT ["go", "test", "-v", "-cover"]
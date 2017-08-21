FROM golang:1.9rc2-alpine
WORKDIR /go/src/github.com/dairycart/dairyclient/v1

ADD . .

ENTRYPOINT ["go", "test", "-v", "-cover"]
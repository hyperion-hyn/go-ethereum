# syntax=docker/dockerfile:experimental

# Build Geth in a stock Go builder container
FROM golang:1.15-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git \
    pkgconf g++ gmp-dev openssl-dev

ADD ./go.* /go-ethereum/
ADD ./third_party/bls/go.* /go-ethereum/third_party/bls/
ADD ./third_party/mcl/go.* /go-ethereum/third_party/mcl/
RUN cd /go-ethereum && go mod vendor

ADD . /go-ethereum
RUN --mount=type=cache,target=/root/.cache/go-build ( cd /go-ethereum && make clean && make third_party && make geth )

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates \
    musl gmp openssl libstdc++ fish shadow \
    && usermod -s `which fish` root
COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/

EXPOSE 8545 8546 8547 30303 30303/udp
ENTRYPOINT ["geth"]

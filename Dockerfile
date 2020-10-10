# syntax=docker/dockerfile:experimental

# Build Geth in a stock Go builder container
FROM golang:1.15-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git \
    pkgconf g++ gmp-dev openssl-dev

ADD ./go.* /go-ethereum/
RUN --mount=type=cache,sharing=locked,id=go-cache,target=/go/pkg/mod \
    mkdir -p /go-ethereum/third_party/mcl /go-ethereum/third_party/bls && \
    (cd /go-ethereum/third_party/mcl && go mod init github.com/hyperion-hyn/mcl) && \
    (cd /go-ethereum/third_party/bls && go mod init github.com/hyperion-hyn/bls) && \
    (cd /go-ethereum && go mod download -x)

ADD . /go-ethereum
RUN cd /go-ethereum && make clean
RUN --mount=type=cache,sharing=locked,id=go-cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    ( cd /go-ethereum && go mod vendor && make third_party && make geth )

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates \
    musl gmp openssl libstdc++ fish bash shadow \
    && usermod -s `which fish` root
COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/

EXPOSE 8545 8546 8547 30303 30303/udp
ENTRYPOINT ["geth"]

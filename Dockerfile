# Build Geth in a stock Go builder container
FROM golang:1.15-alpine as builder

ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=${GOPROXY}

RUN wget -qO- http://rsync.alpinelinux.org/alpine/MIRRORS.txt | xargs -I{} -P 50 time -f "%e %C" wget -T 1 -q {}/MIRRORS.txt -O /dev/null 2>&1 | sed -n '/^[0-9]/s,.*\(http[^ ]*\)/MIRRORS.txt.*,\1,;s,/*$,,p' | head -n 1 | xargs -I{} sed -i 's,http.*/v,{}/v,g' /etc/apk/repositories
RUN apk add --no-cache make gcc musl-dev linux-headers git \
    pkgconf g++ gmp-dev openssl-dev

ADD . /go-ethereum
RUN cd /go-ethereum && make clean && make third_party && make all

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=${GOPROXY}

RUN wget -qO- http://rsync.alpinelinux.org/alpine/MIRRORS.txt | xargs -I{} -P 50 time -f "%e %C" wget -T 1 -q {}/MIRRORS.txt -O /dev/null 2>&1 | sed -n '/^[0-9]/s,.*\(http[^ ]*\)/MIRRORS.txt.*,\1,p' | head -n 1 | xargs -I{} sed -i 's,http.*/v,{}/v,g' /etc/apk/repositories
RUN apk add --no-cache ca-certificates musl gmp openssl libstdc++ fish shadow \
    && usermod -s `which fish` root

COPY --from=builder /go-ethereum/build/bin/* /usr/local/bin/

EXPOSE 8545 8546 8547 30303 30303/udp
ENTRYPOINT ["geth"]

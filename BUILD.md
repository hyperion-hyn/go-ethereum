# Build

## in container

````shell script
docker run -it \
  --name atlas-builder \
  --restart always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v `pwd`:/src/ -v `pwd`/build-alpine/bin:/src/build/bin/ \
  -v `pwd`/build-alpine/mcl/obj:/src/third_party/mcl/obj -v `pwd`/build-alpine/mcl/lib:/src/third_party/mcl/lib \
  -v `pwd`/build-alpine/bls/obj:/src/third_party/bls/obj -v `pwd`/build-alpine/bls/lib:/src/third_party/bls/lib \
  golang:1.15-alpine sh -c "apk add --no-cache make gcc musl-dev linux-headers git pkgconf g++ gmp-dev openssl-dev openssl-libs-static binutils upx git bash fish docker && cd /src && bash"

````
````shell script
make third_party && make geth && docker build -t hyperion/atlas -f ./Dockerfile.dev ./
````

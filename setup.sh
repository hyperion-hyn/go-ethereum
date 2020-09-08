#!/bin/bash

# https://stackoverflow.com/questions/4774054/reliable-way-for-a-bash-script-to-get-the-full-path-to-itself
BASE_DIR="$( cd "$( dirname "$0" )" >/dev/null 2>&1 && pwd -P )"

if [ ! -e third_party/mcl/go.mod ]; then
  pushd third_party/mcl > /dev/null && go mod init github.com/hyperion-hyn/mcl && popd > /dev/null
fi

if [ ! -e third_party/bls/go.mod ]; then
  pushd third_party/bls > /dev/null && go mod init github.com/hyperion-hyn/bls && popd > /dev/null
fi

make -C third_party/mcl lib/libmcl.a -j8
make -C third_party/bls lib/libbls384_256.a BLS_SWAP_G=1 -j8

export CGO_CFLAGS="-I${BASE_DIR}/third_party/mcl/include -I${BASE_DIR}/third_party/bls/include"
export CGO_LDFLAGS="-L${BASE_DIR}/third_party/mcl/lib -L${BASE_DIR}/third_party/bls/lib"

unset BASE_DIR

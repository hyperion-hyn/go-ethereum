.PHONY: generator abigen storegen
.PHONY: build-contract
.PHONY: test clean

BUILDDIR = data/build/contracts
SOURCEDIR = data/contracts
ABIDIR = .
STORAGEDIR = .
BINDIR = .
SOURCES = $(filter-out %/Migrations.sol, $(wildcard $(SOURCEDIR)/*.sol))
OBJECTS = $(patsubst $(SOURCEDIR)/%.sol, $(BUILDDIR)/%.json, $(SOURCES))
ABIS = $(patsubst $(BUILDDIR)/%.json, $(ABIDIR)/%_abi_test.go, $(OBJECTS))
STORAGES = $(patsubst $(BUILDDIR)/%.json, $(STORAGEDIR)/%_storage_test.go, $(OBJECTS))

all: generator build-contract test

abigen:
	go build ../../abigen

storegen:
	go build ../../storegen

generator: abigen storegen

build-contract: $(OBJECTS) $(ABIS) $(STORAGES)

$(BUILDDIR)/%.json: $(SOURCEDIR)/%.sol
	solc $< --overwrite --pretty-json --combined-json abi,bin,bin-runtime,storage-layout \
        | jq '.contracts[] | {"abi": (.abi | fromjson), "bytecode": .bin, "deployedBytecode": ."bin-runtime", "storage-layout": (."storage-layout" | fromjson) }' \
        > $@

$(ABIDIR)/%_abi_test.go: $(BUILDDIR)/%.json abigen
	cat $< | jq '.abi' | ./abigen --abi - --pkg test --type $(patsubst %.json, %, $(notdir $<))Wrapper --out $@

$(STORAGEDIR)/%_storage_test.go: $(BUILDDIR)/%.json storegen
	cat $< | jq '."storage-layout"' | ./storegen --storage-layout - --pkg test --out ./map3pool_storage_test.go

test:
	go test -v -count=1 ./

clean:
	rm -rf data/build/contracts/*
	rm -rf abigen storegen

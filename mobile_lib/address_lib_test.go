package mobile_lib

import (
	"testing"
)

func TestPublicKeyToEthAddress(t *testing.T) {
	addressLib := NewAddressLib()
	pub := "04a5ad8fc58f2afa2a210c156f8a1cfb06bc0b6b46e59ee8295a1fe8fc4ecc6517307f4f2ccea4a9bc0690e3fd1fc28e6d85a8b82e9c0b54cb008f9e4d57485522"
	ethAddress := addressLib.PublicKeyToEthAddress(pub)
	if ethAddress != "0x70247395aFFd13C2347aA8c748225f1bFeD2C32A" {
		t.Error("public key convert to eth address error")
	}
}

func TestPublicKeyToHynAddress(t *testing.T) {
	addressLib := NewAddressLib()
	pub := "04a5ad8fc58f2afa2a210c156f8a1cfb06bc0b6b46e59ee8295a1fe8fc4ecc6517307f4f2ccea4a9bc0690e3fd1fc28e6d85a8b82e9c0b54cb008f9e4d57485522"
	hynAddress := addressLib.PublicKeyToHynAddress(pub, true)
	if hynAddress != "hyn1wqj889d0l5fuydr64rr5sgjlr0ld9se20c3ckj" {
		t.Error("public key convert to hyn mainnet address error")
	}

	hynAddress = addressLib.PublicKeyToHynAddress(pub, false)
	if hynAddress != "thyn1wqj889d0l5fuydr64rr5sgjlr0ld9se2pdcukr" {
		t.Error("public key convert to hyn testnet address error")
	}
}

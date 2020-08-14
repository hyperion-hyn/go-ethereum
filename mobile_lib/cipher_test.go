package mobile_lib

import (
	"strings"
	"testing"
)

func TestCipher_GenKeyPair(t *testing.T) {
	cipher := NewCipher()
	keystr := cipher.GenKeyPair()
	keys := strings.Split(keystr, ",")
	t.Log(keystr)
	if len(keys) != 2 {
		t.Error("key pair make error")
	}
}

func TestCipher_EncryptDecrypt(t *testing.T) {
	cipher := NewCipher()
	keystr := cipher.GenKeyPair()
	keys := strings.Split(keystr, ",")
	if len(keys) != 2 {
		t.Error("key pair make error")
		return
	}
	prvStr := keys[0]
	pubStr := keys[1]

	message := "this is message"
	encryptedMessage := cipher.Encrypt(pubStr, message)
	t.Log(encryptedMessage)
	decryptedMessage := cipher.Decrypt(prvStr, encryptedMessage)
	t.Log(decryptedMessage)
	if message != decryptedMessage {
		t.Error("decrypted result not equal to origin message")
	}
}

func TestDeCompressPubkey(t *testing.T) {
	cipher := NewCipher()
	//eth address is 0x70247395aFFd13C2347aA8c748225f1bFeD2C32A
	pub := cipher.DeCompressPubkey("0x02a5ad8fc58f2afa2a210c156f8a1cfb06bc0b6b46e59ee8295a1fe8fc4ecc6517")
	if pub != "04a5ad8fc58f2afa2a210c156f8a1cfb06bc0b6b46e59ee8295a1fe8fc4ecc6517307f4f2ccea4a9bc0690e3fd1fc28e6d85a8b82e9c0b54cb008f9e4d57485522" {
		t.Error("uncompress error")
	}
}

func TestPublicKeyToEthAddress(t *testing.T) {
	cipher := NewCipher()
	pub := "04a5ad8fc58f2afa2a210c156f8a1cfb06bc0b6b46e59ee8295a1fe8fc4ecc6517307f4f2ccea4a9bc0690e3fd1fc28e6d85a8b82e9c0b54cb008f9e4d57485522"
	ethAddress := cipher.PublicKeyToEthAddress(pub)
	if ethAddress != "0x70247395aFFd13C2347aA8c748225f1bFeD2C32A" {
		t.Error("public key convert to eth address error")
	}
}

func TestPublicKeyToHynAddress(t *testing.T) {
	cipher := NewCipher()
	pub := "04a5ad8fc58f2afa2a210c156f8a1cfb06bc0b6b46e59ee8295a1fe8fc4ecc6517307f4f2ccea4a9bc0690e3fd1fc28e6d85a8b82e9c0b54cb008f9e4d57485522"
	hynAddress := cipher.PublicKeyToHynAddress(pub, true)
	if hynAddress != "hyn1wqj889d0l5fuydr64rr5sgjlr0ld9se20c3ckj" {
		t.Error("public key convert to hyn mainnet address error")
	}

	hynAddress = cipher.PublicKeyToHynAddress(pub, false)
	if hynAddress != "thyn1wqj889d0l5fuydr64rr5sgjlr0ld9se2pdcukr" {
		t.Error("public key convert to hyn testnet address error")
	}
}

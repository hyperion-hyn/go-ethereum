package mobile_lib

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
)

type AddressLib struct {
}

func (m *AddressLib) PublicKeyToEthAddress(pubStr string) string {
	pubKeyByte, err := hex.DecodeString(pubStr)
	if err != nil {
		return ""
	}
	publicKey, err := crypto.UnmarshalPubkey(pubKeyByte)
	if err != nil {
		return ""
	}
	address := crypto.PubkeyToAddress(*publicKey)
	return address.String()
}

func (m *AddressLib) PublicKeyToHynAddress(pubStr string, isMainNet bool) string {
	pubKeyByte, err := hex.DecodeString(pubStr)
	if err != nil {
		return ""
	}
	publicKey, err := crypto.UnmarshalPubkey(pubKeyByte)
	if err != nil {
		return ""
	}
	address := crypto.PubkeyToAddress(*publicKey)
	if isMainNet {
		return address.Bech32()
	} else {
		return address.Bech32T()
	}
}

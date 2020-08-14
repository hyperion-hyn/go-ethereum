package mobile_lib

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

type Cipher struct {
}

func NewCipher() *Cipher {
	return &Cipher{}
}

func (m *Cipher) GenKeyPair() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		print("generate crypt key get error: %v", err)
		return ""
	}
	prvStr := hex.EncodeToString(crypto.FromECDSA(privateKey))
	pubKey := privateKey.Public().(*ecdsa.PublicKey)
	pubStr := hex.EncodeToString(crypto.FromECDSAPub(pubKey))
	keyPairString := fmt.Sprintf("%s,%s", prvStr, pubStr)

	return keyPairString
}

func (m *Cipher) Encrypt(pubStr, message string) string {
	pubKeyByte, err := hex.DecodeString(pubStr)
	if err != nil {
		return ""
	}
	publicKey, err := crypto.UnmarshalPubkey(pubKeyByte)
	if err != nil {
		return ""
	}
	eciesPubKey := ecies.ImportECDSAPublic(publicKey)
	cipherBytes, err := ecies.Encrypt(rand.Reader, eciesPubKey, []byte(message), nil, nil)
	if err != nil {
		return ""
	}
	cipherText := base58.Encode(cipherBytes)
	return cipherText
}

func (m *Cipher) Decrypt(prvStr, cipherText string) string {
	originCipherText := base58.Decode(cipherText)
	prvKey, _ := crypto.HexToECDSA(prvStr)
	eciesPrvKey := ecies.ImportECDSA(prvKey)
	message, err := eciesPrvKey.Decrypt([]byte(originCipherText), nil, nil)
	if err != nil {
		return ""
	}
	return string(message)
}

func (m *Cipher) DeCompressPubkey(pubStrC string) string {
	pubStr, err := crypto.DecompressPubkey(hexutil.MustDecode(pubStrC))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(crypto.FromECDSAPub(pubStr))
}

func (m *Cipher) PublicKeyToEthAddress(pubStr string) string {
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

func (m *Cipher) PublicKeyToHynAddress(pubStr string, isMainNet bool) string {
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

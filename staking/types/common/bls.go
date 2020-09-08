package common

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperion-hyn/bls/ffi/go/bls"
	"github.com/pkg/errors"
)

const (
	BLSSignatureSizeInBytes = 96
	BLSVerificationStr      = "hyperion-hyn"
)

var (
	errBLSKeysNotMatchSigs = errors.New("bls keys and corresponding signatures could not be verified")
)

// BLSSignature defines the bls signature
type BLSSignature [BLSSignatureSizeInBytes]byte

type blsKey interface {
	ToLibBLSPublicKey(key *bls.PublicKey) error
}


// VerifyBLSKey checks if the public BLS key matches the BLS signature
func VerifyBLSKey(pubKey blsKey, pubKeySig *BLSSignature) error {
	if len(pubKeySig) == 0 {
		return errBLSKeysNotMatchSigs
	}

	blsPubKey := new(bls.PublicKey)
	if err := pubKey.ToLibBLSPublicKey(blsPubKey); err != nil {
		return errBLSKeysNotMatchSigs
	}

	msgSig := bls.Sign{}
	if err := msgSig.Deserialize(pubKeySig[:]); err != nil {
		return err
	}

	messageBytes := []byte(BLSVerificationStr)
	msgHash := crypto.Keccak256(messageBytes)
	if !msgSig.VerifyHash(blsPubKey, msgHash[:]) {
		return errBLSKeysNotMatchSigs
	}

	return nil
}


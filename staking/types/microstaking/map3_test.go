package microstaking

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/bls"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"strings"
)

type blsPubSigPair struct {
	pub BLSPublicKey_
	sig common2.BLSSignature
}

func makeBLSPubSigPairs(size int) []blsPubSigPair {
	pairs := make([]blsPubSigPair, 0, size)
	for i := 0; i != size; i++ {
		pairs = append(pairs, makeBLSPubSigPair())
	}
	return pairs
}

func makeBLSPubSigPair() blsPubSigPair {
	blsPriv := bls.RandPrivateKey()
	blsPub := blsPriv.GetPublicKey()
	msgHash := crypto.Keccak256([]byte(common2.BLSVerificationStr))
	blsSig := blsPriv.SignHash(msgHash)

	pub := BLSPublicKey_{}
	copy(pub.Key[:], blsPub.Serialize())

	var sig common2.BLSSignature
	copy(sig[:], blsSig.Serialize())

	return blsPubSigPair{pub, sig}
}

func getPubsFromPairs(pairs []blsPubSigPair, indexes []int) BLSPublicKeys_ {
	pubs := NewEmptyBLSKeys()
	for _, index := range indexes {
		pubs.Keys = append(pubs.Keys, &pairs[index].pub)
	}
	return pubs
}

func getSigsFromPairs(pairs []blsPubSigPair, indexes []int) []common2.BLSSignature {
	sigs := make([]common2.BLSSignature, 0, len(indexes))
	for _, index := range indexes {
		sigs = append(sigs, pairs[index].sig)
	}
	return sigs
}

func assertError(gotErr, expErr error) error {
	if (gotErr == nil) != (expErr == nil) {
		return fmt.Errorf("error unexpected [%v] / [%v]", gotErr, expErr)
	}
	if gotErr == nil {
		return nil
	}
	if !strings.Contains(gotErr.Error(), expErr.Error()) {
		return fmt.Errorf("error unexpected [%v] / [%v]", gotErr, expErr)
	}
	return nil
}